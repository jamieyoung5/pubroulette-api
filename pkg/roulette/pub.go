package roulette

import (
	"github.com/jamieyoung5/pooblet/internal/osm"
	"github.com/jamieyoung5/pooblet/internal/osm/filter"
)

type Pub struct {
	Tags      []string
	Longitude float64
	Latitude  float64
	Address   *osm.Address
	Name      filter.Names
}

var amenities = []string{"pub", "bar"}

func parsePlaceToPub(place osm.Element) (*Pub, error) {
	names := filter.FilterPlaceNameFromTags(place.Tags)
	tags := filter.FilterTags(place.Tags, filter.ValidTags)
	address, err := osm.ReverseGeocode(place.Lat, place.Lon)
	if err != nil {
		return nil, err
	}

	return &Pub{
		Tags:      tags,
		Longitude: place.Lon,
		Latitude:  place.Lat,
		Address:   address,
		Name:      names,
	}, nil
}
