package filter

func FilterTags(tags map[string]string, whitelist []Tag) (filteredTags []string) {

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

func FilterPlaceNameFromTags(tags map[string]string) Names {
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

	return names
}
