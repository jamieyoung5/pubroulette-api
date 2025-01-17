package osm_test

import (
	"testing"

	"github.com/jamieyoung5/pooblet/pkg/osm"
	"github.com/stretchr/testify/assert"
)

func TestReverseGeocode_RealAPI_Success(t *testing.T) {
	lat := 48.8584
	lon := 2.2945

	address, err := osm.ReverseGeocode(lat, lon)

	assert.NoError(t, err)
	assert.NotNil(t, address)

	assert.NotEmpty(t, address.Road)
	assert.Equal(t, "France", address.Country)
}

func TestReverseGeocode_RealAPI_BoundaryConditions(t *testing.T) {
	tests := []struct {
		name      string
		lat       float64
		lon       float64
		expectErr bool
	}{
		{"North Pole", 90.0, 135.0, false},
		{"South Pole", -90.0, 45.0, false},
		{"International Date Line (West)", 0.0, -180.0, false},
		{"International Date Line (East)", 0.0, 180.0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			address, err := osm.ReverseGeocode(tt.lat, tt.lon)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, address)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, address)
			}
		})
	}
}

func TestReverseGeocode_RealAPI_InvalidCoordinates(t *testing.T) {
	tests := []struct {
		name string
		lat  float64
		lon  float64
	}{
		{"Invalid Latitude (too high)", 91.0, 0.0},
		{"Invalid Latitude (too low)", -91.0, 0.0},
		{"Invalid Longitude (too high)", 0.0, 181.0},
		{"Invalid Longitude (too low)", 0.0, -181.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			address, err := osm.ReverseGeocode(tt.lat, tt.lon)

			assert.NoError(t, err)
			assert.Empty(t, address)
		})
	}
}

func TestReverseGeocode_RealAPI_MiddleOfNowhere(t *testing.T) {
	lat := -25.0
	lon := -140.0

	address, err := osm.ReverseGeocode(lat, lon)

	assert.NoError(t, err)
	assert.NotNil(t, address)

	assert.Empty(t, address.Road)
}
