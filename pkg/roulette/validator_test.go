package roulette_test

import (
	"errors"
	"testing"

	"github.com/jamieyoung5/pooblet/pkg/roulette"
	"github.com/stretchr/testify/assert"
)

func TestValidateLocation(t *testing.T) {
	tests := []struct {
		name      string
		long      float64
		lat       float64
		wantLat   string
		wantLong  string
		wantError error
	}{
		{
			name:      "Valid location",
			long:      0.0,
			lat:       51.5074,
			wantLat:   "51.507400",
			wantLong:  "0.000000",
			wantError: nil,
		},
		{
			name:      "Invalid longitude (too small)",
			long:      -181.0,
			lat:       51.5074,
			wantLat:   "",
			wantLong:  "",
			wantError: errors.New("invalid longitude"),
		},
		{
			name:      "Invalid longitude (too large)",
			long:      181.0,
			lat:       51.5074,
			wantLat:   "",
			wantLong:  "",
			wantError: errors.New("invalid longitude"),
		},
		{
			name:      "Invalid latitude (too small)",
			long:      0.0,
			lat:       -91.0,
			wantLat:   "",
			wantLong:  "",
			wantError: errors.New("invalid latitude"),
		},
		{
			name:      "Invalid latitude (too large)",
			long:      0.0,
			lat:       91.0,
			wantLat:   "",
			wantLong:  "",
			wantError: errors.New("invalid latitude"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			latitude, longitude, err := roulette.ValidateLocation(tt.long, tt.lat)

			assert.Equal(t, tt.wantLat, latitude)
			assert.Equal(t, tt.wantLong, longitude)
			if tt.wantError != nil {
				assert.EqualError(t, err, tt.wantError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateRadius(t *testing.T) {
	tests := []struct {
		name      string
		radius    int
		want      string
		wantError error
	}{
		{
			name:      "Valid radius",
			radius:    500,
			want:      "500",
			wantError: nil,
		},
		{
			name:      "Radius at maximum limit",
			radius:    2000,
			want:      "2000",
			wantError: nil,
		},
		{
			name:      "Radius below minimum",
			radius:    -1,
			want:      "",
			wantError: errors.New("invalid radius"),
		},
		{
			name:      "Radius above maximum",
			radius:    2001,
			want:      "",
			wantError: errors.New("invalid radius"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := roulette.ValidateRadius(tt.radius)

			assert.Equal(t, tt.want, result)
			if tt.wantError != nil {
				assert.EqualError(t, err, tt.wantError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
