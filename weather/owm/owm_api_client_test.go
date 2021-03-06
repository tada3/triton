package owm

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"
)

const (
	owmBaseURL string = "https://api.openweathermap.org/data/2.5/"
	owmAPIKey  string = "e3fd219fa4ed7117d68e9fcbda3b298e"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..", "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func Test_GetCurrentWeather(t *testing.T) {
	client, err := NewOwmClient(owmBaseURL, owmAPIKey, 5)
	if err != nil {
		t.Fatal(err)
	}

	// 1609350 Bangkok
	weather, err := client.GetCurrentWeatherByID(1609350)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Result: %v\n", weather)
}

func Test_GetWeatherForecasts(t *testing.T) {
	client, err := NewOwmClient(owmBaseURL, owmAPIKey, 5)
	if err != nil {
		t.Fatal(err)
	}

	// 1848354 Yokohama
	forecasts, err := client.GetWeatherForecastsByID(1848354)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Result: %v\n", forecasts)
}
