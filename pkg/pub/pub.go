package pub

import (
	"github.com/jamieyoung5/pooblet/pkg/osm"
)

type Pub struct {
	Tags         []Tag
	Longitude    float64
	Latitude     float64
	Address      *osm.Address
	Name         osm.Names
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

const (
	usePlacesEnvVar = "USE_GOOGLE_PLACES"
)

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
