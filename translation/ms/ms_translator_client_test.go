package ms

import (
	"fmt"
	"testing"

	"github.com/tada3/triton/translation"
)

const (
	msTranslatorBaseURL = "https://api.cognitive.microsofttranslator.com/"
	msTranslatorAPIKey  = "78f442ae6aa44e52a221ec70441328c6"
)

func Test_Translate(t *testing.T) {
	var tr translation.Translator
	tr, err := NewMSTranslatorClient(msTranslatorBaseURL, msTranslatorAPIKey, 5)
	if err != nil {
		t.Fatal(err)
	}

	result, err := tr.Translate("ローマ")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("result: %v\n", result)

	if result != "Rome" {
		t.Errorf("Unexpected result! expexted: %s, actual: %s", "Rome", result)
	}
}
