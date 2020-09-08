package loadshedding

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
)

// Client export
type Client struct {
	host     string
	client   *resty.Client
	location *time.Location
}

// NewClient export
func NewClient(location *time.Location, debug bool) *Client {
	client := resty.New()

	if debug {
		client.EnableTrace()
	}

	client.SetRetryCount(2)

	// API is flaky
	client.SetRetryWaitTime(1 * time.Second)

	return &Client{
		client:   client,
		location: location,
		host:     "https://loadshedding.eskom.co.za/LoadShedding",
	}
}

func makeTimestamp() string {
	return fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))
}

func (c *Client) createRequest() *resty.Request {
	return c.client.
		R().
		SetHeader("Accept", "application/json").
		SetQueryParams(map[string]string{
			"_": makeTimestamp(),
		})
}

// GetStatus export
func (c *Client) GetStatus() (Stage, error) {
	resp, err := c.createRequest().
		Get(c.host + "/GetStatus")

	if err != nil {
		return StageUnknown, err
	}

	stage := ConvertStage(resp.String())

	return stage, nil
}

// GetMunicipalities export
func (c *Client) GetMunicipalities(province Province) ([]Municipality, error) {
	var results []Municipality

	_, err := c.createRequest().
		SetQueryParam("Id", fmt.Sprintf("%d", province)).
		SetResult(&results).
		Get(c.host + "/GetMunicipalities")

	if err != nil {
		return nil, err
	}

	return results, nil
}

// GetMunicipalitySuburbsRequest export
type GetMunicipalitySuburbsRequest struct {
	MunicipalityID string
	Search         string
	PageSize       int
}

// GetMunicipalitySuburbsResponse export
type GetMunicipalitySuburbsResponse struct {
	Total   int
	Results []Suburb
}

// GetMunicipalitySuburbs export
func (c *Client) GetMunicipalitySuburbs(req GetMunicipalitySuburbsRequest) ([]Suburb, error) {
	var suburbs []Suburb

	if req.PageSize == 0 {
		req.PageSize = 1000
	}

	var results GetMunicipalitySuburbsResponse

	// Get Total Results
	_, err := c.createRequest().
		SetQueryParams(map[string]string{
			"pageSize":   strconv.Itoa(0),
			"pageNum":    strconv.Itoa(1),
			"searchTerm": req.Search,
			"id":         req.MunicipalityID,
		}).
		SetResult(&results).
		Get(c.host + "/GetSurburbData")

	if err != nil {
		return nil, err
	}

	pages := (results.Total / req.PageSize) + 1

	for i := 1; i <= pages; i++ {
		_, err := c.createRequest().
			SetQueryParams(map[string]string{
				"pageSize":   strconv.Itoa(req.PageSize),
				"pageNum":    strconv.Itoa(i),
				"searchTerm": req.Search,
				"id":         req.MunicipalityID,
			}).
			SetResult(&results).
			Get(c.host + "/GetSurburbData")

		if err != nil {
			return nil, err
		}

		suburbs = append(suburbs, results.Results...)
	}

	return suburbs, nil
}

// SearchSuburbsRequest export
type SearchSuburbsRequest struct {
	Search     string
	MaxResults int
}

// SearchSuburbs export
func (c *Client) SearchSuburbs(req SearchSuburbsRequest) ([]SearchSuburb, error) {
	if req.MaxResults == 0 {
		req.MaxResults = 50
	}

	if req.Search == "" {
		return nil, fmt.Errorf("search parameter cannot be empty")
	}

	var results []SearchSuburb

	_, err := c.createRequest().
		SetQueryParams(map[string]string{
			"maxResults": strconv.Itoa(req.MaxResults),
			"searchText": req.Search,
		}).
		SetResult(&results).
		Get(c.host + "/FindSuburbs")

	if err != nil {
		return nil, err
	}

	return results, nil
}

// GetScheduleRequest export
type GetScheduleRequest struct {
	Stages   []Stage
	SuburbID string
}

// GetSchedule export
func (c *Client) GetSchedule(req GetScheduleRequest) (*Schedule, error) {
	// TODO: validate stage

	if req.Stages == nil {
		req.Stages = []Stage{Stage1, Stage2, Stage3, Stage4}
	}

	var days []ScheduleDay

	type block struct {
		dayMonth string
		duration string
		stage    Stage
	}

	var blocks []block

	for _, stage := range req.Stages {

		resp, err := c.createRequest().
			// Needs to make browser like query
			SetHeader("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:69.0) Gecko/20100101 Firefox/69.0").
			Get(c.host + fmt.Sprintf("/GetScheduleM/%s/%d/_/1", req.SuburbID, stage))

		if err != nil {
			return nil, err
		}

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
		if err != nil {
			return nil, err
		}

		doc.Find(".scheduleDay").Each(func(i int, s *goquery.Selection) {
			// Example format Mon, 07 Sep
			dayMonth := strings.TrimSpace(s.Find(".dayMonth").Text())
			duration := strings.TrimSpace(s.Find("a").Text())

			if duration != "" {
				blocks = append(blocks, block{
					dayMonth, duration, stage,
				})
			}
		})
	}

	year := time.Now().Year()

	for _, block := range blocks {
		day, err := parseHTMLTime(year, block.dayMonth, block.duration, c.location, block.stage)

		if err != nil {
			return nil, fmt.Errorf("failed to parse html duration: %w", err)
		}

		days = append(days, day)
	}

	groupedByDate := make(map[time.Time]ScheduleDay)

	// Merge same days
	for _, day := range days {
		date := day.Date

		if val, ok := groupedByDate[date]; ok {
			groupedByDate[date] = ScheduleDay{
				Date:  date,
				Slots: append(val.Slots, day.Slots...),
			}
		} else {
			groupedByDate[date] = day
		}
	}

	// Empty list
	days = nil

	for _, day := range groupedByDate {
		var mergedSlots = make(map[string]ScheduleSlot)

		for _, slot := range day.Slots {
			// Remove duplicate days
			key := fmt.Sprintf("%s %s", slot.Start.Format(time.RFC3339), slot.Duration.String())

			if val, ok := mergedSlots[key]; ok {
				if val.Stage < slot.Stage {
					mergedSlots[key] = slot
				}
			} else {
				mergedSlots[key] = slot
			}
		}

		var slots []ScheduleSlot

		for _, slot := range mergedSlots {
			slots = append(slots, slot)
		}

		sort.Sort(ScheduleSlotByStart(slots))

		days = append(days, ScheduleDay{
			Date:  day.Date,
			Slots: slots,
		})
	}

	sort.Sort(ScheduleSlotByDay(days))

	schedule := Schedule{
		Schedule: days,
	}

	return &schedule, nil
}
