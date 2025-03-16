package roulette

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	maxRadius = 2000 // in meters
	maxLength = 5    // max pub crawl length
)

func ValidateLocation(long float64, lat float64) (latitude, longitude string, err error) {
	if long < -180 || long > 180 {
		return "", "", errors.New("invalid longitude")
	}
	longitude = strconv.FormatFloat(long, 'f', 6, 64)

	if lat < -90 || lat > 90 {
		return "", "", errors.New("invalid latitude")
	}
	latitude = strconv.FormatFloat(lat, 'f', 6, 64)

	return latitude, longitude, nil
}

func ValidateRadius(radius int) (string, error) {
	if radius < 0 || radius > maxRadius {
		return "", errors.New("invalid radius")
	}
	return fmt.Sprintf("%d", radius), nil
}

func ValidateLength(length int) (int, error) {
	if length < 1 || length > maxLength {
		return -1, errors.New("invalid length")
	}
	return length, nil
}
