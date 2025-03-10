package pub

import (
	"github.com/jamieyoung5/pooblet/pkg/googleplaces"
	"github.com/jamieyoung5/pooblet/pkg/osm"
	"go.uber.org/zap"
	"os"
)

type googlePlacesClient interface {
	GetAllAvailablePubs(location string, radius string) (*googleplaces.PlacesAPIResponse, error)
}

type osmApiClient interface {
	GetAmenitiesInRadius(lat, long, radius string, amenity string) (osm.Places, error)
}

type Finder struct {
	placesClient googlePlacesClient
	osmClient    osmApiClient
	logger       *zap.Logger
}

func NewPubFinder(logger *zap.Logger) *Finder {
	usePlaces := os.Getenv(usePlacesEnvVar)

	var (
		placesClient googlePlacesClient
		osmClient    osmApiClient
	)

	if usePlaces == "true" {
		placesClient = googleplaces.NewClient(logger)
	} else {
		osmClient = osm.NewOverpassClient(logger)
	}

	return &Finder{
		placesClient: placesClient,
		osmClient:    osmClient,
		logger:       logger,
	}
}

func (f *Finder) Find() {}
