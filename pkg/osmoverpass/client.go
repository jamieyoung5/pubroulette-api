package osmoverpass

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"

	"github.com/jamieyoung5/pubroulette-api/pkg/pub"
	"go.uber.org/zap"
)

const (
	overpassInterpreter = "https://overpass-api.de/api/interpreter"
	pubAmenity          = "pub"
)

type Client struct {
	logger *zap.Logger
}

func NewOverpassClient(logger *zap.Logger) *Client {
	return &Client{
		logger: logger,
	}
}

func (c *Client) GetRandomPub(lat, lon, radius string) (*pub.Pub, error) {
	amenities, err := c.getAmenitiesInRadius(lat, lon, radius, pubAmenity)
	if err != nil {
		return nil, err
	}

	randomIndex := rand.Intn(len(amenities))

	return amenities[randomIndex].toPub()
}

func (c *Client) GetAllAvailablePubs(lat, lon, radius string) ([]*pub.Pub, error) {
	places, err := c.getAmenitiesInRadius(lat, lon, radius, pubAmenity)
	if err != nil {
		return nil, err
	}

	pubs := make([]*pub.Pub, len(places))
	for i, place := range places {
		parsedPub, parsingErr := place.toPub()
		if parsingErr != nil {
			continue
		}

		pubs[i] = parsedPub
	}

	return pubs, nil
}

func (c *Client) getAmenitiesInRadius(lat, long, radius string, amenity string) ([]Element, error) {
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

	return parsedResponse.Elements, nil
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
