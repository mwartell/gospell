package internal

import (
	"os"
	"time"
	"unicode"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/wav"
)

func IsLower(s *string) bool {
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
    defer f.Close();

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
