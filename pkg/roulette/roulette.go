package roulette

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jamieyoung5/pooblet/pkg/osm"
	"github.com/jamieyoung5/pooblet/pkg/pub"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"maps"
	"math/rand"
	"slices"
)

type OverpassApi interface {
	GetAmenitiesInRadius(lat string, long string, radius string, amenity string) (osm.Places, error)
}

type ScraperFunc func(name string) (pub.Pub, error)

type Scraper struct {
	Source string
	Scrape ScraperFunc
}

type Game struct {
	redisCli    *redis.Client
	overpassApi OverpassApi
	scrapers    []Scraper
	logger      *zap.Logger
	id          uuid.UUID
}

const parsingAttempts = 10

func NewGame(logger *zap.Logger, scrapers []Scraper, overpassApi OverpassApi, redisCli *redis.Client) *Game {
	return &Game{
		redisCli:    redisCli,
		overpassApi: overpassApi,
		scrapers:    scrapers,
		logger:      logger,
		id:          uuid.New(),
	}
}

func (g *Game) Play(lat, lon string, radius string) (*pub.Pub, error) {
	places := make(osm.Places)

	g.logger.Debug(
		"Starting new game",
		zap.String("id", g.id.String()),
		zap.String("latitude", lat),
		zap.String("longitude", lon),
		zap.String("radius", radius),
	)

	// get list of all pubs in radius of lat/lon for each osm pub-related amenity
	for _, amenity := range osm.PubAmenities {
		result, err := g.overpassApi.GetAmenitiesInRadius(lat, lon, radius, amenity)
		if err != nil {
			g.logger.Error("Failed to get places list for amenity",
				zap.String("id", g.id.String()),
				zap.String("amenity", amenity),
				zap.Error(err),
			)
			return nil, ErrSearchFailure
		}
		maps.Copy(places, result)
	}

	if len(places) <= 0 {
		g.logger.Error("No places found", zap.String("id", g.id.String()))
		return nil, ErrNoPubsFound
	}

	// removed blacklisted places from our list
	err := RemoveBlacklistedOsmPlaces(g.redisCli, places)
	if err != nil {
		g.logger.Error(
			"Failed to remove blacklisted osm places",
			zap.String("id", g.id.String()),
			zap.Error(err),
		)
		// continue to attempt to find a valid pub instead of returning an error
	}

	return g.findPub(places)
}

// findPub finds a random pub from gathered places,
// with a max of 10 attempts to allow for potential data anomalies
func (g *Game) findPub(places osm.Places) (*pub.Pub, error) {
	attempts := parsingAttempts
	if len(places) <= parsingAttempts {
		attempts = len(places)
	}

	for i := range attempts {

		randomPlaceId := randomPlace(places)
		randPub, err := g.processRandomPlace(places[randomPlaceId])
		if err == nil {
			return randPub, nil
		}

		g.logger.Error(
			"Failed to process random place",
			zap.String("id", g.id.String()),
			zap.Int("place id", randomPlaceId),
			zap.String("attempt", fmt.Sprintf("%d/%d", i, parsingAttempts)),
			zap.Error(err),
		)
		err = AddToBlacklist(g.redisCli, randomPlaceId)
		if err != nil {
			g.logger.Error(
				"Failed to add blacklisted osm place",
				zap.String("id", g.id.String()),
				zap.Error(err),
			)
		}

		delete(places, i)
	}

	return nil, ErrParsingFailure
}

func (g *Game) processRandomPlace(place osm.Element) (*pub.Pub, error) {

	randPub, err := pub.OsmElementToPub(place)
	if err != nil {
		g.logger.Error("Failed to convert open street maps place to pub", zap.Int("place id", place.ID), zap.Error(err))
		return nil, err
	}

	// scrape for additional data about pub
	for _, scraper := range g.scrapers {
		result, scrapingErr := scraper.Scrape(randPub.Name.Name)

		if scrapingErr != nil {
			g.logger.Warn(
				"Failed to scrape source for additional pub data",
				zap.String("Pub name", randPub.Name.Name),
				zap.String("Source", scraper.Source),
				zap.Error(scrapingErr),
			)
		} else {
			g.logger.Debug(
				"Successfully scraped source for additional pub data",
				zap.String("Pub name", randPub.Name.Name),
				zap.String("Source", scraper.Source),
			)

			pub.Merge(randPub, result)
		}
	}

	return randPub, nil
}

func randomPlace(places osm.Places) int {
	keys := slices.Collect(maps.Keys(places))
	return keys[rand.Intn(len(keys))]
}
