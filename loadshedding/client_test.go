package loadshedding

import (
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func createClient(t *testing.T) (*Client, *time.Location) {
	tz, err := time.LoadLocation("Africa/Johannesburg")

	if err != nil {
		require.NoError(t, err)
	}

	return NewClient(tz, false), tz
}

func Test_GetStatus(t *testing.T) {
	assert := assert.New(t)
	client, _ := createClient(t)
	httpmock.ActivateNonDefault(client.client.GetClient())
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://loadshedding.eskom.co.za/LoadShedding/GetStatus",
		httpmock.NewStringResponder(200, `2`))

	stage, err := client.GetStatus()
	assert.NoError(err)
	assert.Equal(Stage(3), stage)
}

func Test_GetMunicipalities(t *testing.T) {
	assert := assert.New(t)
	client, _ := createClient(t)
	httpmock.ActivateNonDefault(client.client.GetClient())
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://loadshedding.eskom.co.za/LoadShedding/GetMunicipalities",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200,
				`[{"Selected":false,"Text":"Amahlathi","Value":"100"},{"Selected":false,"Text":"Baviaans","Value":"101"},{"Selected":false,"Text":"Blue Crane Route","Value":"102"}]`)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

	results, err := client.GetMunicipalities(EasternCape)
	assert.NoError(err)
	assert.Equal([]Municipality{
		{
			ID:   "100",
			Name: "Amahlathi",
		},
		{
			ID:   "101",
			Name: "Baviaans",
		},
		{
			ID:   "102",
			Name: "Blue Crane Route",
		},
	}, results)
}

func Test_GetMunicipalitySuburbs(t *testing.T) {
	// TODO: pagination test

	assert := assert.New(t)
	client, _ := createClient(t)
	httpmock.ActivateNonDefault(client.client.GetClient())
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://loadshedding.eskom.co.za/LoadShedding/GetSurburbData",
		func(req *http.Request) (*http.Response, error) {
			pageSize := req.URL.Query().Get("pageSize")
			var resp *http.Response

			if pageSize == "0" {
				resp = httpmock.NewStringResponse(200,
					`{
						"Total": 360,
						"Results": []
					  }
					  `)
			} else {
				resp = httpmock.NewStringResponse(200,
					`{
						"Total": 4,
						"Results": [
						  { "id": "1372", "text": "Aandrus", "Tot": 0 },
						  { "id": "1373", "text": "Abelshoek", "Tot": 0 },
						  { "id": "1374", "text": "Adamskraal", "Tot": 272 },
						  { "id": "1375", "text": "Advice", "Tot": 272 }
						]
					  }
					  `)
			}
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

	results, err := client.GetMunicipalitySuburbs(GetMunicipalitySuburbsRequest{
		MunicipalityID: "168",
	})
	assert.NoError(err)
	assert.Equal([]Suburb{
		{
			ID:    "1372",
			Name:  "Aandrus",
			Total: 0,
		},
		{
			ID:    "1373",
			Name:  "Abelshoek",
			Total: 0,
		},
		{
			ID:    "1374",
			Name:  "Adamskraal",
			Total: 272,
		},
		{
			ID:    "1375",
			Name:  "Advice",
			Total: 272,
		},
	}, results)
	assert.Equal(2, httpmock.GetTotalCallCount())
}

func Test_SearchSuburbs(t *testing.T) {
	assert := assert.New(t)
	client, _ := createClient(t)
	httpmock.ActivateNonDefault(client.client.GetClient())
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://loadshedding.eskom.co.za/LoadShedding/FindSuburbs",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200,
				`[
					{
					  "MunicipalityName": "Sundays River Valley",
					  "ProvinceName": "Eastern Cape",
					  "Name": "Allandale",
					  "Id": 10186,
					  "Total": 270
					},
					{
					  "MunicipalityName": "Nelson Mandela Bay",
					  "ProvinceName": "Eastern Cape",
					  "Name": "Allan Heights",
					  "Id": 8187,
					  "Total": 0
					},
					{
					  "MunicipalityName": "Kopanong",
					  "ProvinceName": "Free State",
					  "Name": "Allandale",
					  "Id": 11646,
					  "Total": 270
					}
				]`)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

	results, err := client.SearchSuburbs(SearchSuburbsRequest{
		Search: "Allan",
	})
	assert.NoError(err)
	assert.Equal([]SearchSuburb{
		{
			ID:           10186,
			Province:     "Eastern Cape",
			Municipality: "Sundays River Valley",
			Suburb:       "Allandale",
			Total:        270,
		},
		{
			ID:           8187,
			Province:     "Eastern Cape",
			Municipality: "Nelson Mandela Bay",
			Suburb:       "Allan Heights",
			Total:        0,
		},
		{
			ID:           11646,
			Province:     "Free State",
			Municipality: "Kopanong",
			Suburb:       "Allandale",
			Total:        270,
		},
	}, results)
}

func Test_GetSchedule(t *testing.T) {
	assert := assert.New(t)
	client, tz := createClient(t)
	httpmock.ActivateNonDefault(client.client.GetClient())
	defer httpmock.DeactivateAndReset()

	scheduleHTML, err := ioutil.ReadFile("testdata/schedule.html")
	assert.NoError(err)

	httpmock.RegisterResponder("GET", "https://loadshedding.eskom.co.za/LoadShedding/GetScheduleM/64106/3/_/1",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, string(scheduleHTML))
			resp.Header.Set("Content-Type", "text/html")
			return resp, nil
		})

	results, err := client.GetSchedule(GetScheduleRequest{
		Stage:    Stage1,
		SuburbID: "64106",
	})
	assert.NoError(err)
	assert.Equal(23, len(results.Schedule))

	day := results.Schedule[0]
	assert.Equal(time.Date(
		2020,
		time.September,
		7,
		0,
		0,
		0,
		0,
		tz,
	), day.Day)

	assert.Equal(
		[]ScheduleTime{{
			StartTime: time.Date(
				2020,
				time.September,
				7,
				0,
				0,
				0,
				0,
				tz,
			),
			EndTime: time.Date(
				2020,
				time.September,
				7,
				2,
				30,
				0,
				0,
				tz,
			),
		}}, day.Times)
}
