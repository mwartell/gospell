package tts

import (
	"bytes"
	"context"
	"fmt"

	"log"
	"time"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/wav"
)

// SayWord takes a word and uses the Google Cloud Text-to-Speech API to generate and play the audio for that word.
// It requires a context and a text-to-speech client.
// It saves the audio to a temporary file and plays it using the definition.PlayWav function.
// If any error occurs during synthesis or playback, it panics with an error message.
func SayWord(ctx context.Context, client texttospeech.Client, word string) {
	// call tts api
	audioContent, err := synthesizeSpeech(ctx, &client, word)
	if err != nil {
		panic(fmt.Sprintf("Failed to synthesize speech: %v", err))
	}

	r := bytes.NewReader(audioContent)
	// Play the audio
	streamer, format, err := wav.Decode(r)
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

// synthesizeSpeech takes a context, a text-to-speech client, and a string of text.
// It configures the synthesis request with the desired voice and audio output format.
func synthesizeSpeech(ctx context.Context, client *texttospeech.Client, text string) ([]byte, error) {
	// set up request
	req := texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{
				Text: text,
			},
		},
		// configure voice
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "en-US",
			Name:         "en-US-Chirp3-HD-Fenrir", // this one makes me laugh bc he's zesty
			// Name: "en-US-Wavenet-D", // this one is more natural
			// Name: "en-US-Neural2-J", // this is the og
			SsmlGender: texttospeechpb.SsmlVoiceGender_MALE,
		},
		// Configure the audio output
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_LINEAR16, // WAV format
		},
	}

	// Call the API
	resp, err := client.SynthesizeSpeech(ctx, &req) // this is confusing but it is two diff functions
	if err != nil {
		return nil, err
	}

	return resp.AudioContent, nil
}
