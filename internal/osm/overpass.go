// Interacts with the overpass API

package osm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var PubAmenities = []string{"pub", "bar"}

const overpassInterpreter = "http://overpass-api.de/api/interpreter"

func Query(query string) (response []byte, err error) {
	resp, err := http.Post(overpassInterpreter, "application/x-www-form-urlencoded", bytes.NewBufferString("data="+query))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func GetAmenitiesInRadius(lat, long, radius string, amenity string) (Places, error) {
	locationRadiusParameter := fmt.Sprintf("(around:%s,%s,%s);", radius, lat, long)
	query := `[out:json];
    (
      node["amenity"="` + amenity + `"]` + locationRadiusParameter + `
      way["amenity"="` + amenity + `"]` + locationRadiusParameter + `
      relation["amenity"="` + amenity + `"]` + locationRadiusParameter + `
    );
    out body;
    >;
    out skel qt;`

	response, err := Query(query)
	if err != nil {
		return nil, err
	}

	var parsedResponse *Response
	if err = json.Unmarshal(response, &parsedResponse); err != nil {
		return nil, err
	}

	return mapPlaces(parsedResponse)
}

func mapPlaces(response *Response) (Places, error) {
	places := make(Places)
	for _, element := range response.Elements {
		places[element.ID] = element
	}

	return places, nil
}
