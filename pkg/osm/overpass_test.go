package osm_test

import (
	"testing"

	"github.com/jamieyoung5/pooblet/pkg/osm"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestOverpassApi_GetAmenitiesInRadius_Success(t *testing.T) {
	logger := zap.NewExample()
	api := osm.NewOverpassApi(logger)

	lat := "51.5074"
	long := "-0.1278"
	radius := "500"
	amenity := "pub"

	places, err := api.GetAmenitiesInRadius(lat, long, radius, amenity)

	assert.NoError(t, err)
	assert.NotNil(t, places)

	assert.Greater(t, len(places), 0, "Expected to find at least one pub")

	for _, place := range places {
		assert.NotZero(t, place.ID)
	}
}

func TestOverpassApi_GetAmenitiesInRadius_NoResults(t *testing.T) {
	logger := zap.NewExample()
	api := osm.NewOverpassApi(logger)

	lat := "0.0"
	long := "0.0"
	radius := "100"
	amenity := "pub"

	places, err := api.GetAmenitiesInRadius(lat, long, radius, amenity)

	assert.NoError(t, err)

	assert.Equal(t, 0, len(places), "Expected no pubs to be found in the ocean")
}

func TestOverpassApi_GetAmenitiesInRadius_InvalidAmenity(t *testing.T) {
	logger := zap.NewExample()
	api := osm.NewOverpassApi(logger)

	lat := "51.5074"
	long := "-0.1278"
	radius := "500"
	amenity := "invalid_amenity"

	places, err := api.GetAmenitiesInRadius(lat, long, radius, amenity)

	assert.NoError(t, err)
	assert.Equal(t, 0, len(places), "Expected no results for invalid amenity type")
}

func TestOverpassApi_GetAmenitiesInRadius_InvalidCoordinates(t *testing.T) {
	logger := zap.NewExample()
	api := osm.NewOverpassApi(logger)

	lat := "91.0"
	long := "0.0"
	radius := "500"
	amenity := "pub"

	places, err := api.GetAmenitiesInRadius(lat, long, radius, amenity)

	assert.Error(t, err)
	assert.Nil(t, places)
}

func TestOverpassApi_GetAmenitiesInRadius_LargeRadius(t *testing.T) {
	logger := zap.NewExample()
	api := osm.NewOverpassApi(logger)

	lat := "51.5074"
	long := "-0.1278"
	radius := "10000"
	amenity := "pub"

	places, err := api.GetAmenitiesInRadius(lat, long, radius, amenity)

	assert.NoError(t, err)
	assert.NotNil(t, places)

	assert.Greater(t, len(places), 0, "Expected to find pubs in a 10km radius around London")
}
