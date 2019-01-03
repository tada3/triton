package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/tada3/triton/weather"
	"github.com/tada3/triton/weather/model"

	"github.com/tada3/triton/translation"

	"github.com/tada3/triton/game"
	"github.com/tada3/triton/protocol"

	"github.com/tada3/triton/tritondb"
)

var (
	masterRepo *game.GameMasterRepo
)

func init() {
	masterRepo = game.NewGameMasterRepo()
}

func Dispatch(w http.ResponseWriter, r *http.Request) {
	req, err := parseRequest(r)
	if err != nil {
		fmt.Printf("JSON decoding failed, %v\n", err.Error())
		respondError(w, "Invalid reqeuest!")
		return
	}

	reqType := req.Request.Type

	userId := getUserId(req)

	var response protocol.CEKResponse

	switch reqType {
	case "LaunchRequest":
		fmt.Println("LaunchRequest")
		response = handleLaunchRequest()
	case "SessionEndedRequest":
		fmt.Println("SessionEndedRequest")
		response = protocol.MakeCEKResponse(handleEndRequest())

	case "IntentRequest":
		fmt.Println("IntentRequest")
		intentName := getIntentName(req)
		if intentName == "CurrentWeather" {
			response = handleCurrentWeather(req, userId)
		} else if intentName == "Tomete" {
			response = handleTomete(req, userId)
		} else if intentName == "Arigato" {
			response = handleArigato(req, userId)
		} else if intentName == "Sugoine" {
			response = handleSugoine(req, userId)
		} else if intentName == "Retry" {
			response = handleStartOver(req, userId)
		} else {
			response = handleUnknownRequest(req)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(&response)
	fmt.Printf("<<< %s\n", string(b))
	w.Write(b)
}

func parseRequest(r *http.Request) (protocol.CEKRequest, error) {
	defer r.Body.Close()

	var req protocol.CEKRequest

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return req, err
	}
	fmt.Printf(">>> %s\n", string(body))

	err = json.Unmarshal(body, &req)
	return req, err
}

func getUserId(req protocol.CEKRequest) string {
	system0 := req.Contexts["System"]
	system, ok := system0.(map[string]interface{})
	if !ok {
		return ""
	}
	user0 := system["user"]
	user, ok := user0.(map[string]string)
	if !ok {
		return ""
	}
	return user["userId"]
}

func handleCurrentWeather(req protocol.CEKRequest, userID string) protocol.CEKResponse {
	var msg string
	var p protocol.CEKResponsePayload

	// 0. Get City
	city := genCityInfoFromSlots(req)
	if city == nil {
		fmt.Printf("LOG No slots were passed: %+v", req.Request.Intent)
		msg = game.GetMessage2(game.NoCity)
		p = protocol.MakeCEKResponsePayload(msg, false)
		return protocol.MakeCEKResponse(p)
	}
	fmt.Printf("city0: %v\n", city)

	// 1. Translation
	var err error
	city.CityNameEN, err = translation.Translate(city.CityName)
	if err != nil {
		fmt.Println("ERROR!", err)
		msg := "ごめんなさい、システムの調子が良くないようです。しばらくしてからもう一度お試しください。"
		return getErrorResponse(msg)
	}
	fmt.Printf("city1: %v\n", city)

	// 2. Get weather
	weather, err := weather.GetCurrentWeather2(city)
	if err != nil {
		fmt.Println("Error!", err.Error())
		msg := "ごめんなさい、システムの調子が良くないみたいです。しばらくしてからもう一度お試しください。"
		return getErrorResponse(msg)
	}

	if weather != nil {
		countryName := ""
		if city.CountryCode != "" && city.CountryCode != "HK" && city.CountryCode != "JP" {
			cn, found := tritondb.CountryCode2CountryName(city.CountryCode)
			if found {
				countryName = cn
			} else {
				fmt.Printf("CountryName is not found: %s\n", city.CountryCode)
			}
		}
		if countryName != "" {
			msg = game.GetMessage(game.CurrentWeather2, countryName, city.CityName, weather.Weather, weather.Temp)
		} else {
			msg = game.GetMessage(game.CurrentWeather, city.CityName, weather.Weather, weather.Temp)
		}
	} else {
		fmt.Printf("Weather for %v is not found.\n", city)
		msg = game.GetMessage2(game.WeatherNotFound, city.CityName)
	}
	p = protocol.MakeCEKResponsePayload(msg, false)
	return protocol.MakeCEKResponse(p)
}

func handleTomete(req protocol.CEKRequest, userID string) protocol.CEKResponse {
	msg := game.GetMessage2(game.Tomete)
	p := protocol.MakeCEKResponsePayload(msg, true)
	return protocol.MakeCEKResponse(p)
}

func handleArigato(req protocol.CEKRequest, userID string) protocol.CEKResponse {
	msg := game.GetMessage2(game.Arigato)
	p := protocol.MakeCEKResponsePayload(msg, false)
	return protocol.MakeCEKResponse(p)
}

func handleSugoine(req protocol.CEKRequest, userID string) protocol.CEKResponse {
	msg := game.GetMessage2(game.Sugoine)
	p := protocol.MakeCEKResponsePayload(msg, false)
	return protocol.MakeCEKResponse(p)
}

func getCityFromCountrySlot2(req protocol.CEKRequest) *model.CityInfo {
	intent := req.Request.Intent
	slots := intent.Slots
	country := protocol.GetStringSlot(slots, "country")
	if country != "" {
		city, found, err := tritondb.CountryCode2City2(country)
		if err != nil {
			fmt.Println("ERROR!", err.Error())
			return nil
		}
		if !found {
			fmt.Printf("WARN: city not found: %s\n", country)
			return nil
		}
		return city
	}

	country = protocol.GetStringSlot(slots, "country_snt")
	if country != "" {
		city, found := tritondb.CountryName2City2(country)
		if !found {
			fmt.Printf("WARN: city not found: %s\n", country)
			return nil
		}
		return city
	}

	country = protocol.GetStringSlot(slots, "ken_jp")
	if country != "" {
		city, found := tritondb.CountryName2City2(country)
		if !found {
			fmt.Printf("WARN: city not found: %s\n", country)
			return nil
		}
		city.CountryCode = "JP"
		return city
	}

	return nil
}

func getCityFromCitySlot2(req protocol.CEKRequest) *model.CityInfo {
	intent := req.Request.Intent
	slots := intent.Slots
	cityInfo := &model.CityInfo{}

	cityInfo.CityName = protocol.GetStringSlot(slots, "city")
	if cityInfo.CityName != "" {
		return cityInfo
	}

	cityInfo.CityName = protocol.GetStringSlot(slots, "city_snt")
	if cityInfo.CityName != "" {
		return cityInfo
	}

	cityInfo.CityName = protocol.GetStringSlot(slots, "city_jp")
	if cityInfo.CityName != "" {
		if strings.HasSuffix(cityInfo.CityName, "市") {
			cityInfo.CityName = strings.TrimRight(cityInfo.CityName, "市")
		}
		cityInfo.CountryCode = "JP"
		return cityInfo
	}

	return nil
}

func getCityFromCountrySlot3(slots map[string]protocol.CEKSlot) *model.CityInfo {
	country := protocol.GetStringSlot(slots, "country")
	if country != "" {
		fmt.Println("AAAA", country)
		city, found, err := tritondb.CountryCode2City2(country)
		fmt.Println("XXXXXX city=%v\n", city)
		if err != nil {
			fmt.Println("ERROR!", err.Error())
			return nil
		}
		if !found {
			fmt.Printf("WARN: city not found: %s\n", country)
			return nil
		}
		return city
	}

	country = protocol.GetStringSlot(slots, "country_snt")
	if country != "" {
		city, found := tritondb.CountryName2City2(country)
		if !found {
			fmt.Printf("WARN: city not found: %s\n", country)
			return nil
		}
		return city
	}

	country = protocol.GetStringSlot(slots, "ken_jp")
	if country != "" {
		city, found := tritondb.CountryName2City2(country)
		if !found {
			fmt.Printf("WARN: city not found: %s\n", country)
			return nil
		}
		city.CountryCode = "JP"
		return city
	}

	return nil
}

func getCityFromCitySlot3(slots map[string]protocol.CEKSlot, cityInfo *model.CityInfo) *model.CityInfo {
	if cityInfo == nil {
		cityInfo = &model.CityInfo{}
	}

	var cityName string

	cityName = protocol.GetStringSlot(slots, "city")
	if cityName != "" {
		cityInfo.CityName = cityName
		return cityInfo
	}

	cityName = protocol.GetStringSlot(slots, "city_snt")
	if cityName != "" {
		cityInfo.CityName = cityName
		return cityInfo
	}

	cityName = protocol.GetStringSlot(slots, "city_jp")
	if cityName != "" {
		if strings.HasSuffix(cityName, "市") {
			cityName = strings.TrimRight(cityName, "市")
		}
		cityInfo.CityName = cityName
		if cityInfo.CountryCode == "" {
			cityInfo.CountryCode = "JP"
		}
		return cityInfo
	}
	return cityInfo
}

func genCityInfoFromSlots(req protocol.CEKRequest) *model.CityInfo {
	intent := req.Request.Intent
	slots := intent.Slots
	cityInfo := getCityFromCountrySlot3(slots)
	fmt.Printf("XXX cityInfo=%v\n", cityInfo)
	return getCityFromCitySlot3(slots, cityInfo)
}

func handleStartNew(req protocol.CEKRequest, userId string) protocol.CEKResponse {
	prevSize := 0
	gm := masterRepo.GetGameMaster(userId)
	if gm != nil {
		gm.Stop()
		prevSize, _ = gm.GetSize()
	}
	gm = game.NewGameMaster()
	masterRepo.AddGameMaster(userId, gm)
	err := gm.StartNew()
	if err != nil {
		fmt.Println("ERROR!", err)
		return handleInvalidRequest(req)
	}

	// After calling of gm.StartNew(). So maze should exist.
	size, _ := gm.GetSize()
	start, err := gm.GetStart()
	if err != nil {
		fmt.Println("ERROR!", err)
		return handleInvalidRequest(req)
	}
	goal, err := gm.GetGoal()
	if err != nil {
		fmt.Println("ERROR!", err)
		return handleInvalidRequest(req)
	}

	var msg1 string
	if size != prevSize {
		msg1 = game.GetMessage(game.START_MSG_NEW, size, size, start, goal)
	} else {
		msg1 = game.GetMessage(game.START_MSG_NEW_SIMPLE, start, goal)
	}
	msg2 := game.GetMessage(game.RepromptMsg2)

	p := protocol.MakeCEKResponsePayload２(msg1, msg2, false)
	return protocol.MakeCEKResponse(p)
}

func handleStartOver(req protocol.CEKRequest, userId string) protocol.CEKResponse {
	gm := masterRepo.GetGameMaster(userId)
	if gm == nil {
		return handleStartNew(req, userId)
	}
	gm.Stop()
	err := gm.StartOver()
	if err != nil {
		fmt.Println("ERROR!", err)
		return handleInvalidRequest(req)
	}

	start, err := gm.GetStart()
	if err != nil {
		fmt.Println("ERROR!", err)
		return handleInvalidRequest(req)
	}
	goal, err := gm.GetGoal()
	if err != nil {
		fmt.Println("ERROR!", err)
		return handleInvalidRequest(req)
	}

	msg1 := game.GetMessage(game.START_MSG_REPEAT, start, goal)
	msg2 := game.GetMessage(game.RepromptMsg2)

	p := protocol.MakeCEKResponsePayload２(msg1, msg2, false)
	return protocol.MakeCEKResponse(p)
}

func location2String222(loc game.Location) string {
	return fmt.Sprintf("%dの%d", loc.X, loc.Y)
}

func handleLaunchRequest() protocol.CEKResponse {
	osVal1 := protocol.MakeOutputSpeechUrlValue(game.GetSoundURL(game.OpeningSound))
	osVal2 := protocol.MakeOutputSpeechTextValue(game.GetMessage(game.WelcomeMsg))
	os := protocol.MakeOutputSpeechList(osVal1, osVal2)
	p := protocol.CEKResponsePayload{
		OutputSpeech:     os,
		ShouldEndSession: false,
	}
	return protocol.MakeCEKResponse(p)
}

func handleEndRequest() protocol.CEKResponsePayload {
	msg := game.GetMessage(game.GoodbyMsg)
	return protocol.CEKResponsePayload{
		OutputSpeech:     protocol.MakeSimpleOutputSpeech(msg),
		ShouldEndSession: true,
	}
}

func handleUnknownRequest(req protocol.CEKRequest) protocol.CEKResponse {
	msg := game.GetMessage(game.InquirelyMsg)
	p := protocol.CEKResponsePayload{
		OutputSpeech:     protocol.MakeSimpleOutputSpeech(msg),
		ShouldEndSession: false,
	}
	return protocol.MakeCEKResponse(p)
}

func handleInvalidRequest(req protocol.CEKRequest) protocol.CEKResponse {
	msg := game.GetMessage2(game.InvalidActionMsg)
	p := protocol.CEKResponsePayload{
		OutputSpeech:     protocol.MakeSimpleOutputSpeech(msg),
		ShouldEndSession: false,
	}
	return protocol.MakeCEKResponse(p)
}

func getErrorResponse(msg string) protocol.CEKResponse {
	p := protocol.CEKResponsePayload{
		OutputSpeech:     protocol.MakeSimpleOutputSpeech(msg),
		ShouldEndSession: true,
	}
	return protocol.MakeCEKResponse(p)
}

func respondError(w http.ResponseWriter, msg string) {
	response := protocol.MakeCEKResponse(
		protocol.CEKResponsePayload{
			OutputSpeech: protocol.MakeSimpleOutputSpeech(msg),
		})

	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(&response)
	w.Write(b)
}

func getIntentName(req protocol.CEKRequest) string {
	name := req.Request.Intent.Name
	return name
}

// HealthCheck just returns 'OK'.
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}
