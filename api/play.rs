
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

#[derive(Debug, Serialize)]
pub struct PubResponse {
    name: String,
    latitude: f64,
    longitude: f64,
}

#[tokio::main]
async fn main() -> Result<(), Error> {
    let app = Router::new().route("/", get(find_random_pub));

    Ok(run(app).await?)
}