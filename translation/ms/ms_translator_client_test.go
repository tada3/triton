package ms

import (
	"fmt"
	"testing"
)

const (
	msTranslatorBaseURL = "https://api.cognitive.microsofttranslator.com/"
	msTranslatorAPIKey  = "78f442ae6aa44e52a221ec70441328c6"
)

func Test_Translate(t *testing.T) {
	c, err := NewMSTranslatorClient(msTranslatorBaseURL, msTranslatorAPIKey, 5)
	if err != nil {
		t.Fatal(err)
	}

	result, err := c.Translate("ローマ")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("result: %v\n", result)

	if result != "Rome" {
		t.Errorf("Unexpected result! expexted: %s, actual: %s", "Rome", result)
	}
}
