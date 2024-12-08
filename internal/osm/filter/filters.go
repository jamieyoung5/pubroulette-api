package filter

func filterBoolean(value string) bool {
	if value != "no" {
		return true
	}

	return false
}

func filterPresence(value string) bool {
	return true
}

func filterPrimary(value string) bool {
	if value == "primary" {
		return true
	}

	return false
}
