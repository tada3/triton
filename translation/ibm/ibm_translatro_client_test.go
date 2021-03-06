package ibm

import (
	"fmt"
	"testing"
)

const (
	ibmTranslatorBaseURL = "https://gateway.watsonplatform.net/"
	ibmTranslatorAPIKey  = "0X0NrhL0wUYsDmWKAAIegoUhOcUp3cazYNz6-Rv-81uJ"
)

func Test_Translate(t *testing.T) {
	tr, err := NewIBMTranslatorClient(ibmTranslatorBaseURL, ibmTranslatorAPIKey, 5)
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
