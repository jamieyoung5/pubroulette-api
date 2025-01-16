package whatpub

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/jamieyoung5/pooblet/internal/pub"
	"strings"
)

func scrapeOpeningTimes(doc *goquery.Document) []string {
	var openingTimes []string

	section := findSectionByHeading(doc, "Opening Times")
	if section == nil {
		return openingTimes
	}

	section.Find(".opening-times-table table tr").Each(func(_ int, tr *goquery.Selection) {
		cols := tr.Find("td")
		if cols.Length() == 2 {
			day := strings.TrimSpace(cols.Eq(0).Text())
			times := strings.TrimSpace(cols.Eq(1).Text())
			if day != "" && times != "" {
				openingTimes = append(openingTimes, fmt.Sprintf("%s: %s", day, times))
			}
		}
	})

	return openingTimes
}

func scrapeFacilities(doc *goquery.Document) []pub.Tag {
	var facilities []pub.Tag

	section := findSectionByHeading(doc, "Facilities")
	if section == nil {
		return facilities
	}

	section.Find("ul.pub_features li").Each(func(_ int, li *goquery.Selection) {
		name := strings.TrimSpace(li.Find("span").Text())
		description := strings.TrimSpace(li.Find("p.pub_feature_comment").Text())

		if name != "" {
			facilities = append(facilities, pub.Tag{
				Name:        name,
				Description: description,
			})
		}
	})

	return facilities
}

func findSectionByHeading(doc *goquery.Document, heading string) *goquery.Selection {
	var matched *goquery.Selection

	doc.Find("section").EachWithBreak(func(_ int, section *goquery.Selection) bool {
		text := strings.TrimSpace(section.Find("p.pub_heading").Text())
		if text == heading {
			matched = section
			return false // stop iteration
		}
		return true
	})

	return matched
}
