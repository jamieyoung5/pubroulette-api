package osm_test

import (
	"testing"

	"github.com/jamieyoung5/pooblet/pkg/osm"
	"github.com/stretchr/testify/assert"
)

func TestFilterTags(t *testing.T) {
	tests := []struct {
		name         string
		inputTags    map[string]string
		expectedTags map[string]string
	}{
		{
			name: "Valid boolean tags",
			inputTags: map[string]string{
				"diet:vegetarian": "yes",
				"brewery":         "yes",
				"wheelchair":      "no",
			},
			expectedTags: map[string]string{
				"diet:vegetarian": "Vegetarian Options",
				"brewery":         "Brewery",
			},
		},
		{
			name: "Invalid tags are removed",
			inputTags: map[string]string{
				"random_tag":      "value",
				"unknown_feature": "true",
			},
			expectedTags: map[string]string{},
		},
		{
			name: "Presence filter always passes",
			inputTags: map[string]string{
				"real ale": "any value",
			},
			expectedTags: map[string]string{
				"real ale": "Real Ale",
			},
		},
		{
			name: "Primary filter with valid value",
			inputTags: map[string]string{
				"lgbtq": "primary",
			},
			expectedTags: map[string]string{
				"lgbtq": "LGBTQ+",
			},
		},
		{
			name: "Primary filter with invalid value",
			inputTags: map[string]string{
				"lgbtq": "secondary",
			},
			expectedTags: map[string]string{},
		},
		{
			name: "Mixed valid and invalid tags",
			inputTags: map[string]string{
				"diet:vegan":         "yes",
				"outdoor_seating":    "no",
				"toilets:wheelchair": "yes",
				"unknown_tag":        "true",
			},
			expectedTags: map[string]string{
				"diet:vegan":         "Vegan Options",
				"toilets:wheelchair": "Wheelchair Accessible Toilet",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			osm.FilterTags(tt.inputTags)
			assert.Equal(t, tt.expectedTags, tt.inputTags)
		})
	}
}
