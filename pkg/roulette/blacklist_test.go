package roulette_test

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/jamieyoung5/pooblet/pkg/osm"
	"github.com/jamieyoung5/pooblet/pkg/roulette"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRemoveBlacklistedOsmPlaces(t *testing.T) {
	srv, err := miniredis.Run()
	assert.NoError(t, err)
	defer srv.Close()

	client := redis.NewClient(&redis.Options{
		Addr: srv.Addr(),
	})

	places := osm.Places{
		1: {ID: 1},
		2: {ID: 2},
		3: {ID: 3},
	}

	err = roulette.AddToBlacklist(client, 2)
	assert.NoError(t, err)

	err = roulette.RemoveBlacklistedOsmPlaces(client, places)
	assert.NoError(t, err)

	assert.Len(t, places, 2)
	_, stillPresent := places[2]
	assert.False(t, stillPresent, "expected place #2 to be removed")
}

func TestIsBlacklisted(t *testing.T) {
	srv, err := miniredis.Run()
	assert.NoError(t, err)
	defer srv.Close()

	client := redis.NewClient(&redis.Options{
		Addr: srv.Addr(),
	})

	blacklisted, err := roulette.IsBlacklisted(client, 42)
	assert.NoError(t, err)
	assert.False(t, blacklisted, "expected #42 not to be blacklisted yet")

	err = roulette.AddToBlacklist(client, 42)
	assert.NoError(t, err)

	blacklisted, err = roulette.IsBlacklisted(client, 42)
	assert.NoError(t, err)
	assert.True(t, blacklisted, "expected #42 to be blacklisted now")
}
