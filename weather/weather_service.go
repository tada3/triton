package weather

import (
	"encoding/json"
	"fmt"
	"net/http"

	"time"

	"github.com/tada3/triton/redis"
	"github.com/tada3/triton/tritondb"
	"github.com/tada3/triton/weather/model"
	"github.com/tada3/triton/weather/owm"
)

const (
	OwmBaseURL         string        = "https://api.openweathermap.org/data/2.5/"
	OwmAPIKey          string        = "e3fd219fa4ed7117d68e9fcbda3b298e"
	CurrentWeatherPath string        = "weather"
	cacheTimeout       time.Duration = 30 * time.Minute
	keyFmt             string        = "triton:weather:%s:%s"
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

func GetCurrentWeather2(city *model.CityInfo) (*model.CurrentWeather, error) {
	// 1. Check cache
	cw, found := checkCache2(city)
	if found {
		fmt.Printf("LOG Cache2 Hit: %v\n", city)
		if cw == nil {
			fmt.Println("LOG Cach2 cached data is nil.")
		}
		return cw, nil
	}
	fmt.Printf("LOG Cache2 Miss: %v\n", city)

	// 2. Get CityID
	cityID, countryCode, found := tritondb.GetCityID2(city)

	// 3. Get Weather from OWM
	var err error
	if found {
		// use cityID
		fmt.Printf("cityID: %d, countryCode: %s\n", cityID, countryCode)
		cw, err = owmc.GetCurrentWeatherByID(cityID)
		if err != nil {
			return nil, err
		}

		setCache2(city, cw)
		// key of the cache data should be the original city struct
		city.CountryCode = countryCode
		return cw, nil
	} else {
		// use cityID
		fmt.Printf("cityID not found: %v\n", *city)
		cw, err = owmc.GetCurrentWeatherByName2(city.CityNameEN, city.CountryCode)
		if err != nil {
			return nil, err
		}

		setCache2(city, cw)
		return cw, nil
	}
}

func GetCurrentWeather(cityName string) (*model.CurrentWeather, error) {
	var cw *model.CurrentWeather
	cw, ok := checkCache(cityName)
	if ok {
		if cw == nil {
			fmt.Printf("LOG Cach2 Hit, but no weather data: %s\n", cityName)
			return nil, nil
		}
		fmt.Printf("LOG Cache2 Hit: %s\n", cityName)
		return cw, nil
	}
	fmt.Printf("LOG Cache2 Miss: %s\n", cityName)

	cityID, found, err := tritondb.GetCityID(cityName)
	if err != nil {
		// ignore error hrere
		fmt.Printf("LOG DB error: %s\n", err.Error())
		found = false
	}
	var err2 error
	if !found {
		// use cityEn as it is
		cw, err2 = owmc.GetCurrentWeatherByName(cityName)
	} else {
		// use cityID
		fmt.Printf("cityID: %d\n", cityID)
		cw, err2 = owmc.GetCurrentWeatherByID(cityID)
	}
	if err2 != nil {
		return nil, err2
	}
	// Note that setCache() is called even when cw is nil.
	setCache(cityName, cw)
	return cw, nil
}

// DecodeBody decode JSON response body to the specified struct.
func DecodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}

func checkCache2(city *model.CityInfo) (*model.CurrentWeather, bool) {
	v, ok := redis.Get(getRedisKey2(city))
	if !ok {
		return nil, false
	}
	if v == "null" {
		// Weather was not found last time.
		return nil, true
	}
	cw := &model.CurrentWeather{}
	err := json.Unmarshal([]byte(v), cw)
	if err != nil {
		fmt.Printf("LOG Unmarshal failed: %s\n", err.Error())
		return nil, false
	}
	return cw, true
}

func checkCache(cityName string) (*model.CurrentWeather, bool) {
	v, ok := redis.Get(getRedisKey(cityName))
	if !ok {
		return nil, false
	}
	if v == "null" {
		// Weather was not found last time.
		return nil, true
	}
	cw := &model.CurrentWeather{}
	err := json.Unmarshal([]byte(v), cw)
	if err != nil {
		fmt.Printf("LOG Unmarshal failed: %s\n", err.Error())
		return nil, false
	}
	return cw, true
}

func setCache(cityName string, cw *model.CurrentWeather) {
	b, err := json.Marshal(cw)
	if err != nil {
		fmt.Printf("LOG Marshal failed: %s\n", err.Error())
		return
	}
	v := string(b)
	redis.Set(getRedisKey(cityName), v, cacheTimeout)
}

func setCache2(city *model.CityInfo, cw *model.CurrentWeather) {
	b, err := json.Marshal(cw)
	if err != nil {
		fmt.Printf("LOG Marshal failed: %s\n", err.Error())
		return
	}
	v := string(b)
	redis.Set(getRedisKey2(city), v, cacheTimeout)
}

func getRedisKey2(city *model.CityInfo) string {
	key1 := city.CountryCode
	if key1 == "" {
		key1 = "NONE"
	}
	key2 := city.CityNameEN
	return fmt.Sprintf(keyFmt, key1, key2)
}

func getRedisKey(w string) string {
	return "triton:weather:" + w
}
