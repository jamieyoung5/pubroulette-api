use std::collections::HashMap;
use serde::{Deserialize, Serialize};

#[derive(Debug, Clone)]
pub struct Element {
    pub elem_type: String,
    pub id: i64,
    pub location: Location,
    pub names: Names,
    pub tags: HashMap<String, String>,
}

#[derive(Debug, Clone)]
pub struct Location {
    pub latitude: f64,
    pub longitude: f64,
}

#[derive(Debug, Clone)]
pub struct Names {
    pub name: String,
    pub alt_name: Option<String>,
    pub old_name: Option<String>,
}

const OVERPASS_INTERPRETER: &str = "https://overpass-api.de/api/interpreter";
const PUB_AMENITY: &str = "pub";

pub async fn get_elements_in_radius(
    latitude: f64,
    longitude: f64,
    radius: i16,
    bearing: Option<f64>,
    amenity: String) -> Result<Vec<Element>, String> {

    #[derive(Serialize, Deserialize, Debug)]
    struct RawLocation { lat: f64, lon: f64, }
    #[derive(Serialize, Deserialize, Debug)]
    struct Center { center: RawLocation, }
    #[derive(Serialize, Deserialize, Debug)]
    #[serde(untagged)]
    enum LocationData { Direct(RawLocation), Nested(Center), }
    #[derive(Serialize, Deserialize, Debug)]
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

    let elements = raw_response.elements.into_iter().filter_map(|raw| {
        
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

    if let Some(bearing_value) = bearing {
        Ok(filter_by_bearing(elements, latitude, longitude, bearing_value))
    } else {
        Ok(elements)
    }
}

fn filter_by_bearing(elements: Vec<Element>, lat: f64, lon: f64, bearing: f64) -> Vec<Element> {
    elements.into_iter().filter(|item| {
        let item_bearing = calculate_bearing(
            lat,
            lon,
            item.location.latitude,
            item.location.longitude
        );

        is_within_cone(
            item_bearing,
            bearing,
            70.0,
        )
    }).collect::<Vec<_>>()
}

fn calculate_bearing(lat1: f64, lon1: f64, lat2: f64, lon2: f64) -> f64 {
    let lat1_rad = lat1.to_radians();
    let lat2_rad = lat2.to_radians();
    let delta_lon_rad = (lon2 - lon1).to_radians();

    let y = delta_lon_rad.sin() * lat2_rad.cos();
    let x = lat1_rad.cos() * lat2_rad.sin() -
        lat1_rad.sin() * lat2_rad.cos() * delta_lon_rad.cos();

    let bearing_rad = y.atan2(x);

    // convert bearing from radians to degrees
    (bearing_rad.to_degrees() + 360.0) % 360.0
}

fn is_within_cone(bearing: f64, target_bearing: f64, cone_width: f64) -> bool {
    let half_width = cone_width / 2.0;

    let target = (target_bearing + 360.0) % 360.0;
    let bearing = (bearing + 360.0) % 360.0;

    let mut diff = (bearing - target).rem_euclid(360.0);
    if diff > 180.0 {
        diff -= 360.0;
    }

    diff.abs() <= half_width
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_calculate_bearing() {
        let precision = 1e-3;

        let bearing_north = calculate_bearing(0.0, 0.0, 1.0, 0.0);
        assert!((bearing_north - 0.0).abs() < precision);

        let bearing_east = calculate_bearing(0.0, 0.0, 0.0, 1.0);
        assert!((bearing_east - 90.0).abs() < precision);

        let bearing_south = calculate_bearing(0.0, 0.0, -1.0, 0.0);
        assert!((bearing_south - 180.0).abs() < precision);

        let bearing_west = calculate_bearing(0.0, 0.0, 0.0, -1.0);
        assert!((bearing_west - 270.0).abs() < precision);

        let bearing_ne = calculate_bearing(51.5, -0.1, 52.5, 0.9);
        assert!((bearing_ne - 31.226).abs() < precision);
    }

    #[test]
    fn test_is_within_cone() {
        let cone_width = 70.0;

        let target_90 = 90.0;
        assert!(is_within_cone(90.0, target_90, cone_width));
        assert!(is_within_cone(55.0, target_90, cone_width));
        assert!(is_within_cone(125.0, target_90, cone_width));
        assert!(!is_within_cone(54.9, target_90, cone_width));
        assert!(!is_within_cone(125.1, target_90, cone_width));

        let target_0 = 0.0;
        assert!(is_within_cone(0.0, target_0, cone_width));
        assert!(is_within_cone(360.0, target_0, cone_width));
        assert!(is_within_cone(35.0, target_0, cone_width));
        assert!(is_within_cone(325.0, target_0, cone_width));
        assert!(is_within_cone(325.0, 360.0, cone_width));
        assert!(!is_within_cone(35.1, target_0, cone_width));
        assert!(!is_within_cone(324.9, target_0, cone_width));

        let target_350 = 350.0;

        assert!(is_within_cone(340.0, target_350, cone_width));
        assert!(is_within_cone(10.0, target_350, cone_width));
        assert!(is_within_cone(25.0, target_350, cone_width));
        assert!(is_within_cone(315.0, target_350, cone_width));
        assert!(!is_within_cone(25.1, target_350, cone_width));
        assert!(!is_within_cone(314.9, target_350, cone_width));
    }

    fn create_mock_element(id: i64, lat: f64, lon: f64) -> Element {
        Element {
            elem_type: "node".to_string(),
            id,
            location: Location { latitude: lat, longitude: lon },
            names: Names {
                name: "Test Pub".to_string(),
                alt_name: None,
                old_name: None,
            },
            tags: HashMap::new(),
        }
    }

    #[test]
    fn test_filter_by_bearing() {
        let start_lat = 0.0;
        let start_lon = 0.0;

        let el_north = create_mock_element(1, 1.0, 0.0);
        let el_east = create_mock_element(2, 0.0, 1.0);
        let el_south = create_mock_element(3, -1.0, 0.0);
        let el_west = create_mock_element(4, 0.0, -1.0);

        let elements = vec![el_north, el_east, el_south, el_west];

        let filtered_north = filter_by_bearing(elements.clone(), start_lat, start_lon, 0.0);
        assert_eq!(filtered_north.len(), 1);
        assert_eq!(filtered_north[0].id, 1);

        let filtered_east = filter_by_bearing(elements.clone(), start_lat, start_lon, 90.0);
        assert_eq!(filtered_east.len(), 1);
        assert_eq!(filtered_east[0].id, 2);

        let filtered_ne = filter_by_bearing(elements.clone(), start_lat, start_lon, 45.0);

        assert_eq!(filtered_ne.len(), 0);
    }
}