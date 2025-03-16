package roulette

import (
	"os"
	"reflect"

	"github.com/jamieyoung5/go-strc-yourself/pkg/sliceutil"
	"github.com/jamieyoung5/pubroulette-api/pkg/googleplaces"
	"github.com/jamieyoung5/pubroulette-api/pkg/osmoverpass"
	"github.com/jamieyoung5/pubroulette-api/pkg/osrm"
	"github.com/jamieyoung5/pubroulette-api/pkg/pub"
	"go.uber.org/zap"
)

type pubFinder interface {
	GetRandomPub(lat, lon, radius string) (*pub.Pub, error)
	GetAllAvailablePubs(lat, lon, radius string) ([]*pub.Pub, error)
}

func Play(lat, lon, radius string, logger *zap.Logger) (*pub.Pub, error) {

	finder := getFinder(logger)

	logger.Info(
		"Starting new game",
		zap.String("latitude", lat),
		zap.String("longitude", lon),
		zap.String("radius", radius),
		zap.String("finder", reflect.TypeOf(finder).String()),
	)

	pub, err := finder.GetRandomPub(lat, lon, radius)
	if err != nil {
		logger.Error(
			"Error occurred while trying to find pub",
			zap.Error(err),
		)

		return nil, err
	}

	pub.Tags = filterTags(pub.Tags)

	logger.Info(
		"Found pub",
		zap.Any("pub", pub),
	)

	return pub, nil
}

func Crawl(lat, lon, radius string, maxLength int, logger *zap.Logger) ([]*pub.Pub, error) {
	finder := getFinder(logger)

	logger.Info(
		"Starting new crawl",
		zap.String("latitude", lat),
		zap.String("longitude", lon),
		zap.String("radius", radius),
		zap.Int("max length", maxLength),
		zap.String("finder", reflect.TypeOf(finder).String()),
	)

	pubs, err := finder.GetAllAvailablePubs(lat, lon, radius)
	if err != nil {
		logger.Error(
			"Error occurred while trying to find pubs",
			zap.Error(err),
		)

		return nil, err
	}

	for i := range pubs {
		pubs[i].Tags = filterTags(pubs[i].Tags)
	}

	crawlSize := maxLength
	if maxLength > len(pubs) {
		crawlSize = len(pubs)
	}

	subset := sliceutil.RandomSubset(pubs, crawlSize)

	return osrm.GetOptimizedOrder(subset)
}

func getFinder(logger *zap.Logger) pubFinder {
	usePlaces := os.Getenv("USE_GOOGLE_PLACES")
	if usePlaces == "true" {
		return googleplaces.NewClient(logger)
	} else {
		return osmoverpass.NewOverpassClient(logger)
	}
}
