package roulette

import (
	"context"
	"github.com/jamieyoung5/pooblet/pkg/osm"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func RemoveBlacklistedOsmPlaces(cli *redis.Client, places osm.Places) error {
	for _, place := range places {
		blacklisted, err := IsBlacklisted(cli, place.ID)
		if err != nil {
			return err
		}

		if blacklisted {
			delete(places, place.ID)
		}
	}

	return nil
}

func IsBlacklisted(cli *redis.Client, id int) (bool, error) {
	inSet, err := cli.SIsMember(ctx, "blacklist", id).Result()
	if err != nil {
		return false, err
	}
	return inSet, nil
}

func AddToBlacklist(cli *redis.Client, id int) error {
	return cli.SAdd(ctx, "blacklist", id).Err()
}
