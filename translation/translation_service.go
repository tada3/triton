package translation

import (
	"fmt"
	"time"

	"github.com/tada3/triton/redis"
	"github.com/tada3/triton/translation/ms"
	"github.com/tada3/triton/tritondb"
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
	var err error
	var hit bool
	// 1. Check cache
	tw, hit = checkCache(w)
	if hit {
		fmt.Printf("LOG Cache Hit %s\n", w)
		return tw, nil
	}

	fmt.Printf("LOG Cache Miss %s\n", w)
	// 2. Translate by DB
	tw, hit, err = tritondb.TranslateByDB(w)
	if err != nil {
		// Just ignore here
		fmt.Printf("ERROR! TranslateByDB() failed: %v", err.Error())
	}
	if !hit {
		// 3. Translate by MS API
		tw, err = tr.Translate(w)
		if err != nil {
			return "", err
		}
	}

	// 4. Update cache
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
