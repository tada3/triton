package ms

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"
)

type MSTranslatorClient struct {
	baseURL    *url.URL
	apiKey     string
	httpClient *http.Client
}

type TranslationRequest struct {
	Text string
}

type TranslationResponse struct {
	Translations []TextAndTo `json:"translations"`
}

type TextAndTo struct {
	Text string `json:"text"`
	To   string `json:"to"`
}

func NewMSTranslatorClient(baseURL, apiKey string, timeout int) (*MSTranslatorClient, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	c := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	return &MSTranslatorClient{
		baseURL:    u,
		apiKey:     apiKey,
		httpClient: c,
	}, nil
}

func (c *MSTranslatorClient) Translate(w string) (string, error) {
	req, err := c.NewPostRequest("translate", w)
	if err != nil {
		return "", err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	result, err := DecodeResponse(res)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (c *MSTranslatorClient) NewPostRequest(spath, w string) (*http.Request, error) {

	// body
	tr := TranslationRequest{
		Text: w,
	}
	root := []TranslationRequest{tr}

	b := new(bytes.Buffer)
	enc := json.NewEncoder(b)
	enc.Encode(root)

	// Request
	u := *c.baseURL
	u.Path = path.Join(c.baseURL.Path, spath)

	q := u.Query()
	q.Set("api-version", "3.0")
	q.Set("from", "ja")
	q.Set("to", "en")
	u.RawQuery = q.Encode()

	fmt.Printf("url=%v\n", u.String())

	req, err := http.NewRequest("POST", u.String(), b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Ocp-Apim-Subscription-Key", c.apiKey)

	return req, nil
}

// DecodeResponse does not work.
func DecodeResponse(resp *http.Response) (string, error) {
	defer resp.Body.Close()

	root := make([]TranslationResponse, 0, 1)

	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&root)
	if err != nil {
		return "", err
	}

	fmt.Printf("root=%+v\n", root)
	if len(root) > 0 {
		t := root[0].Translations
		if len(t) > 0 {
			return t[0].Text, nil
		}
	}
	return "", errors.New("invalid response")

}

func DecodeResponse2(resp *http.Response) (string, error) {
	defer resp.Body.Close()

	var root []map[string]interface{}

	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&root)
	if err != nil {
		return "", err
	}

	fmt.Printf("XXX root=%+v\n", root)
	if len(root) > 0 {
		r0 := root[0]
		if t, ok := r0["translations"]; ok {
			ts := t.([]interface{})
			if len(ts) > 0 {
				ts0 := ts[0]
				ts0m := ts0.(map[string]interface{})
				if txt, ok := ts0m["text"]; ok {
					result := txt.(string)
					return result, nil
				}
			}
		}
	}
	return "__NA__", nil

}
