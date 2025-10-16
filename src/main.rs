use axum::{
    routing::get,
    Router,
};
use std::net::SocketAddr;
use crate::roulette::find_random_pub;

mod osm;
mod roulette;

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt()
        .with_env_filter(tracing_subscriber::EnvFilter::from_default_env())
        .init();

    let app = Router::new().route("/pub", get(find_random_pub));
    let addr = SocketAddr::from(([0, 0, 0, 0], 3000));

    let listener = tokio::net::TcpListener::bind(addr).await.unwrap();
    axum::serve(listener, app).await.unwrap();
}
