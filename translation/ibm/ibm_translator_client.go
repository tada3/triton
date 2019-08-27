package ibm

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

type IBMTranslatorClient struct {
	baseURL    *url.URL
	apiKey     string
	httpClient *http.Client
}

type TranslationRequest struct {
	Text    string `json:text`
	ModelId string `json:model_id`
}

type TranslationResponse struct {
	Results        []TranslationResult `json:"translations"`
	WordCount      int                 `json:"word_count"`
	CharacterCount int                 `json:"character_count"`
}

type TranslationResult struct {
	Text string `json:"translation"`
}

func NewIBMTranslatorClient(baseURL, apiKey string, timeout int) (*IBMTranslatorClient, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	c := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	return &IBMTranslatorClient{
		baseURL:    u,
		apiKey:     apiKey,
		httpClient: c,
	}, nil
}

func (c *IBMTranslatorClient) Translate(w string) (string, error) {
	req, err := c.NewPostRequest("language-translator/api/v3/translate", w)
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

func (c *IBMTranslatorClient) NewPostRequest(spath, w string) (*http.Request, error) {

	// body
	tr := TranslationRequest{
		Text:    w,
		ModelId: "ja-en",
	}

	b := new(bytes.Buffer)
	enc := json.NewEncoder(b)
	enc.Encode(tr)

	// Request
	u := *c.baseURL
	u.Path = path.Join(c.baseURL.Path, spath)

	q := u.Query()
	q.Set("version", "2018-05-01")
	u.RawQuery = q.Encode()

	fmt.Printf("XXX url=%v\n", u.String())

	req, err := http.NewRequest("POST", u.String(), b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	req.SetBasicAuth("apikey", c.apiKey)
	return req, nil
}

// DecodeResponse does not work.
func DecodeResponse(resp *http.Response) (string, error) {
	defer resp.Body.Close()

	var tr TranslationResponse
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&tr)
	if err != nil {
		return "", err
	}

	fmt.Printf("XXX root=%+v\n", tr)
	if len(tr.Results) > 0 {
		r := tr.Results[0]
		return r.Text, nil
	}
	return "", errors.New("invalid response")
}
