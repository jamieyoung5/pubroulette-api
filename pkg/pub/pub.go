package pub

import (
	osm2 "github.com/jamieyoung5/pooblet/pkg/osm"
)

type Pub struct {
	Tags         []Tag
	Longitude    float64
	Latitude     float64
	Address      *osm2.Address
	Name         osm2.Names
	OpeningTimes []OpeningHour
}

type Tag struct {
	Name        string
	Description string
}

type OpeningHour struct {
	Day     string
	Open24  string
	Close24 string
	Closed  bool
}

func Merge(subject *Pub, merging Pub) {
	if merging.Tags != nil {
		subject.Tags = merging.Tags
	}

	if merging.OpeningTimes != nil {
		subject.OpeningTimes = merging.OpeningTimes
	}

	if merging.Name.OldName != "" {
		subject.Name.OldName = merging.Name.OldName
	}

	if merging.Name.AltName != "" {
		subject.Name.AltName = merging.Name.AltName
	}
}

func OsmElementToPub(element osm2.Element) (*Pub, error) {
	names, err := element.FindNames()
	if err != nil {
		return nil, err
	}

	address, err := osm2.ReverseGeocode(element.Lat, element.Lon)
	if err != nil {
		return nil, err
	}

	osm2.FilterTags(element.Tags)

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
