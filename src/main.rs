use std::env;
use axum::{
    routing::get,
    Router,
};
use std::net::SocketAddr;
use reqwest::header::ACCEPT;
use reqwest::Method;
use serde::Deserialize;
use tower_http::cors::{AllowOrigin, Any, CorsLayer};
use crate::roulette::{find_random_crawl, find_random_pub};
use http::HeaderValue;

use tower::ServiceBuilder;
use tower_governor::governor::GovernorConfigBuilder;
use tower_governor::GovernorLayer;

mod osm;
mod roulette;

const PORT_ENV_NAME: &str = "PORT";
const OSM_HOST_ENV_NAME: &str = "OSM_HOST";
const ALLOW_ORIGINS_ENV_NAME: &str = "ALLOW_ORIGINS";

struct ApiConfig {
    port: String,
    osm_host: String,
    allow_origins: String,
}

#[derive(Clone)]
struct AppState {
    http_client: reqwest::Client,
    osm_host: String,
}

#[tokio::main]
async fn main() {
    let cfg = bootstrap();

    let origins = match cfg.allow_origins.as_str() {
        "*" => AllowOrigin::any(),
        _ => {
            let parsed_origins = cfg.allow_origins
                .split(',')
                .map(|origin| {
                    origin.parse::<HeaderValue>()
                        .expect("failed to parse origin into HeaderValue")
                })
                .collect::<Vec<HeaderValue>>();
            AllowOrigin::list(parsed_origins)
        }
    };

    tracing_subscriber::fmt()
        .with_env_filter(tracing_subscriber::EnvFilter::from_default_env())
        .init();

    let cors = CorsLayer::new()
        .allow_origin(origins)
        .allow_methods(Method::GET)
        .allow_headers([ACCEPT]);

    let app_state = AppState {
        http_client: reqwest::Client::new(),
        osm_host: cfg.osm_host,
    };

    let governor_config = Box::new(
        GovernorConfigBuilder::default().per_second(120)
            .per_minute(30)
            .burst_size(5)
            .finish()
            .expect("Failed to build governor config"),
    );

    let app = Router::new()
        .route("/pub", get(find_random_pub))
        .route("/crawl", get(find_random_crawl))
        .route("/health", get(|| async { "ok" }))
        .with_state(app_state)
        .layer(
            ServiceBuilder::new()
                .layer(cors)
                .layer(GovernorLayer {
                    config: Box::leak(governor_config),
                })
        );

    let addr = SocketAddr::from(([0, 0, 0, 0], cfg.port.parse().unwrap()));

    let listener = tokio::net::TcpListener::bind(addr).await.unwrap();
    axum::serve(listener, app).await.unwrap();
}

fn bootstrap() -> ApiConfig {

    ApiConfig{
        port: env::var(PORT_ENV_NAME).expect("failed to read PORT"),
        osm_host: env::var(OSM_HOST_ENV_NAME).expect("failed to read OSM_HOST"),
        allow_origins: env::var(ALLOW_ORIGINS_ENV_NAME).expect("failed to read ALLOW_ORIGINS"),
    }

}