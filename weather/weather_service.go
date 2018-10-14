package weather

import (
	"encoding/json"
	"net/http"
)

const (
	OwmBaseURL string = "https://api.openweathermap.org/data/2.5/"
	OwmAPIKey  string = "e3fd219fa4ed7117d68e9fcbda3b298e"
)

type CurrentWeather struct {
	weather string
	temp    int
}

func GetCurrentWeather(city string) *CurrentWeather {

	return nil
}

// DecodeBody decode JSON response body to the specified struct.
func DecodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}
