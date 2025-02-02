package roulette

import (
	"errors"
)

var (
	ErrNoPubsFound          = errors.New("no pubs within your radius of the provided lat/lon were found")
	ErrParsingFailure       = errors.New("failed to parse a valid pub after 3 attempts")
	ErrSearchFailure        = errors.New("failed to search for pubs")
	ErrReverseGeocodeFailed = errors.New("failed to reverse geocode")
)

func GetErrorCode(err error) int {
	switch {
	case errors.Is(err, ErrNoPubsFound):
		return 2
	default:
		return 1
	}
}
