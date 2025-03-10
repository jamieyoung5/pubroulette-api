package googleplaces

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"os"
)

const (
	apiKeyEnvVar = "GOOGLE_API_KEY"

	statusNoResults = "ZERO_RESULTS"
)

var (
	placeTypes = []string{"bar", "pub"}
)

type Client struct {
	logger *zap.Logger
	apiKey string
}

func NewClient(logger *zap.Logger) *Client {
	apiKey := os.Getenv(apiKeyEnvVar)

	if apiKey == "" {
		panic("no google places api key set")
	}

	return &Client{logger: logger, apiKey: apiKey}
}

func (c *Client) GetAllAvailablePubs(location string, radius string) (*PlacesAPIResponse, error) {
	googlePlacesApiKey := os.Getenv(apiKeyEnvVar)
	if googlePlacesApiKey == "" {
		fmt.Println("api key not found")
	}

	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/place/nearbysearch/json?location=%s&radius=%s&type=%s&key=%s&keyword=%s&opennow=true",
		location,
		radius,
		placeTypes[0],
		googlePlacesApiKey,
		placeTypes[1])

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

	return &placesAPIResponse, nil

}
