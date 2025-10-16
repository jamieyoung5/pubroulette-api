use std::collections::HashMap;
use serde::Deserialize;

#[derive(Debug)]
pub struct Element {
    pub elem_type: String,
    pub id: i64,
    pub location: Location,
    pub names: Names,
    pub tags: HashMap<String, String>,
}

#[derive(Debug)]
pub struct Location {
    pub latitude: f64,
    pub longitude: f64,
}

#[derive(Debug)]
pub struct Names {
    pub name: String,
    pub alt_name: Option<String>,
    pub old_name: Option<String>,
}

const OVERPASS_INTERPRETER: &str = "https://overpass-api.de/api/interpreter";

pub async fn get_elements_in_radius(
    latitude: f64,
    longitude: f64,
    radius: i16,
    amenity: String) -> Result<Vec<Element>, String> {

    #[derive(Deserialize)]
    struct RawLocation { lat: f64, lon: f64, }
    #[derive(Deserialize)]
    struct Center { center: RawLocation, }
    #[derive(Deserialize)]
    #[serde(untagged)]
    enum LocationData { Direct(RawLocation), Nested(Center), }
    #[derive(Deserialize)]
    struct RawElement {
        #[serde(rename = "type")]
        elem_type: String,
        id: i64,
        tags: Option<HashMap<String, String>>,
        #[serde(flatten)]
        location_data: Option<LocationData>,
    }
    #[derive(Deserialize)]
    struct RawResponse { elements: Vec<RawElement>, }

    let locator_parameter = format!("(around:{},{},{});", radius, latitude, longitude);
    let query = format!(
        r#"[out:json];
    (
      node["amenity"="{}"]{}
      way["amenity"="{}"]{}
      relation["amenity"="{}"]{}
    );
    out center;
    >;
    out skel qt;"#,
        amenity,
        locator_parameter,
        amenity,
        locator_parameter,
        amenity,
        locator_parameter
    );

    let client = reqwest::Client::new();
    let response = client
        .post(OVERPASS_INTERPRETER)
        .body(query)
        .send()
        .await.map_err(|e| e.to_string())?;

    if !response.status().is_success() {
        return Err(format!(
            "Overpass API returned a non-success status: {}",
            response.status()
        ));
    }

    let response_text = response.text().await.map_err(|e| e.to_string())?;
    let raw_response: RawResponse =
        serde_json::from_str(&response_text).map_err(|e| e.to_string())?;

    let final_elements = raw_response.elements.into_iter().filter_map(|raw| {
        let location_data = raw.location_data?;
        let tags = raw.tags?;
        let name = tags.get("name")?.clone();

        let mut alt_name = tags.get("alt_name").cloned();
        if alt_name.as_deref() == Some(&name) {
            alt_name = None;
        }

        let mut old_name = tags.get("old_name").cloned();
        if old_name.as_deref() == Some(&name) {
            old_name = None;
        }

        let raw_location = match location_data {
            LocationData::Direct(loc) => loc,
            LocationData::Nested(center) => center.center,
        };

        Some(Element {
            elem_type: raw.elem_type,
            id: raw.id,
            tags,
            names: Names { name, alt_name, old_name },
            location: Location {
                latitude: raw_location.lat,
                longitude: raw_location.lon,
            },
        })
    }).collect();

    Ok(final_elements)
}

