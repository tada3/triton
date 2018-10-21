package weather

import (
	"encoding/json"
	"net/http"

	"github.com/tada3/triton/weather/owm"
)

const (
	OwmBaseURL         string = "https://api.openweathermap.org/data/2.5/"
	OwmAPIKey          string = "e3fd219fa4ed7117d68e9fcbda3b298e"
	CurrentWeatherPath string = "weather"
)

var (
	owmc *owm.OwmClient
)

type CurrentWeather struct {
	weather string
	temp    int
}

func init() {
	var err error
	owmc, err = owm.NewOwmClient(OwmBaseURL, OwmAPIKey, 5)
	if err != nil {
		panic(err)
	}
}

func GetCurrentWeatherByID(id string) (*CurrentWeather, error) {
	req, err := owmc.NewGetRequest(CurrentWeatherPath, id)
	if err != nil {
		return nil, err
	}

	res, err := owmc.

	return nil, nil
}

// DecodeBody decode JSON response body to the specified struct.
func DecodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}
