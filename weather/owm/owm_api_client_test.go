package owm

import (
	"fmt"
	"testing"

	"github.com/tada3/triton/weather"
)

func Test_GetCurrentWeather(t *testing.T) {
	client, err := NewOwmClient(weather.OwmBaseURL, weather.OwmAPIKey, 5)
	if err != nil {
		t.Fatal(err)
	}
	spath := "weather"
	cityID := "1848354"
	req, err := client.NewGetRequest(spath, cityID)
	if err != nil {
		t.Fatal(err)
	}

	res, err := client.httpClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	var cw CurrentWeather

	if err := weather.DecodeBody(res, &cw); err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Result: %v\n", cw)
}
