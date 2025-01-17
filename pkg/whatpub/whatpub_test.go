package whatpub_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jamieyoung5/pooblet/pkg/pub"
	"github.com/jamieyoung5/pooblet/pkg/whatpub"
	"github.com/stretchr/testify/assert"
)

func TestScrape_Success(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.String(), "/search/autocomplete") {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"request": "test",
				"results": [
					{"type": "pub", "match": "The Example Pub", "href": "/pubs/example-pub"}
				]
			}`))
		} else if strings.Contains(r.URL.String(), "/pubs/example-pub") {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`
				<html>
					<body>
						<section>
							<p class="pub_heading">Opening Times</p>
							<div class="opening-times-table">
								<table>
									<tr><td>Monday</td><td>12:00 pm - 11:00 pm</td></tr>
									<tr><td>Tuesday</td><td>Closed</td></tr>
								</table>
							</div>
						</section>
						<section>
							<p class="pub_heading">Facilities</p>
							<ul class="pub_features">
								<li><span>Wheelchair Accessible</span><p class="pub_feature_comment">Fully accessible</p></li>
								<li><span>Outdoor Seating</span><p class="pub_feature_comment">Available in the garden</p></li>
							</ul>
						</section>
					</body>
				</html>
			`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	whatpub.BaseUrl = mockServer.URL

	pubName := "The Example Pub"
	result, err := whatpub.Scrape(pubName)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	expectedOpeningTimes := []pub.OpeningHour{
		{Day: "Monday", Open24: "12:00", Close24: "23:00", Closed: false},
		{Day: "Tuesday", Open24: "", Close24: "", Closed: true},
	}
	expectedFacilities := []pub.Tag{
		{Name: "Wheelchair Accessible", Description: "Fully accessible"},
		{Name: "Outdoor Seating", Description: "Available in the garden"},
	}

	assert.Equal(t, expectedOpeningTimes, result.OpeningTimes)
	assert.Equal(t, expectedFacilities, result.Tags)
}

func TestScrape_NoResults(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.String(), "/search/autocomplete") {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"request": "test", "results": []}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	whatpub.BaseUrl = mockServer.URL

	pubName := "Nonexistent Pub"
	result, err := whatpub.Scrape(pubName)

	assert.Error(t, err)
	assert.Empty(t, result)
}

func TestScrape_InvalidHTML(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.String(), "/search/autocomplete") {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"request": "test",
				"results": [
					{"type": "pub", "match": "Invalid Pub", "href": "/pubs/invalid-pub"}
				]
			}`))
		} else if strings.Contains(r.URL.String(), "/pubs/invalid-pub") {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("<html><body><invalid></invalid></body></html>"))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	whatpub.BaseUrl = mockServer.URL

	pubName := "Invalid Pub"
	result, _ := whatpub.Scrape(pubName)

	assert.Empty(t, result)
}

func TestScrape_HTTPError(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	whatpub.BaseUrl = mockServer.URL

	pubName := "Error Pub"
	result, err := whatpub.Scrape(pubName)

	assert.Error(t, err)
	assert.Empty(t, result)
}
