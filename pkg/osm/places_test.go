package osm_test

import (
	"testing"

	"github.com/jamieyoung5/pooblet/pkg/osm"
	"github.com/stretchr/testify/assert"
)

func TestElement_FindNames_Success(t *testing.T) {
	tests := []struct {
		name     string
		element  osm.Element
		expected osm.Names
	}{
		{
			name: "Valid element with all names",
			element: osm.Element{
				Tags: map[string]string{
					"name":     "name",
					"alt_name": "alt name",
					"old_name": "old name",
				},
			},
			expected: osm.Names{
				Name:    "name",
				AltName: "alt name",
				OldName: "old name",
			},
		},
		{
			name: "Valid element with only a name",
			element: osm.Element{
				Tags: map[string]string{
					"name": "name",
				},
			},
			expected: osm.Names{
				Name:    "name",
				AltName: "",
				OldName: "",
			},
		},
		{
			name: "Element with alt_name equal to name",
			element: osm.Element{
				Tags: map[string]string{
					"name":     "name",
					"alt_name": "name",
				},
			},
			expected: osm.Names{
				Name:    "name",
				AltName: "",
				OldName: "",
			},
		},
		{
			name: "Element with old_name equal to name",
			element: osm.Element{
				Tags: map[string]string{
					"name":     "name",
					"old_name": "name",
				},
			},
			expected: osm.Names{
				Name:    "name",
				AltName: "",
				OldName: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			names, err := tt.element.FindNames()
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, names)
		})
	}
}

func TestElement_FindNames_NoValidName(t *testing.T) {
	element := osm.Element{
		Tags: map[string]string{},
	}

	names, err := element.FindNames()

	assert.Error(t, err)
	assert.Equal(t, "no valid place found", err.Error())
	assert.Equal(t, "unknown", names.Name)
}

func TestElement_FindNames_InvalidCombination(t *testing.T) {
	element := osm.Element{
		Tags: map[string]string{
			"name":     "unknown",
			"alt_name": "alt name",
			"old_name": "old name",
		},
	}

	names, err := element.FindNames()

	assert.Error(t, err)
	assert.Equal(t, "unknown", names.Name)
	assert.Equal(t, "alt name", names.AltName)
	assert.Equal(t, "old name", names.OldName)
}
