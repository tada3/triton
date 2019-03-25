package owm

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/tada3/triton/weather/model"
)

const (
	CurrentWeatherPath string = "weather"
	tempStrFormatP     string = "%d度"
	tempStrFormatN     string = "氷点下%d度"
)

// OwmClient is a Client for OpenWeatherMap API
type OwmClient struct {
	baseURL    *url.URL
	apiKey     string
	httpClient *http.Client
}

type OwmCurrentWeather struct {
	Name    string
	Weather []OwmWeather
	Main    OwmMain
	// Usually integer but string in case of 404
	Cod     json.Number
	Message string
}

type OwmWeather struct {
	Id          int64
	Main        string
	Description string
	Icon        string
}

type OwmMain struct {
	Temp     float64
	Pressure float64
	Humidity int64
	Temp_min float64
	Temp_max float64
}

func NewOwmClient(baseURL, apiKey string, timeout int) (*OwmClient, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	c := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	return &OwmClient{
		baseURL:    u,
		apiKey:     apiKey,
		httpClient: c,
	}, nil
}

func (c *OwmClient) NewGetRequest(spath string, cityID int64, cityName, countryCode string) (*http.Request, error) {
	u := *c.baseURL
	u.Path = path.Join(c.baseURL.Path, spath)

	q := u.Query()
	q.Set("appid", c.apiKey)
	q.Set("units", "metric")
	if cityID > 0 {
		q.Set("id", strconv.FormatInt(cityID, 10))
	} else {
		qParam := cityName
		if countryCode != "" {
			qParam = qParam + "," + countryCode
		}
		q.Set("q", qParam)
	}
	u.RawQuery = q.Encode()

	fmt.Printf("XXX url=%v\n", u.String())

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

// GetCurrentWeatherByID gets current weather using OWM API
// and returns it as CurrentWeather.
func (c *OwmClient) GetCurrentWeatherByID(id int64) (*model.CurrentWeather, error) {

	req, err := c.NewGetRequest(CurrentWeatherPath, id, "", "")
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	ocw := new(OwmCurrentWeather)

	if err := decodeBody2(res, ocw); err != nil {
		return nil, err
	}

	fmt.Printf("XXX ocw=%+v\n", ocw)

	return normalize(ocw)
}

func (c *OwmClient) GetCurrentWeatherByName(name, code string) (*model.CurrentWeather, error) {

	req, err := c.NewGetRequest(CurrentWeatherPath, -1, name, code)
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	ocw := new(OwmCurrentWeather)

	if err := decodeBody2(res, ocw); err != nil {
		return nil, err
	}

	fmt.Printf("YYY ocw=%+v\n", ocw)

	return normalize(ocw)
}

func decodeBody2(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}

func normalize(ocw *OwmCurrentWeather) (*model.CurrentWeather, error) {
	var weather string
	var temp int64
	cod := ocw.Cod.String()
	if cod != "200" || len(ocw.Weather) == 0 {
		fmt.Printf("LOG Errorneous response: %+v\n", ocw)
		if cod == "404" {
			return nil, nil
		}
		return nil, fmt.Errorf("Received error response from OWM: %s, %s", cod, ocw.Message)
	}
	weather = GetWeatherCondition(ocw.Weather[0].Id)
	temp = marume(ocw.Main.Temp)
	tempStr := getTempStr(temp)
	fmt.Printf(" tempStr = %s\n", tempStr)
	return &model.CurrentWeather{
		Weather: weather,
		Temp:    temp,
		TempStr: tempStr}, nil
}

func marume(t float64) int64 {
	if t < 0 {
		// '通常は地上1.25～2.0mの大気の温度を摂氏（℃）単位で表す。度の単位に丸めるときは十分位を四捨五入するが、０度未満は五捨六入する。'
		// by 気象庁
		return int64(math.Ceil(t - 0.5))
	}
	return int64(math.Floor(t + 0.5))
}

func getTempStr(t int64) string {
	if t < 0 {
		return fmt.Sprintf(tempStrFormatN, -1*t)
	}
	return fmt.Sprintf(tempStrFormatP, t)
}
