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

func (c *OwmClient) NewGetRequest(spath string, cityID int64, cityName string) (*http.Request, error) {
	u := *c.baseURL
	u.Path = path.Join(c.baseURL.Path, spath)

	q := u.Query()
	q.Set("appid", c.apiKey)
	q.Set("units", "metric")
	if cityID > 0 {
		q.Set("id", strconv.FormatInt(cityID, 10))
	} else {
		q.Set("q", cityName)
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

	req, err := c.NewGetRequest(CurrentWeatherPath, id, "")
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

	return normalize(ocw), nil
}

func (c *OwmClient) GetCurrentWeatherByName(name string) (*model.CurrentWeather, error) {

	req, err := c.NewGetRequest(CurrentWeatherPath, -1, name)
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

	return normalize(ocw), nil
}

func decodeBody2(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}

func normalize(ocw *OwmCurrentWeather) *model.CurrentWeather {
	weather := getWeatherDescription(ocw.Weather[0].Id)
	temp := marume(ocw.Main.Temp)
	return &model.CurrentWeather{weather, temp}
	// return &CurrentWeatherSummary{ocw.Weather[0].Id, ocw.Main.Temp}
}

func getWeatherDescription(code int64) string {
	return "hoge"
}

func marume(t float64) int64 {
	if t >= 0 {
		return int64(math.Floor(t + 0.5))
	} else {
		// '通常は地上1.25～2.0mの大気の温度を摂氏（℃）単位で表す。度の単位に丸めるときは十分位を四捨五入するが、０度未満は五捨六入する。'
		// by 気象庁
		return int64(math.Ceil(t - 0.5))
	}
}
