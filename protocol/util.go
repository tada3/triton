package protocol

// MakeCEKResponse creates CEKResponse instance with given params
func MakeCEKResponse(responsePayload CEKResponsePayload) CEKResponse {
	response := CEKResponse{
		Version:  "0.1.0",
		Response: responsePayload,
	}

	return response
}

func MakeCEKResponsePayload２(msg1 string, msg2 string, ses bool) CEKResponsePayload {
	p := CEKResponsePayload{
		OutputSpeech:     MakeSimpleOutputSpeech(msg1),
		Reprompt:         MakeReprompt(MakeSimpleOutputSpeech(msg2)),
		ShouldEndSession: ses,
	}
	return p
}

func MakeCEKResponsePayload3(os OutputSpeech, rep Reprompt, ses bool) CEKResponsePayload {
	return CEKResponsePayload{
		OutputSpeech:     os,
		Reprompt:         rep,
		ShouldEndSession: ses,
	}
}

func MakeReprompt(os OutputSpeech) Reprompt {
	return Reprompt{
		OutputSpeech: os,
	}
}

// MakeOutputSpeech creates OutputSpeech instance with given params
func MakeSimpleOutputSpeech(msg string) OutputSpeech {
	return OutputSpeech{
		Type: "SimpleSpeech",
		Values: Value{
			Lang:  "ja",
			Value: msg,
			Type:  "PlainText",
		},
	}
}

// MakeOutputSpeechList creates OutputSpeech of type 'SpeechList'.
func MakeOutputSpeechList(value ...Value) OutputSpeech {
	return OutputSpeech{
		Type:   "SpeechList",
		Values: value,
	}
}

func MakeOutputSpeechTextValue(msg string) Value {
	return Value{
		Lang:  "ja",
		Value: msg,
		Type:  "PlainText",
	}
}

func MakeOutputSpeechUrlValue(url string) Value {
	return Value{
		Lang:  "",
		Value: url,
		Type:  "URL",
	}
}

func GetStringSlot(slots map[string]CEKSlot, name string) string {
	slot, ok := slots[name]
	if ok {
		return slot.Value
	}
	return ""
}
