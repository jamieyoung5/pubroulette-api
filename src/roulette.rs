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
use vercel_runtime::{StatusCode as VercelStatusCode};

const PUB_AMENITY: &str = "pub";

#[derive(Debug, Serialize, Deserialize)]
pub struct PubResponse {
    pub name: String,
    pub latitude: f64,
    pub longitude: f64,
}

pub async fn get_random_crawl(length: i16, lat: f64, lon: f64, radius: i16, bearing: Option<f64>) -> Result<Vec<PubResponse>, AppError> {

    let mut all_elements = osm::get_elements_in_radius(
        lat,
        lon,
        radius,
        bearing,
        PUB_AMENITY.to_string(),
    )
        .await
        .map_err(|e| AppError::OsmError(e))?;

    let num_pubs_requested = length as usize;
    let num_pubs_found = all_elements.len();

    if num_pubs_found == 0 {
        return Err(AppError::NoPubsFound);
    }

    if num_pubs_found < num_pubs_requested {
        return Err(AppError::NotEnoughPubsFound {
            requested: length,
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

    Ok(responses)
}

pub async fn get_random_pub(lat: f64, lon: f64, radius: i16, bearing: Option<f64>) -> Result<PubResponse, AppError> {
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
    NotEnoughPubsFound { requested: i16, found: usize },
    ParseError(String),
}

impl AppError {
    pub fn to_response(&self) -> (VercelStatusCode, String) {
        match self {
            AppError::Validation(e) => (
                VercelStatusCode::BAD_REQUEST,
                format!("Input validation failed: {}", e),
            ),
            AppError::NoPubsFound => (
                VercelStatusCode::NOT_FOUND,
                "No pubs found within the specified radius.".to_string(),
            ),
            AppError::NotEnoughPubsFound  {requested, found} => (
                VercelStatusCode::NOT_FOUND,
                format!("Not enough pubs found. Request {}, but only {} found.", requested, found),
            ),
            AppError::OsmError(_e) => {
                (
                    VercelStatusCode::INTERNAL_SERVER_ERROR,
                    "An internal service error occurred.".to_string(),
                )
            },
            AppError::ParseError(e) => (
                VercelStatusCode::BAD_REQUEST,
                format!("Failed to parse query parameters: {}", e),
            ),
        }
    }
}

impl From<validator::ValidationErrors> for AppError {
    fn from(err: validator::ValidationErrors) -> Self {
        AppError::Validation(err)
    }
}