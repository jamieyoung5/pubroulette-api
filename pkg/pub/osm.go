package pub

import (
	"errors"
	"github.com/jamieyoung5/pooblet/pkg/osm"
)

func OsmElementToPub(element osm.Element) (*Pub, error) {
	names, err := element.FindNames()
	if err != nil {
		return nil, err
	}

	address, err := osm.ReverseGeocode(element.Lat, element.Lon)
	if err != nil {
		return nil, err
	}

	if address == nil {
		return nil, errors.New("reverse geocode failed")
	}

	osm.FilterTags(element.Tags)

	return &Pub{
		Tags:      convertOsmTags(element.Tags),
		Longitude: element.Lon,
		Latitude:  element.Lat,
		Address:   address,
		Name:      names,
	}, nil
}

func convertOsmTags(tags map[string]string) []Tag {
	convertedTags := make([]Tag, 0, len(tags))

	for _, value := range tags {
		convertedTags = append(convertedTags, Tag{Name: value})
	}

	return convertedTags
}
