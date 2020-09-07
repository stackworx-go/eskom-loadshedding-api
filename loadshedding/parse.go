package loadshedding

import (
	"strings"
	"time"
)

func parseDate(value string, location *time.Location) (time.Time, error) {
	return time.ParseInLocation("Mon, 02 Jan", value, location)
}

func parseTime(value string, location *time.Location) (time.Time, error) {
	return time.ParseInLocation("15:04", value, location)
}

func parseHTMLTime(year int, dateOfMonth, duration string, location *time.Location) (*time.Time, *time.Time, error) {
	times := strings.Split(duration, " - ")

	date, err := parseDate(dateOfMonth, location)

	if err != nil {
		return nil, nil, err
	}

	startTime, err := parseTime(times[0], location)

	if err != nil {
		return nil, nil, err
	}

	endTime, err := parseTime(times[1], location)

	if err != nil {
		return nil, nil, err
	}

	return combineDateAndTime(year, date, startTime), combineDateAndTime(year, date, endTime), nil
}

func combineDateAndTime(year int, date time.Time, time_ time.Time) *time.Time {
	t := time.Date(
		year,
		date.Month(),
		date.Day(),
		time_.Hour(),
		time_.Minute(),
		time_.Second(),
		time_.Nanosecond(),
		time_.Location(),
	)

	return &t
}
