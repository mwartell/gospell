package definition

import (
	"log"
	"os"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/wav"
)

// PlayWav plays a WAV file using the beep library.
func PlayWav(filepath string) {
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	streamer, format, err := wav.Decode(f)
	if err != nil {
        log.Fatal(err)
	}
	defer streamer.Close()

	sr := format.SampleRate
	speaker.Init(sr, sr.N(time.Second/10))

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	<-done
}
