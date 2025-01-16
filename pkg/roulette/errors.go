package roulette

import "errors"

var (
	ErrNoPubsFound    = errors.New("no pubs within your radius of the provided lat/lon were found")
	ErrParsingFailure = errors.New("failed to parse a valid pub after 3 attempts")
)
