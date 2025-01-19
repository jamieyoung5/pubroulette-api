package roulette_test

import (
	"errors"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"testing"

	"github.com/jamieyoung5/pooblet/pkg/osm"
	"github.com/jamieyoung5/pooblet/pkg/pub"
	"github.com/jamieyoung5/pooblet/pkg/roulette"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
)

type MockOverpassApi struct {
	mock.Mock
}

func (m *MockOverpassApi) GetAmenitiesInRadius(lat, long, radius, amenity string) (osm.Places, error) {
	args := m.Called(lat, long, radius, amenity)
	return args.Get(0).(osm.Places), args.Error(1)
}

func MockScraperFunc(name string) (pub.Pub, error) {
	return pub.Pub{}, errors.New("mock scraper error")
}

func createMockRedis(t *testing.T) (*redis.Client, *miniredis.Miniredis) {
	srv, err := miniredis.Run()
	assert.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: srv.Addr(),
	})

	return client, srv
}

func TestNewGame(t *testing.T) {
	logger := zaptest.NewLogger(t)
	scrapers := []roulette.Scraper{
		{Source: "test_source", Scrape: MockScraperFunc},
	}
	overpassApi := &MockOverpassApi{}

	client, srv := createMockRedis(t)
	defer srv.Close()

	game := roulette.NewGame(logger, scrapers, overpassApi, client)

	assert.NotNil(t, game)
}

func TestGame_Play_NoPlacesFound(t *testing.T) {
	logger := zaptest.NewLogger(t)
	scrapers := []roulette.Scraper{
		{Source: "test_source", Scrape: MockScraperFunc},
	}
	overpassApi := &MockOverpassApi{}
	overpassApi.On("GetAmenitiesInRadius", "51.5074", "-0.1278", "500", mock.Anything).Return(osm.Places{}, nil)

	client, srv := createMockRedis(t)
	defer srv.Close()

	game := roulette.NewGame(logger, scrapers, overpassApi, client)

	_, err := game.Play("51.5074", "-0.1278", "500")
	assert.Error(t, err)
	assert.Equal(t, roulette.ErrNoPubsFound, err)
	overpassApi.AssertCalled(t, "GetAmenitiesInRadius", "51.5074", "-0.1278", "500", "pub")
}

func TestGame_Play_Success(t *testing.T) {
	logger := zaptest.NewLogger(t)
	scrapers := []roulette.Scraper{
		{Source: "test_source", Scrape: func(name string) (pub.Pub, error) {
			return pub.Pub{Name: osm.Names{Name: name}}, nil
		}},
	}
	overpassApi := &MockOverpassApi{}
	places := osm.Places{
		1: {ID: 1, Tags: map[string]string{"name": "Test Pub"}},
	}
	overpassApi.On("GetAmenitiesInRadius", "51.5074", "-0.1278", "500", mock.Anything).Return(places, nil)

	client, srv := createMockRedis(t)
	defer srv.Close()

	game := roulette.NewGame(logger, scrapers, overpassApi, client)

	result, err := game.Play("51.5074", "-0.1278", "500")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Pub", result.Name.Name)
	overpassApi.AssertCalled(t, "GetAmenitiesInRadius", "51.5074", "-0.1278", "500", "pub")
}

func TestGame_Play_ScraperErrorHandling(t *testing.T) {
	logger := zaptest.NewLogger(t)
	scrapers := []roulette.Scraper{
		{Source: "test_source", Scrape: MockScraperFunc},
	}
	overpassApi := &MockOverpassApi{}
	places := osm.Places{
		1: {ID: 1, Tags: map[string]string{"name": "Scraper Test Pub"}},
	}
	overpassApi.On("GetAmenitiesInRadius", "51.5074", "-0.1278", "500", mock.Anything).Return(places, nil)

	client, srv := createMockRedis(t)
	defer srv.Close()

	game := roulette.NewGame(logger, scrapers, overpassApi, client)

	result, err := game.Play("51.5074", "-0.1278", "500")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Scraper Test Pub", result.Name.Name)
}

func TestGame_Play_MultipleAmenities(t *testing.T) {
	logger := zaptest.NewLogger(t)
	scrapers := []roulette.Scraper{}
	overpassApi := &MockOverpassApi{}
	places1 := osm.Places{
		1: {ID: 1, Tags: map[string]string{"name": "Amenity 1 Pub"}},
	}
	places2 := osm.Places{
		2: {ID: 2, Tags: map[string]string{"name": "Amenity 2 Pub"}},
	}
	overpassApi.On("GetAmenitiesInRadius", "51.5074", "-0.1278", "500", "pub").Return(places1, nil)
	overpassApi.On("GetAmenitiesInRadius", "51.5074", "-0.1278", "500", "bar").Return(places2, nil)

	client, srv := createMockRedis(t)
	defer srv.Close()

	game := roulette.NewGame(logger, scrapers, overpassApi, client)

	result, err := game.Play("51.5074", "-0.1278", "500")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Contains(t, []string{"Amenity 1 Pub", "Amenity 2 Pub"}, result.Name.Name) // Ensure one of the pubs is returned
}
