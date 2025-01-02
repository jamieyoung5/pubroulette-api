//Responsible for validating different inputs that can be provided to the api.

package verification

import (
	"errors"
	"fmt"
	"strconv"
)

const maxRadius = 2000

func VerifyLocation(long float64, lat float64) (latitude, longitude string, err error) {
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

func VerifyRadius(radius int) (string, error) {
	if radius < 0 || radius > maxRadius {
		return "", errors.New("invalid radius")
	}
	return fmt.Sprintf("%d", radius), nil
}
