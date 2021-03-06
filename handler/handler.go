package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/tada3/triton/logging"
	"github.com/tada3/triton/weather"
	"github.com/tada3/triton/weather/model"
	"github.com/tada3/triton/weather/util"

	"github.com/tada3/triton/game"
	"github.com/tada3/triton/protocol"

	"github.com/tada3/triton/tritondb"
)

var (
	log        *logging.Entry
	masterRepo *game.GameMasterRepo
)

func init() {
	log = logging.NewEntry("handler")
	masterRepo = game.NewGameMasterRepo()
}

func Dispatch(w http.ResponseWriter, r *http.Request) {
	req, err := parseRequest(r)
	if err != nil {
		log.Error("JSON decoding failed!", err)
		respondError(w, "Invalid reqeuest!")
		return
	}

	reqType := req.Request.Type
	log.Info("request type: %s", reqType)

	userId := getUserId(req)

	var response protocol.CEKResponse

	switch reqType {
	case "LaunchRequest":
		response = handleLaunchRequest()
	case "SessionEndedRequest":
		response = protocol.MakeCEKResponse(handleEndRequest())

	case "IntentRequest":
		intentName := getIntentName(req)
		if intentName == "CurrentWeather" {
			response = handleCurrentWeather(req, userId)
		} else if intentName == "TomorrowWeather" {
			response = handleTomorrowWeather(req, userId)
		} else if intentName == "Tomete" {
			response = handleTomete(req, userId)
		} else if intentName == "Arigato" {
			response = handleArigato(req, userId)
		} else if intentName == "Sugoine" {
			response = handleSugoine(req, userId)
		} else if intentName == "Doita" {
			response = handleDoita(req, userId)
		} else if intentName == "Question" {
			response = handleQuestion(req, userId)
		} else if intentName == "Samui" {
			response = handleSamui(req, userId)
		} else if intentName == "Clova.YesIntent" {
			response = handleYesIntent(req, userId)
		} else if intentName == "ClovaNoIntent" {
			response = handleNoIntent(req, userId)
		} else {
			response = handleUnknownRequest(req)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(&response)
	log.Info("<<< %s", string(b))
	w.Write(b)
}

func parseRequest(r *http.Request) (protocol.CEKRequest, error) {
	defer r.Body.Close()

	var req protocol.CEKRequest

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return req, err
	}
	log.Info(">>> %s", string(body))

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

func handleTomorrowWeather(req protocol.CEKRequest, userID string) protocol.CEKResponse {
	var msg string
	var p protocol.CEKResponsePayload

	// 0. Get city
	city := genCityInfoFromSlots(req)
	if city == nil || city.CityName == "" {
		log.Info("Cannot get city from slots: %+v", req.Request.Intent)
		msg = game.GetMessage2(game.NoCity)
		p = protocol.MakeCEKResponsePayload(msg, false)
		return protocol.MakeCEKResponse(p)
	}
	var found bool
	msg, found = game.GetMessageForSpecialCity(city.CityName)
	if found {
		p = protocol.MakeCEKResponsePayload(msg, false)
		return protocol.MakeCEKResponse(p)
	}
	log.Info("city: %v", city)

	// 1. Check cache

	// 2. Get tomorrow weather
	tw, err := weather.GetTomorrowWeather(city)
	if err != nil {
		log.Error("Error!", err)
		msg = "ごめんなさい、システムの調子が良くないみたいです。しばらくしてからもう一度お試しください。"
		return getErrorResponse(msg)
	}
	if tw == nil {
		log.Info("Weather for %v is not found.", city)
		msg = game.GetMessage2(game.WeatherNotFound, city.CityName)
	}

	// 3. Generate message
	countryName := ""
	if city.CountryCode != "" && city.CountryCode != "HK" && city.CountryCode != "JP" {
		cn, found := tritondb.CountryCode2CountryName(city.CountryCode)
		if found {
			countryName = cn
		} else {
			log.Info("CountryName is not found: %s\n", city.CountryCode)
		}
	}
	cityName := convertCityName(city.CityName)
	if countryName != "" && countryName != cityName {
		msg = game.GetMessage(game.TomorrowWeather2, util.GetDayStr(tw.Day),
			countryName, cityName, tw.Weather,
			util.GetTempRangeStr(tw.TempMin, tw.TempMax))
	} else {
		msg = game.GetMessage(game.TomorrowWeather, util.GetDayStr(tw.Day),
			cityName, tw.Weather,
			util.GetTempRangeStr(tw.TempMin, tw.TempMax))
	}

	// 5. Make response
	p = protocol.MakeCEKResponsePayload(msg, false)
	return protocol.MakeCEKResponse(p)
}

func handleCurrentWeather(req protocol.CEKRequest, userID string) protocol.CEKResponse {
	var msg string
	var p protocol.CEKResponsePayload

	// 0. Get City
	city := genCityInfoFromSlots(req)
	if city == nil || city.CityName == "" {
		log.Info("Cannot get city from slots: %+v", req.Request.Intent)
		msg = game.GetMessage2(game.NoCity)
		p = protocol.MakeCEKResponsePayload(msg, false)
		return protocol.MakeCEKResponse(p)
	}
	var found bool
	msg, found = game.GetMessageForSpecialCity(city.CityName)
	if found {
		p = protocol.MakeCEKResponsePayload(msg, false)
		return protocol.MakeCEKResponse(p)
	}

	log.Info("city: %v", city)

	// 1. Check cache
	cityInput := city.Clone()
	cw, found := weather.GetCurrentWeatherFromCache(cityInput)
	if !found {
		log.Info("Cache miss: %v", cityInput)
		// 2. Get weather
		var err error
		cw, err = weather.GetCurrentWeather(city)
		if err != nil {
			log.Error("Error!", err)
			msg = "ごめんなさい、システムの調子が良くないみたいです。しばらくしてからもう一度お試しください。"
			return getErrorResponse(msg)
		}

		// 3. Set cache
		// Although city info should have been added more details in GetCurrentWeather,
		// we use the original cityInput as the cache key. If you use the elaborated
		// city info as the cache key, cache would not hit.
		weather.SetCurrentWeatherToCache(cityInput, cw)
	}

	// 4. Generate message
	if cw != nil {
		countryName := ""
		if cw.CountryCode != "" && cw.CountryCode != "HK" && cw.CountryCode != "JP" {
			cn, found := tritondb.CountryCode2CountryName(cw.CountryCode)
			if found {
				countryName = cn
			} else {
				log.Info("CountryName is not found: %s", city.CountryCode)
			}
		}
		cityName := convertCityName(city.CityName)
		if countryName != "" && countryName != cityName {
			msg = game.GetMessage(game.CurrentWeather2, countryName, cityName, cw.Weather, cw.TempStr)
		} else {
			msg = game.GetMessage(game.CurrentWeather, cityName, cw.Weather, cw.TempStr)
		}
	} else {
		log.Info("Weather for %v is not found.", city)
		msg = game.GetMessage2(game.WeatherNotFound, city.CityName)
	}

	// 5. Make response
	p = protocol.MakeCEKResponsePayload(msg, false)
	return protocol.MakeCEKResponse(p)
}

func convertCityName(name string) string {
	if name == "HK" {
		return "香港"
	}
	return name
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

func handleQuestion(req protocol.CEKRequest, userID string) protocol.CEKResponse {
	intent := req.Request.Intent
	slots := intent.Slots
	qitem := protocol.GetStringSlot(slots, "qitem")
	var msg string
	if qitem == "煙霧" {
		msg = game.GetMessage(game.Enmu)
	} else if qitem == "もや" {
		msg = game.GetMessage(game.Moya)
	} else {
		if qitem == "" {
			msg = game.GetMessage2(game.NoCity)
		} else {
			msg = game.GetMessage2(game.UnknownQItem, qitem)
		}
	}
	p := protocol.MakeCEKResponsePayload(msg, false)
	return protocol.MakeCEKResponse(p)
}

func handleSamui(req protocol.CEKRequest, userID string) protocol.CEKResponse {
	msg := game.GetMessage2(game.Samui)
	p := protocol.MakeCEKResponsePayload(msg, false)
	return protocol.MakeCEKResponse(p)
}

func handleYesIntent(req protocol.CEKRequest, userID string) protocol.CEKResponse {
	msg := game.GetMessage2(game.Yes)
	p := protocol.MakeCEKResponsePayload(msg, false)
	return protocol.MakeCEKResponse(p)
}

func handleNoIntent(req protocol.CEKRequest, userID string) protocol.CEKResponse {
	msg := game.GetMessage2(game.No)
	p := protocol.MakeCEKResponsePayload(msg, false)
	return protocol.MakeCEKResponse(p)
}

func handleDoita(req protocol.CEKRequest, userID string) protocol.CEKResponse {
	msg := game.GetMessage2Random(game.Doita, 0.85)
	if msg == "" {
		msg = game.GetOsusumeMessage()
	}
	p := protocol.MakeCEKResponsePayload(msg, false)
	return protocol.MakeCEKResponse(p)
}

// getCityFromCountrySlot3 checks country type slots and create CityInfo with it.
// Second return value represents weather country type slots exists or not.
func getCityFromCountrySlot3(slots map[string]protocol.CEKSlot) *model.CityInfo {
	country := protocol.GetStringSlot(slots, "country_snt")
	if country != "" {
		city, found := tritondb.CountryName2City2(country)
		if found {
			return city
		}
		log.Warn("country not found: %s", country)
	}

	country = protocol.GetStringSlot(slots, "ken_jp")
	if country != "" {
		city, found := tritondb.CountryName2City2(country)
		if found {
			return city
		}
		log.Warn("country not found: %s", country)
	}

	country = protocol.GetStringSlot(slots, "country")
	if country != "" {
		city, found := tritondb.Country2City(country)
		if found {
			return city
		}
		log.Warn("country not found: %s", country)
	}

	return nil
}

// getCityFromPoiSlots checks poi type slots and populates the passed CityInfo.
// Second return value represents weather poi type slots exists or not.
func getCityFromPoiSlots(slots map[string]protocol.CEKSlot, cityInfo *model.CityInfo) (*model.CityInfo, bool) {
	poi := protocol.GetStringSlot(slots, "poi_snt")
	if poi == "" {

		return cityInfo, false
	}

	log.Debug("poi: %s", poi)

	cityInfo, found, err := tritondb.Poi2City(poi, cityInfo)
	if err != nil {
		fmt.Println("ERROR!", err.Error())
	}
	if !found {
		fmt.Printf("WARN: POI not found: %s\n", poi)
	}
	return cityInfo, true
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

	cityInfo, poiExists := getCityFromPoiSlots(slots, cityInfo)
	if poiExists {
		return cityInfo
	}

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
