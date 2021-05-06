package translation

import (
	"time"

	"github.com/tada3/triton/logging"
	"github.com/tada3/triton/redis"
	"github.com/tada3/triton/translation/ibm"
	"github.com/tada3/triton/tritondb"
)

const (
	ibmTranslatorBaseURL = "https://api.us-south.language-translator.watson.cloud.ibm.com/instances/863c1970-124c-45c1-b86b-a01ecc6fbe57"
	ibmTranslatorAPIKey  = "c4FFCRaBeJs4YlMBrASPRkkYa98yk40ChuguUSrn9Y8J"
	msTranslatorBaseURL  = "https://api.cognitive.microsofttranslator.com/"
	msTranslatorAPIKey   = "987ce01818af428eab13cf5c1b99605a"

	cacheTimeout time.Duration = 24 * time.Hour
)

var (
	tr  Translator
	log *logging.Entry
)

type Translator interface {
	Translate(string) (string, error)
}

func init() {
	log = logging.NewEntry("trans")
	var err error
	//tr, err = ms.NewMSTranslatorClient(msTranslatorBaseURL, msTranslatorAPIKey, 5)
	tr, err = ibm.NewIBMTranslatorClient(ibmTranslatorBaseURL, ibmTranslatorAPIKey, 1)
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
		log.Info("LOG Cache Hit %s\n", w)
		return tw, nil
	}

	log.Info("LOG Cache Miss %s\n", w)
	// 2. Translate by DB
	tw, hit, err = tritondb.TranslateByDB(w)
	if err != nil {
		// Just ignore here
		log.Error("ERROR! TranslateByDB() failed.", err)
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
