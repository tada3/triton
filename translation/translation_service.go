package translation

import "github.com/tada3/triton/translation/ms"

const (
	msTranslatorBaseURL = "https://api.cognitive.microsofttranslator.com/"
	msTranslatorAPIKey  = "78f442ae6aa44e52a221ec70441328c6"
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
	return tr.Translate(w)
}
