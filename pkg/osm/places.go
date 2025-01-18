package osm

import "errors"

type Places map[int]Element

type Response struct {
	Elements []Element `json:"elements"`
}

type Element struct {
	Type  string            `json:"type"`
	ID    int               `json:"id"`
	Lat   float64           `json:"lat"`
	Lon   float64           `json:"lon"`
	Nodes []int             `json:"nodes"`
	Tags  map[string]string `json:"tags"`
}

type Names struct {
	AltName string
	Name    string
	OldName string
}

func (e *Element) FindNames() (Names, error) {
	names := Names{
		Name: "unknown", // Setting this as the default value is kinda an anti-pattern. What if the pub is actually called "unknown"?
	}

	for tagName, tagValue := range e.Tags {
		switch tagName {
		case "alt_name":
			names.AltName = tagValue
		case "old_name":
			names.OldName = tagValue
		case "name":
			names.Name = tagValue
		}
	}

	return sanitizeNames(names)
}

func sanitizeNames(names Names) (Names, error) {
	if names.AltName == names.Name {
		names.AltName = ""
	}

	if names.OldName == names.Name {
		names.OldName = ""
	}

	if names.Name == "unknown" {
		return names, errors.New("no valid place found")
	}

	return names, nil
}
