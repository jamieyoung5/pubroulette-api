package whatpub

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/jamieyoung5/pooblet/pkg/pub"
	"net/http"
	"net/url"
)

type PubInfo struct {
	OpeningTimes []pub.OpeningHour
	Features     []string
	Facilities   []pub.Tag
}

type Facility struct {
	Name    string
	Comment string
}

type searchApiResponse struct {
	Request string `json:"request"`
	Results []struct {
		Type  string `json:"type"`
		Match string `json:"match"`
		Href  string `json:"href"`
	} `json:"results"`
}

// Why cant this be a const?
var BaseUrl = "https://whatpub.com"

const searchApi = "/search/autocomplete?q=%s&features=&limit=10&AdditionalServices=false&home=1"

func Scrape(pubName string) (pub.Pub, error) {
	var result pub.Pub

	pubUrl, err := findPubUrl(pubName)
	if err != nil {
		return result, err
	}

	resp, err := http.Get(pubUrl)
	if err != nil {
		return result, fmt.Errorf("failed to GET pub page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return result, fmt.Errorf("failed to parse HTML: %w", err)
	}

	scrapedOpeningTimes := scrapeOpeningTimes(doc)
	openingTimes, err := parseOpeningTimes(scrapedOpeningTimes)
	if err != nil {
		return result, fmt.Errorf("failed to parse opening times: %w", err)
	}
	result.OpeningTimes = openingTimes
	result.Tags = scrapeFacilities(doc)

	return result, nil
}

func findPubUrl(name string) (string, error) {
	encodedName := url.QueryEscape(name)
	searchRequestUrl := BaseUrl + fmt.Sprintf(searchApi, encodedName)

	resp, err := http.Get(searchRequestUrl)
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
	}

	var searchResp searchApiResponse
	if err = json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return "", err
	}

	if len(searchResp.Results) != 1 {
		return "", errors.New(fmt.Sprintf("unexpected number of results: %d", len(searchResp.Results)))
	}

	href := searchResp.Results[0].Href

	return BaseUrl + href, nil
}
