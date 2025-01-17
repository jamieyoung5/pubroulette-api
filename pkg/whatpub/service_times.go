package whatpub

import (
	"fmt"
	"github.com/jamieyoung5/pooblet/pkg/pub"
	"strings"
	"time"
)

func parseOpeningTimes(lines []string) ([]pub.OpeningHour, error) {
	var result []pub.OpeningHour

	for _, line := range lines {
		oh, err := parseLine(line)
		if err != nil {
			return nil, err
		}
		result = append(result, oh)
	}
	return result, nil
}

func parseLine(line string) (pub.OpeningHour, error) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return pub.OpeningHour{}, fmt.Errorf("invalid line format: %s", line)
	}

	day := strings.TrimSpace(parts[0])
	rawTimes := strings.TrimSpace(parts[1])

	if strings.EqualFold(rawTimes, "Closed") {
		return pub.OpeningHour{
			Day:     day,
			Open24:  "",
			Close24: "",
			Closed:  true,
		}, nil
	}

	ranges := strings.SplitN(rawTimes, "-", 2)
	if len(ranges) != 2 {
		return pub.OpeningHour{}, fmt.Errorf("invalid time range for: %s", rawTimes)
	}

	openStr := strings.TrimSpace(ranges[0])
	closeStr := strings.TrimSpace(ranges[1])

	openTime, err := parseTime(openStr)
	if err != nil {
		return pub.OpeningHour{}, fmt.Errorf("error parsing open time '%s': %v", openStr, err)
	}

	closeTime, err := parseTime(closeStr)
	if err != nil {
		return pub.OpeningHour{}, fmt.Errorf("error parsing close time '%s': %v", closeStr, err)
	}

	return pub.OpeningHour{
		Day:     day,
		Open24:  openTime.Format("15:04"),
		Close24: closeTime.Format("15:04"),
		Closed:  false,
	}, nil
}

func parseTime(tStr string) (time.Time, error) {
	if strings.EqualFold(tStr, "Noon") {
		tStr = "12:00 pm"
	} else if strings.EqualFold(tStr, "Midnight") {
		tStr = "12:00 am"
	}
	tStr = strings.Replace(tStr, ".", ":", 1)

	layout := "3:04 pm"
	parsed, err := time.ParseInLocation(layout, tStr, time.Local)
	if err != nil {
		return time.Time{}, err
	}

	return parsed, nil
}
