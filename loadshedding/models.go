package loadshedding

import "time"

// Schedule export
type Schedule struct {
	Schedule []ScheduleDay `json:"schedule"`
}

// ScheduleDay export
type ScheduleDay struct {
	Day   time.Time      `json:"schedule"`
	Times []ScheduleTime `json:"scheduleTime"`
}

type ByDay []ScheduleDay

func (a ByDay) Len() int           { return len(a) }
func (a ByDay) Less(i, j int) bool { return a[i].Day.Before(a[j].Day) }
func (a ByDay) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// ScheduleTime export
type ScheduleTime struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

// Suburb export
type Suburb struct {
	ID    string `json:"id"`
	Name  string `json:"text"`
	Total int    `json:"Tot"`
}

// SearchSuburb export
type SearchSuburb struct {
	ID    int `json:"Id"`
	Total int `json:"Total"`

	Municipality string `json:"MunicipalityName"`
	Province     string `json:"ProvinceName"`
	Suburb       string `json:"Name"`
}

// Municipality export
type Municipality struct {
	ID   string `json:"Value"`
	Name string `json:"Text"`
}
