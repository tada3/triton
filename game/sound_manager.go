package game

import (
	"strings"

	"github.com/tada3/triton/config"
)

const (
	_ = iota
	OpeningSound
)

var (
	baseURL string
	names   = [...]string{"", "b_099.mp3"}
)

// SoundType is a type representing sound type.
type SoundType int

func init() {
	cfg := config.GetConfig()
	baseURL = cfg.SoundFileBaseUrl
	if !strings.HasSuffix(baseURL, "/") {
		baseURL = baseURL + "/"
	}
}

// GetSoundURL returns url of the sound file for the specified type.
func GetSoundURL(t SoundType) string {
	if t < OpeningSound || t > OpeningSound {
		return ""
	}
	return baseURL + names[t]
}
