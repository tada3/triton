package translation

import (
	"fmt"
	"time"

	"github.com/tada3/triton/redis"
	"github.com/tada3/triton/translation/ms"
)

const (
	msTranslatorBaseURL               = "https://api.cognitive.microsofttranslator.com/"
	msTranslatorAPIKey                = "78f442ae6aa44e52a221ec70441328c6"
	cacheTimeout        time.Duration = 24 * time.Hour
)

var tr Translator

type Translator interface {
	Translate(string) (string, error)
}

func init() {
	var err error
	tr, err = ms.NewMSTranslatorClient(msTranslatorBaseURL, msTranslatorAPIKey, 5)
	if err != nil {
		panic(err)
	}
}

func Translate(w string) (string, error) {
	var tw string
	tw, hit := checkCache(w)
	if hit {
		fmt.Printf("LOG Cache Hit %s\n", w)
		return tw, nil
	}
	fmt.Printf("LOG Cache Miss %s\n", w)
	tw, err := tr.Translate(w)
	if err != nil {
		return "", err
	}
	setCache(w, tw)
	return tw, nil
}

func checkCache(w string) (string, bool) {
	return redis.Get(getRedisKey(w))
}

func setCache(w, tw string) {
	redis.Set(getRedisKey(w), tw, cacheTimeout)
}

func getRedisKey(w string) string {
	return "triton:translation:" + w
}
