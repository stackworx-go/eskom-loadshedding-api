package loadshedding

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

func parseDate(value string, location *time.Location) (time.Time, error) {
	return time.ParseInLocation("Mon, 02 Jan", value, location)
}

func parseTime(value string, location *time.Location) (time.Time, error) {
	return time.ParseInLocation("15:04", value, location)
}

var slotRe = regexp.MustCompile(`(\d{2}:\d{2} - \d{2}:\d{2})`)

type Slot struct {
	Start time.Time
	End   time.Time
}

func parseHTMLTime(year int, dateOfMonth, slotsRaw string, location *time.Location) ([]Slot, error) {
	matches := slotRe.FindAllString(slotsRaw, -1)

	if matches == nil {
		return nil, fmt.Errorf("invalid slot: %s", slotsRaw)
	}

	var slots []Slot

	for _, match := range matches {
		times := strings.Split(match, " - ")

		date, err := parseDate(dateOfMonth, location)

		if err != nil {
			return nil, err
		}

		startTime, err := parseTime(times[0], location)

		if err != nil {
			return nil, err
		}

		endTime, err := parseTime(times[1], location)

		if err != nil {
			return nil, err
		}

		start := combineDateAndTime(year, date, startTime)
		end := combineDateAndTime(year, date, endTime)

		// Shift end to next day
		if end.Before(start) {
			end = end.Add(time.Hour * 24 * 1)
		}

		slots = append(slots, Slot{
			Start: start,
			End:   end,
		})
	}

	return slots, nil
}

func combineDateAndTime(year int, date time.Time, time_ time.Time) time.Time {
	return time.Date(
		year,
		date.Month(),
		date.Day(),
		time_.Hour(),
		time_.Minute(),
		time_.Second(),
		time_.Nanosecond(),
		time_.Location(),
	)
}
