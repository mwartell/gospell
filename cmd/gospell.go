package main

import (
	"context"
	"fmt"
	"gospell/internal"
	"os"
	"strings"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"github.com/tjarratt/babble"
	"google.golang.org/api/option"
)

func main() {
	// this will be a program to help study for spelling tests.
	// it will read a list of words from a file or maybe a db, and then quiz the user on them.
	// the user will be able to choose the difficulty level, and the program will
	// generate a quiz based on that level.
	// the program will also keep track of the user's progress, and provide feedback
	// on their performance.
	fmt.Println("gospell - a spelling quiz program")

	// start the tts client
	ctx := context.Background()
	client, err := texttospeech.NewClient(ctx, option.WithCredentialsFile("/Users/25jaso/Library/.jack/secret.json"))
	if err != nil {
		panic(fmt.Sprintf("Failed to create client: %v", err))
	}
	defer client.Close()

	var babbler = babble.NewBabbler() // babbler gets a random word from /usr/share/dict/words
	babbler.Count = 1

	var userInput string
	var word string

	for true { // infinite loop (why does go not have while loops?)
		word = getAcceptableWord(babbler)
		responseObject := internal.GetResponse(&word)

		internal.PrintDefinition(responseObject)
		sayWord(ctx, *client, word)

		fmt.Scan(&userInput) // this is temporary, we will use a TUI later

		if userInput == word {
			fmt.Println("Correct!")
		} else {
			fmt.Println("Incorrect.")
			// this will be a function that will provide the correct spelling with error highlighting
			// for now, we will just print the correct spelling
			fmt.Printf("You typed: %s\n", userInput)
			fmt.Printf("The correct spelling is: %s\n", word)
		}

		// provide feedback on the user's performance
	}
}

func getAcceptableWord(babbler babble.Babbler) string {
	word := babbler.Babble()

	if isAcceptableWord(word) {
		return word
	} else {
		return getAcceptableWord(babbler)
	}
}

// isAcceptableWord checks if a word is acceptable for the quiz
// acceptable words are lowercase and contain no special characters & are defined in the dictionary
func isAcceptableWord(word string) bool {
	return !strings.ContainsAny(word, "-_'") && internal.IsLower(&word) && internal.IsDefined(internal.GetResponse(&word))
}

func sayWord(ctx context.Context, client texttospeech.Client, word string) {
	// call tts api
	audioContent, err := synthesizeSpeech(ctx, &client, word)
	if err != nil {
		panic(fmt.Sprintf("Failed to synthesize speech: %v", err))
	}

	// save audio to temporary file
	tempFile := "temp.wav"
	if err := os.WriteFile(tempFile, audioContent, 0644); err != nil {
		panic(fmt.Sprintf("Failed to write audio content to file: %v", err))
	}

	// Play the audio
	if err := internal.PlayWav(tempFile); err != nil {
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
			Name:         "en-US-Neural2-J", // lowkey just picked this one for fun
			SsmlGender:   texttospeechpb.SsmlVoiceGender_MALE,
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
