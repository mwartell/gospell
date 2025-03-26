package tts

import (
	"context"
	"fmt"
	"os"

	"github.com/jharlan-hash/gospell/internal/definition"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
)

func SayWord(ctx context.Context, client texttospeech.Client, word string) {
	// call tts api
	audioContent, err := synthesizeSpeech(ctx, &client, word)
	if err != nil {
		panic(fmt.Sprintf("Failed to synthesize speech: %v", err))
	}

	// save audio to temporary file
	tempFile := "./audio/temp.wav"
	if err := os.WriteFile(tempFile, audioContent, 0644); err != nil {
		panic(fmt.Sprintf("Failed to write audio content to file: %v", err))
	}

	// Play the audio
	if err := definition.PlayWav(tempFile); err != nil {
		panic(err)
	}
}

// synthesizeSpeech calls the Google Cloud Text-to-Speech API to generate speech
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
