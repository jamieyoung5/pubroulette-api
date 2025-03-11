package osmoverpass

import (
	"errors"
	"maps"
	"slices"
	"strings"

	"github.com/jamieyoung5/pubroulette-api/pkg/pub"
)

type Response struct {
	Elements []Element `json:"elements"`
}

type Element struct {
	Type  string            `json:"type"`
	ID    int               `json:"id"`
	Lat   float64           `json:"lat"`
	Lon   float64           `json:"lon"`
	Nodes []int             `json:"nodes"`
	Tags  map[string]string `json:"tags"`
}

func (e *Element) toPub() (*pub.Pub, error) {
	names, err := findNameTags(e.Tags)
	if err != nil {
		return nil, err
	}

	address, err := ReverseGeocode(e.Lat, e.Lon)
	if err != nil {
		return nil, err
	}

	if address == nil {
		return nil, errors.New("reverse geocode failed")
	}

	tags := slices.Collect(maps.Keys(e.Tags))

	return &pub.Pub{
		Tags:      tags,
		Longitude: e.Lon,
		Latitude:  e.Lat,
		Address:   normalizeAddress(address),
		Name:      names,
	}, nil
}

func normalizeAddress(addr *Address) string {
	parts := []string{}

	if addr.Road != "" {
		parts = append(parts, addr.Road)
	}
	if addr.City != "" {
		parts = append(parts, addr.City)
	}
	if addr.State != "" {
		parts = append(parts, addr.State)
	}
	if addr.Postcode != "" {
		parts = append(parts, addr.Postcode)
	}
	if addr.Country != "" {
		parts = append(parts, addr.Country)
	}

	return strings.Join(parts, ", ")
}

func findNameTags(tags map[string]string) (pub.Names, error) {
	names := pub.Names{
		Name: "unknown",
	}

	for tagName, tagValue := range tags {
		switch tagName {
		case "alt_name":
			names.AltName = tagValue
		case "old_name":
			names.OldName = tagValue
		case "name":
			names.Name = tagValue
		}
	}

	return sanitizeNames(names)
}

func sanitizeNames(names pub.Names) (pub.Names, error) {
	if names.AltName == names.Name {
		names.AltName = ""
	}

	if names.OldName == names.Name {
		names.OldName = ""
	}

	if names.Name == "unknown" {
		return names, errors.New("no valid place found")
	}

	return names, nil
}
