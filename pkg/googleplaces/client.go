package googleplaces

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"

	"github.com/jamieyoung5/pubroulette-api/pkg/pub"
	"go.uber.org/zap"
)

const (
	apiKeyEnvVar = "GOOGLE_API_KEY"
	openNowOnly  = "OPENNOW_ONLY"

	statusNoResults = "ZERO_RESULTS"

	pubType    = "bar"
	pubKeyword = "pub"
)

var (
	placeTypes = []string{"bar", "pub"}
)

type Client struct {
	logger   *zap.Logger
	apiKey   string
	openOnly string
}

func NewClient(logger *zap.Logger) *Client {
	apiKey := os.Getenv(apiKeyEnvVar)

	if apiKey == "" {
		panic("no google places api key set")
	}

	openOnly := ""
	openOnlyEnv := os.Getenv(openNowOnly)
	if openOnlyEnv == "yes" {
		openOnly = "&opennow=true"
	}

	return &Client{
		logger:   logger,
		apiKey:   apiKey,
		openOnly: openOnly,
	}
}

// TODO: is this really the best approach in terms of extensibility when dealing with multiple data sources?

// GetRandomPub gets on single pub at random from the google places api
func (c *Client) GetRandomPub(lat, lon, radius string) (*pub.Pub, error) {
	location := formatLocation(lat, lon)
	places, err := c.getAllAvailablePlacesInRadius(location, radius, pubType, pubKeyword)
	if err != nil {
		return nil, err
	}

	randomIndex := rand.Intn(len(places))

	return places[randomIndex].toPub(), nil
}

// GetAllAvailablePubs gets all pubs available from google places api
func (c *Client) GetAllAvailablePubs(lat, lon, radius string) ([]*pub.Pub, error) {
	location := formatLocation(lat, lon)
	places, err := c.getAllAvailablePlacesInRadius(location, radius, pubType, pubKeyword)
	if err != nil {
		return nil, err
	}

	pubs := make([]*pub.Pub, len(places))
	for i, place := range places {
		pubs[i] = place.toPub()
	}

	return pubs, nil
}

func (c *Client) getAllAvailablePlacesInRadius(location, radius, placeType, placeKeyword string) ([]Result, error) {

	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/place/nearbysearch/json?location=%s&radius=%s&type=%s&key=%s&keyword=%s",
		location,
		radius,
		placeType,
		c.apiKey,
		placeKeyword,
	) + c.openOnly

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer response.Body.Close()

	var placesAPIResponse PlacesAPIResponse
	err = json.NewDecoder(response.Body).Decode(&placesAPIResponse)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	return placesAPIResponse.Results, nil

}

func formatLocation(lat, lon string) string {
	return lat + "%2C" + lon
}
