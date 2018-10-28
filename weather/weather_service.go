package weather

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tada3/triton/weather/model"
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

func init() {
	var err error
	owmc, err = owm.NewOwmClient(OwmBaseURL, OwmAPIKey, 5)
	if err != nil {
		panic(err)
	}
}

//func GetCurrentWeatherByID(id string) (*model.CurrentWeather, error) {
//idNum, _ := strconv.ParseInt(id, 10, 64)
//return owmc.GetCurrentWeatherByID(idNum)

//return owmc.GetCurrentWeatherByID(idNum)

//}

func GetCurrentWeather(cityName string) (*model.CurrentWeather, error) {
	cityID, found, err := owm.GetCityID(cityName)
	if err != nil {
		// errMsg := "ごめんなさい、システムの調子が良くないようです。しばらくしてからもう一度お試しください。"
		return nil, err
	}
	if !found {
		// use cityEn as it is
		return owmc.GetCurrentWeatherByName(cityName)
	} else {
		// use cityID
		fmt.Printf("cityID: %d\n", cityID)
		return owmc.GetCurrentWeatherByID(cityID)
	}
}

// DecodeBody decode JSON response body to the specified struct.
func DecodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}
