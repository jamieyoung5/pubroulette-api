package osrm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jamieyoung5/pubroulette-api/pkg/pub"
)

type Response struct {
	Code      string `json:"code"`
	Waypoints []struct {
		WaypointIndex int       `json:"waypoint_index"`
		Location      []float64 `json:"location"`
	} `json:"waypoints"`
}

const baseUrl = "https://router.project-osrm.org"

func GetOptimizedOrder(pubs []*pub.Pub) ([]*pub.Pub, error) {
	if len(pubs) < 2 {
		return pubs, nil
	}

	var coords []string
	for _, p := range pubs {
		coords = append(coords, fmt.Sprintf("%f,%f", p.Longitude, p.Latitude))
	}

	url := fmt.Sprintf("%s/trip/v1/walking/%s?source=first&roundtrip=false", baseUrl, strings.Join(coords, ";"))

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var osrmResp Response
	if err := json.NewDecoder(resp.Body).Decode(&osrmResp); err != nil {
		return nil, err
	}

	if osrmResp.Code != "Ok" {
		return nil, fmt.Errorf("OSRM error: %s", osrmResp.Code)
	}

	orderedPubs := make([]*pub.Pub, len(pubs))
	for i, index := range osrmResp.Waypoints {
		orderedPubs[index.WaypointIndex] = pubs[i]
	}

	return orderedPubs, nil
}
