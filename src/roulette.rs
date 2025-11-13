use axum::{
    extract::Query,
    http::StatusCode,
    response::{IntoResponse, Response},
    Json,
};
use rand::seq::SliceRandom;
use serde::{Deserialize, Serialize};
use tracing::error;
use validator::Validate;
use crate::osm;

const PUB_AMENITY: &str = "pub";

#[derive(Debug, Deserialize, Validate)]
pub struct GetPubParams {
    #[validate(range(min = -90.0, max = 90.0))]
    latitude: f64,
    #[validate(range(min = -180.0, max = 180.0))]
    longitude: f64,
    #[validate(range(min = 1, max = 1000))]
    radius: i16,
    #[validate(range(min = 0.0, max = 359.0))]
    bearing: Option<f64>,
}

#[derive(Debug, Deserialize, Validate, Clone, Copy)]
pub struct GetCrawlParams {
    #[validate(range(min = -90.0, max = 90.0))]
    latitude: f64,
    #[validate(range(min = -180.0, max = 180.0))]
    longitude: f64,
    #[validate(range(min = 1, max = 1000))]
    radius: i16,
    #[validate(range(min = 0.0, max = 359.0))]
    bearing: Option<f64>,
    #[validate(range(min = 0, max = 4))]
    num_pubs: i16
}

#[derive(Debug, Serialize)]
pub struct PubResponse {
    name: String,
    latitude: f64,
    longitude: f64,
}

pub async fn find_random_crawl(Query(params): Query<GetCrawlParams>) -> Result<Json<Vec<PubResponse>>, AppError> {
    params.validate()?;

    let mut all_elements = osm::get_elements_in_radius(
        params.latitude,
        params.longitude,
        params.radius,
        params.bearing,
        PUB_AMENITY.to_string(),
    )
        .await
        .map_err(|e| AppError::OsmError(e))?;

    let num_pubs_requested = params.num_pubs as usize;
    let num_pubs_found = all_elements.len();

    if num_pubs_found == 0 {
        return Err(AppError::NoPubsFound);
    }

    if num_pubs_found < num_pubs_requested {
        return Err(AppError::NotEnoughPubsFound {
            requested: params.num_pubs,
            found: num_pubs_found,
        });
    }

    all_elements.shuffle(&mut rand::thread_rng());

    let responses = all_elements
        .into_iter()
        .take(num_pubs_requested)
        .map(|element| PubResponse {
            name: element.names.name,
            latitude: element.location.latitude,
            longitude: element.location.longitude,
        })
        .collect();

    Ok(Json(responses))
}

pub async fn find_random_pub(Query(params): Query<GetPubParams>) -> Result<Json<PubResponse>, AppError> {
    params.validate()?;

    let response = get_random_pub(params.latitude, params.longitude, params.radius, params.bearing).await?;

    Ok(Json(response))
}

async fn get_random_pub(lat: f64, lon: f64, radius: i16, bearing: Option<f64>) -> Result<PubResponse, AppError> {
    let elements = osm::get_elements_in_radius(
        lat,
        lon,
        radius,
        bearing,
        PUB_AMENITY.to_string(),
    )
        .await.map_err(|e| AppError::OsmError(e))?;

    let random_element = elements
        .choose(&mut rand::thread_rng())
        .ok_or(AppError::NoPubsFound)?;

    Ok(PubResponse {
        name: random_element.names.name.clone(),
        latitude: random_element.location.latitude,
        longitude: random_element.location.longitude,
    })
}

pub enum AppError {
    Validation(validator::ValidationErrors),
    NoPubsFound,
    OsmError(String),
    NotEnoughPubsFound { requested: i16, found: usize }
}

impl IntoResponse for AppError {
    fn into_response(self) -> Response {
        let (status, error_message) = match self {
            AppError::Validation(e) => (
                StatusCode::BAD_REQUEST,
                format!("Input validation failed: {}", e),
            ),
            AppError::NoPubsFound => (
                StatusCode::NOT_FOUND,
                "No pubs found within the specified radius.".to_string(),
            ),
            AppError::NotEnoughPubsFound  {requested, found} => (
                StatusCode::NOT_FOUND,
                format!("Not enough pubs found. Request {}, but only {} found.", requested, found),
            ),
            AppError::OsmError(e) => {
                error!("External service error: {}", e);

                (
                    StatusCode::INTERNAL_SERVER_ERROR,
                    "An internal service error occurred.".to_string(),
                )
            },
        };

        let body = Json(serde_json::json!({ "error": error_message }));
        (status, body).into_response()
    }
}

impl From<validator::ValidationErrors> for AppError {
    fn from(err: validator::ValidationErrors) -> Self {
        AppError::Validation(err)
    }
}