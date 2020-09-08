package example

import (
	"fmt"
	"github.com/stackworx-go/eskom-loadshedding-api/loadshedding"
	"time"
)

func ExampleClient_getStatus() {
	tz, err := time.LoadLocation("Africa/Johannesburg")

	if err != nil {
		panic(err)
	}

	client := loadshedding.NewClient(tz, false)

	result, err := client.GetStatus()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Status: %v", result)
}

func ExampleClient_GetMunicipalities() {
	tz, err := time.LoadLocation("Africa/Johannesburg")

	if err != nil {
		panic(err)
	}

	client := loadshedding.NewClient(tz, false)

	result, err := client.GetMunicipalities(loadshedding.ProvinceEasternCape)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Municipalities: %v", result)
}

func ExampleClient_GetMunicipalitySuburbs() {
	tz, err := time.LoadLocation("Africa/Johannesburg")

	if err != nil {
		panic(err)
	}

	client := loadshedding.NewClient(tz, false)

	result, err := client.GetMunicipalitySuburbs(loadshedding.GetMunicipalitySuburbsRequest{
		MunicipalityID: "168",
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("Municipalities: %v", result)
}

func ExampleClient_SearchSuburbs() {
	tz, err := time.LoadLocation("Africa/Johannesburg")

	if err != nil {
		panic(err)
	}

	client := loadshedding.NewClient(tz, false)

	result, err := client.SearchSuburbs(loadshedding.SearchSuburbsRequest{
		Search: "Search",
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("Search Suburbs: %v", result)
}

func ExampleClient_GetSchedule() {
	tz, err := time.LoadLocation("Africa/Johannesburg")

	if err != nil {
		panic(err)
	}

	client := loadshedding.NewClient(tz, false)

	result, err := client.GetSchedule(loadshedding.GetScheduleRequest{
		SuburbID: "1",
		Stage:    loadshedding.Stage3,
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("Schedule: %v", result)
}
