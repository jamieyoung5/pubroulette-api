use serde::{Deserialize};
use validator::Validate;
use vercel_runtime::{run, Body, Error, Request, Response as VercelResponse, StatusCode as VercelStatusCode};
use pubroulette_api::roulette::{get_random_crawl, AppError};

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

#[tokio::main]
async fn main() -> Result<(), Error> {
    run(handler).await
}

pub async fn handler(req: Request) -> Result<VercelResponse<Body>, Error> {
    let query_string = req.uri().query().unwrap_or("");
    let params: GetCrawlParams = serde_urlencoded::from_str(query_string)
        .map_err(|e| AppError::ParseError(e.to_string()))
        .and_then(|p: GetCrawlParams| {
            p.validate().map_err(AppError::from)?;
            Ok(p)
        })
        .map_err(|e| {
            let (status, message) = e.to_response();
            return Error::from(message);
        })?;

    params.validate()?;
    
    match get_random_crawl(
        params.num_pubs, 
        params.latitude, 
        params.longitude, 
        params.radius,
        params.bearing).await {
        
        Ok(pub_responses) => {
            let json_body = serde_json::to_string(&pub_responses)?;
            Ok(VercelResponse::builder()
                .status(VercelStatusCode::OK)
                .header("Content-Type", "application/json")
                .body(json_body.into())?)
        }
        Err(e) => {
            let (status, message) = e.to_response();
            let error_json = serde_json::json!({ "error": message }).to_string();
            Ok(VercelResponse::builder()
                .status(status)
                .header("Content-Type", "application/json")
                .body(error_json.into())?)
        }
    }
}
