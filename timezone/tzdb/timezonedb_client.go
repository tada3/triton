package tzdb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"
)

type TzdbClient struct {
	baseURL    *url.URL
	apiKey     string
	httpClient *http.Client
}

type TzdbTimezone struct {
	Status    string
	GmtOffset int64
	Dst       string // Daylight Saving Time or not
}

func NewTzdbClient(baseURL, apiKey string, timeout int) (*TzdbClient, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	c := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	return &TzdbClient{
		baseURL:    u,
		apiKey:     apiKey,
		httpClient: c,
	}, nil
}

// GetTimezone gets jisa of the specified location at the specified time t
// from timezonedb API. t is a unix timestamp in seconds.
func (c *TzdbClient) GetTimezone(lon float64, lat float64, t int64) (int64, error) {
	req, err := c.NewRequest("v2.1/get-time-zone", lon, lat, t)
	if err != nil {
		return 0, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, err
	}

	tz := new(TzdbTimezone)
	if err := decodeBody2(resp, tz); err != nil {
		return 0, err
	}

	if tz.Status != "OK" {
		return 0, err
	}
	return tz.GmtOffset, nil
}

func (c *TzdbClient) NewRequest(spath string, lon, lat float64, t int64) (*http.Request, error) {
	// url
	u := *c.baseURL
	u.Path = path.Join(c.baseURL.Path, spath)

	q := u.Query()
	q.Set("key", c.apiKey)
	q.Set("format", "json")
	q.Set("by", "position")
	q.Set("lat", strconv.FormatFloat(lat, 'f', -1, 64))
	q.Set("lng", strconv.FormatFloat(lon, 'f', -1, 64))
	q.Set("time", strconv.FormatInt(t, 10))
	u.RawQuery = q.Encode()

	fmt.Printf("XXX url=%v\n", u.String())

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func decodeBody2(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}
