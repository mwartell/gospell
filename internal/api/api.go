package api

import (
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/jharlan-hash/gospell/internal/definition"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/wav"
	"github.com/tjarratt/babble"
)

func GetAcceptableWord(babbler babble.Babbler, cache map[string]struct{}) string {
	for {
		word := babbler.Babble()
		if isAcceptableWord(word, cache) {
			return word
		}
	}
}

// Acceptable words are lowercase and contain no special characters & are defined in the dictionary
func isAcceptableWord(word string, cache map[string]struct{}) bool {
	return !strings.ContainsAny(word, "-_'") && isLower(&word) && definition.IsDefined(word, cache)
}

func isLower(s *string) bool {
	for _, r := range *s {
		if !unicode.IsLower(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func PlayWav(filepath string) error {
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	streamer, format, err := wav.Decode(f)
	if err != nil {
		return err
	}
	defer streamer.Close()

	sr := format.SampleRate
	speaker.Init(sr, sr.N(time.Second/10))

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	<-done
	return nil
}
