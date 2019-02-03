package weather

import (
	"encoding/json"
	"fmt"
	"net/http"

	"time"

	"github.com/tada3/triton/redis"
	"github.com/tada3/triton/translation"
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

		setCache3(city, cw, countryCode)
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

func getCityIDOrCityNameEN(city *model.CityInfo) (*model.CityInfo, error) {
	// 1. preferred_city
	city, found, err := getCityIDFromPreferredCity(city)
	if err != nil {
		// ignore
		fmt.Printf("ERROR getCityIDFromPreferredCity() failed: %+v\n", err)
	} else if found {
		return city, nil
	}

	// 2. Translate (CityNameEN)
	if city.CityNameEN == "" {
		var ename string
		ename, err = translation.Translate(city.CityName)
		if err != nil {
			// fmt.Println("ERROR Translate() failed: %+v", err)
			//msg := "ごめんなさい、システムの調子が良くないようです。しばらくしてからもう一度お試しください。"
			//return getErrorResponse(msg)
			return city, err
		}
		city.CityNameEN = ename
	}

	// 3. city_list
	return getCityIDFromCityList(city)
}

func getCityIDFromPreferredCity(city *model.CityInfo) (*model.CityInfo, bool, error) {
	if city.CityName == "" {
		return city, false, nil
	}
	if city.CountryCode != "" {
		id, found := tritondb.GetCityIDFromPreferredCity(city.CityName, city.CountryCode)
		if found {
			city.CityID = id
			return city, true, nil
		}
	} else {
		id, code, found := tritondb.GetCityIDFromPreferredCityNoCountry(city.CityName)
		if found {
			city.CityID = id
			city.CountryCode = code
			return city, true, nil
		}
	}
	return city, false, nil
}

func getCityIDFromCityList(city *model.CityInfo) (*model.CityInfo, error) {
	cityID, countryCode, found := tritondb.GetCityID2(city)
	if found {
		city.CityID = cityID
		if city.CountryCode == "" {
			city.CountryCode = countryCode
		}
		return city, nil
	}
	return city, nil
}

func GetCurrentWeather3(city *model.CityInfo) (*model.CurrentWeather, error) {
	// 1. Get CityID or Translate
	city, err := getCityIDOrCityNameEN(city)
	if err != nil {
		return nil, err
	}

	fmt.Printf("INFO city1: %v\n", city)

	// 2. Get Weather from OWM
	var cw *model.CurrentWeather
	if city.CityID != 0 {
		// 2-1. Use cityID
		cw, err = owmc.GetCurrentWeatherByID(city.CityID)
		if err != nil {
			return nil, err
		}
	} else {
		// 2-2 Use CityName and CountryCpde
		fmt.Printf("WARN cityID not found: %v\n", *city)
		cw, err = owmc.GetCurrentWeatherByName2(city.CityNameEN, city.CountryCode)
		if err != nil {
			return nil, err
		}
	}
	cw.CountryCode = city.CountryCode
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
	if city.CountryCode == "" {
		city.CountryCode = cw.CountryCode
	}
	return cw, true
}

func GetCurrentWeatherFromCache(city *model.CityInfo) (*model.CurrentWeather, bool) {
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

func SetCurrentWeatherToCache(city *model.CityInfo, cw *model.CurrentWeather) {
	b, err := json.Marshal(cw)
	if err != nil {
		fmt.Printf("LOG Marshal failed: %s\n", err.Error())
		return
	}
	v := string(b)
	redis.Set(getRedisKey2(city), v, cacheTimeout)
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

func setCache3(city *model.CityInfo, cw *model.CurrentWeather, countryCode string) {
	if city.CountryCode == "" {
		cw.CountryCode = countryCode
	}
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
