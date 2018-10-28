package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tada3/triton/weather"

	"github.com/tada3/triton/translation"

	"github.com/tada3/triton/game"
	"github.com/tada3/triton/protocol"
)

const (
	OPENING_SOUND_URL   string = "https://clova-common.line-scdn-dev.net/test/b_099.mp3"
	DEAD_SOUND_URL      string = "https://clova-common.line-scdn-dev.net/test/dead-sound.mp3"
	BUTSUKARU_SOUND_URL string = "https://clova-common.line-scdn-dev.net/test/butsukaru_04.mp3"
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
		} else if intentName == "Retry" {
			response = handleStartOver(req, userId)
		} else if intentName == "Move" {
			response = handleMove(req, userId)
		} else if intentName == "Location" {
			response = handleLocation(req, userId)
		} else if intentName == "GiveUp" {
			response = handleGiveUp(req, userId)
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
	intent := req.Request.Intent
	slots := intent.Slots

	city := protocol.GetStringSlot(slots, "city")
	fmt.Printf("city: %s\n", city)

	// 1. Translation
	cityEn, err := translation.Translate(city)
	if err != nil {
		errMsg := "ごめんなさい、システムの調子が良くないようです。しばらくしてからもう一度お試しください。"
		return getErrorResponse(errMsg)
	}

	fmt.Printf("cityEn: %s\n", cityEn)

	// 2. Get weather
	weather, err := weather.GetCurrentWeather(cityEn)
	if err != nil {
		fmt.Println("Error!", err.Error())
		errMsg := "ごめんなさい、システムの調子が良くないようです。しばらくしてからもう一度お試しください。"
		return getErrorResponse(errMsg)
	}

	msg := game.GetMessage(game.CurrentWeather, city, weather.Weather, weather.Temp)

	p := protocol.MakeCEKResponsePayload(msg, false)
	return protocol.MakeCEKResponse(p)
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

func handleMove(req protocol.CEKRequest, userId string) protocol.CEKResponse {
	gm := masterRepo.GetGameMaster(userId)
	if gm == nil {
		return handleInvalidRequest(req)
	}

	slots := req.Request.Intent.Slots
	direction := protocol.GetStringSlot(slots, "direction")

	var d game.Direction
	if direction == "上" {
		d = game.NORTH
	} else if direction == "下" {
		d = game.SOUTH
	} else if direction == "右" {
		d = game.EAST
	} else if direction == "左" {
		d = game.WEST
	}
	result, err := gm.Move(d)

	if err != nil {
		fmt.Println("ERROR! Failed to move.", err)
		return handleInvalidRequest(req)
	}

	var payload protocol.CEKResponsePayload
	if result {
		if gm.State() == game.GOALED {
			msg1 := game.GetMessage(game.GOAL_MSG, gm.MoveCount(), gm.LocateCount())
			msg2 := game.GetMessage(game.RepromptMsg3)

			os := protocol.MakeSimpleOutputSpeech(msg1)
			rep := protocol.MakeReprompt(protocol.MakeSimpleOutputSpeech(msg2))
			payload = protocol.MakeCEKResponsePayload3(os, rep, false)

		} else {
			msg1 := game.GetMessage2(game.MoveMsg, direction) + game.GetMessage2(game.RepromptMsg2)
			msg2 := game.GetMessage2(game.RepromptMsg2)
			os := protocol.MakeSimpleOutputSpeech(msg1)
			rep := protocol.MakeReprompt(protocol.MakeSimpleOutputSpeech(msg2))
			payload = protocol.MakeCEKResponsePayload3(os, rep, false)
		}
	} else {
		if gm.State() == game.DEAD {
			osVal1 := protocol.MakeOutputSpeechUrlValue(BUTSUKARU_SOUND_URL)
			osVal2 := protocol.MakeOutputSpeechUrlValue(DEAD_SOUND_URL)
			osVal3 := protocol.MakeOutputSpeechTextValue(
				game.GetMessage(game.GameoverMsg) + game.GetMessage(game.RepromptMsg1))
			os := protocol.MakeOutputSpeechList(osVal1, osVal2, osVal3)
			rep := protocol.MakeReprompt(protocol.MakeSimpleOutputSpeech(game.GetMessage(game.RepromptMsg1)))
			payload = protocol.MakeCEKResponsePayload3(os, rep, false)
		} else {
			osVal1 := protocol.MakeOutputSpeechUrlValue(BUTSUKARU_SOUND_URL)
			msg2 := game.GetMessage2Random(game.ItaiMsg, 0.3) + game.GetMessage2(game.ButsukaruMsg)
			osVal2 := protocol.MakeOutputSpeechTextValue(msg2)
			os := protocol.MakeOutputSpeechList(osVal1, osVal2)
			rep := protocol.MakeReprompt(protocol.MakeSimpleOutputSpeech(game.GetMessage(game.RepromptMsg2)))
			payload = protocol.MakeCEKResponsePayload3(os, rep, false)
		}
	}

	return protocol.MakeCEKResponse(payload)
}

func handleLocation(req protocol.CEKRequest, userID string) protocol.CEKResponse {
	gm := masterRepo.GetGameMaster(userID)
	if gm == nil {
		return handleInvalidRequest(req)
	}

	loc, err := gm.Locate()
	if err != nil {
		fmt.Println("ERROR!", err)
		return handleInvalidRequest(req)
	}

	msg1 := game.GetMessage(game.LocationMsg, loc)
	msg2 := game.GetMessage(game.RepromptMsg2)

	p := protocol.MakeCEKResponsePayload２(msg1, msg2, false)
	return protocol.MakeCEKResponse(p)
}

func handleGiveUp(req protocol.CEKRequest, userID string) protocol.CEKResponse {
	gm := masterRepo.GetGameMaster(userID)
	if gm == nil {
		return handleInvalidRequest(req)
	}

	if gm.State() != game.STARTED && gm.State() != game.SEARCHING {
		return handleUnknownRequest(req)
	}

	msg1 := game.GetMessage(game.GiveUpMsg) + game.GetMessage(game.RepromptMsg1)
	msg2 := game.GetMessage(game.RepromptMsg4)

	p := protocol.MakeCEKResponsePayload２(msg1, msg2, false)
	return protocol.MakeCEKResponse(p)
}

func location2String222(loc game.Location) string {
	return fmt.Sprintf("%dの%d", loc.X, loc.Y)
}

func handleLaunchRequest() protocol.CEKResponse {
	osVal1 := protocol.MakeOutputSpeechUrlValue(OPENING_SOUND_URL)
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