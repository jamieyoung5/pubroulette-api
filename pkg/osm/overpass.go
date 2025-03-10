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

const overpassInterpreter = "https://overpass-api.de/api/interpreter"

type Client struct {
	logger *zap.Logger
}

func NewOverpassClient(logger *zap.Logger) *Client {
	return &Client{
		logger: logger,
	}
}

func (c *Client) GetAmenitiesInRadius(lat, long, radius string, amenity string) (Places, error) {
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

	response, err := executeQuery(query)
	if err != nil {
		c.logger.Error("Failed to execute query to overpass api", zap.Error(err))
		return nil, err
	}

	var parsedResponse *Response
	if err = json.Unmarshal(response, &parsedResponse); err != nil {
		c.logger.Error("Failed to parse overpass api response", zap.Error(err))
		return nil, err
	}

	return mapPlaces(parsedResponse)
}

func executeQuery(query string) (response []byte, err error) {
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

func mapPlaces(response *Response) (Places, error) {
	places := make(Places)
	for _, element := range response.Elements {
		places[element.ID] = element
	}

	return places, nil
}
