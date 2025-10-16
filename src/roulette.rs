use axum::{
    extract::Query,
    http::StatusCode,
    response::{IntoResponse, Response},
    Json,
};
use rand::seq::SliceRandom;
use serde::{Deserialize, Serialize};
use validator::Validate;
use crate::osm;

#[derive(Debug, Deserialize, Validate)]
pub struct GetPubParams {
    #[validate(range(min = -90.0, max = 90.0))]
    latitude: f64,
    #[validate(range(min = -180.0, max = 180.0))]
    longitude: f64,
    #[validate(range(min = 1, max = 5000))]
    radius: i16,
}

#[derive(Debug, Serialize)]
pub struct PubResponse {
    name: String,
    latitude: f64,
    longitude: f64,
}

pub async fn find_random_pub(Query(params): Query<GetPubParams>) -> Result<Json<PubResponse>, AppError> {
    params.validate()?;

    const AMENITY: &str = "pub";
    let elements = osm::get_elements_in_radius(
        params.latitude,
        params.longitude,
        params.radius,
        AMENITY.to_string(),
    )
        .await.map_err(|e| AppError::OsmError(e))?;

    let random_element = elements
        .choose(&mut rand::thread_rng())
        .ok_or(AppError::NoPubsFound)?;

    let response = PubResponse {
        name: random_element.names.name.clone(),
        latitude: random_element.location.latitude,
        longitude: random_element.location.longitude,
    };

    Ok(Json(response))
}

pub enum AppError {
    Validation(validator::ValidationErrors),
    NoPubsFound,
    OsmError(String),
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
            AppError::OsmError(e) => (
                StatusCode::INTERNAL_SERVER_ERROR,
                format!("An external service error occurred: {}", e),
            ),
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