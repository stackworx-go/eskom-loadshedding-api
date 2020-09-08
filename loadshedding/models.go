package loadshedding

import "time"

// Schedule export
type Schedule struct {
	Schedule []ScheduleDay
}

// ScheduleDay export
type ScheduleDay struct {
	Date  time.Time
	Slots []ScheduleSlot
}

type ScheduleSlotByStart []ScheduleSlot

func (a ScheduleSlotByStart) Len() int           { return len(a) }
func (a ScheduleSlotByStart) Less(i, j int) bool { return a[i].Start.Before(a[j].Start) }
func (a ScheduleSlotByStart) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type ScheduleSlotByDay []ScheduleDay

func (a ScheduleSlotByDay) Len() int           { return len(a) }
func (a ScheduleSlotByDay) Less(i, j int) bool { return a[i].Date.Before(a[j].Date) }
func (a ScheduleSlotByDay) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// ScheduleSlot export
type ScheduleSlot struct {
	Start    time.Time
	Stage    Stage
	Duration time.Duration
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
