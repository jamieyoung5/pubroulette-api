package osmoverpass

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Address struct {
	Road        string `json:"road"`
	City        string `json:"city"`
	State       string `json:"state"`
	Postcode    string `json:"postcode"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
}

// nominatim is a tool that provides a sort of 'search engine' for OpenStreetMap data
const nominatimReverseGeocodeUrl = "https://nominatim.openstreetmap.org/reverse?format=json&lat=%f&lon=%f"

// ReverseGeocode uses the nominatim api to take a longitude and latitude and provide an address correlating to that position
func ReverseGeocode(lat, lon float64) (*Address, error) {
	address := fmt.Sprintf(nominatimReverseGeocodeUrl, lat, lon)

	response, err := http.Get(address)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Address Address `json:"address"`
	}

	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result.Address, nil
}
