// Interacts with the overpass API

package osm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
)

var PubAmenities = []string{"pub", "bar"}

// This API serves over http? big yikes
const overpassInterpreterUrl = "http://overpass-api.de/api/interpreter"

type OverpassApi struct {
	logger *zap.Logger
}

func NewOverpassApi(logger *zap.Logger) *OverpassApi {
	return &OverpassApi{
		logger: logger,
	}
}

// 'Amenities'? Is it not only getting pubs?
func (oa *OverpassApi) GetAmenitiesInRadius(lat, long, radius string, amenity string) (Places, error) {
	locationRadiusParameter := fmt.Sprintf("(around:%s,%s,%s);", radius, lat, long)
	// This query is a mess
	query := `[out:json];
    (
      node["amenity"="` + amenity + `"]` + locationRadiusParameter + `
      way["amenity"="` + amenity + `"]` + locationRadiusParameter + `
      relation["amenity"="` + amenity + `"]` + locationRadiusParameter + `
    );
    out body;
    >;
    out skel qt;`

	response, err := executeQuery(query)
	if err != nil {
		oa.logger.Error("Failed to execute query to overpass api", zap.Error(err))
		return nil, err
	}

	var parsedResponse *Response
	if err = json.Unmarshal(response, &parsedResponse); err != nil {
		oa.logger.Error("Failed to parse overpass api response", zap.Error(err))
		return nil, err
	}

	return mapPlaces(parsedResponse)
}

func executeQuery(query string) (response []byte, err error) {
	resp, err := http.Post(overpassInterpreterUrl, "application/x-www-form-urlencoded", bytes.NewBufferString("data="+query))
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

func mapPlaces(response *Response) (Places, error) {
	places := make(Places)
	for _, element := range response.Elements {
		places[element.ID] = element
	}

	return places, nil
}
