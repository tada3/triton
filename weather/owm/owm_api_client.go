package owm

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"
)

// OwmClient is a Client for OpenWeatherMap API
type OwmClient struct {
	baseURL    *url.URL
	apiKey     string
	httpClient *http.Client
}

type CurrentWeather struct {
	name    string
	weather Weather
}

type Weather struct {
	id          int
	main        string
	description string
	icon        string
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

func (c *OwmClient) newGetRequest(spath, cityID string) (*http.Request, error) {
	u := *c.baseURL
	u.Path = path.Join(c.baseURL.Path, spath)

	q := u.Query()
	q.Set("appid", c.apiKey)
	q.Set("id", cityID)
	u.RawQuery = q.Encode()

	fmt.Printf("XXX url=%v\\n", u.String())

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
