package filter

import "errors"

func Tags(tags map[string]string, whitelist []Tag) (filteredTags []string) {

	for tagName, tagValue := range tags {
		for _, validTag := range whitelist {
			if tagName == validTag.Name {
				if validTag.Filter(tagValue) {
					filteredTags = append(filteredTags, validTag.Alias)
				}
			}
		}
	}

	return filteredTags
}

func PlaceNameFromTags(tags map[string]string) (Names, error) {
	names := Names{
		Name: "unknown",
	}

	for tagName, tagValue := range tags {
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
