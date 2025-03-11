package main

import (
	"context"
	"encoding/json"
	"fmt"
	"gospell/internal"
	"io"
	"log"
	"net/http"
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
	var responseObject Welcome

	for true { // infinite loop (why does go not have while loops?)
        word = getAcceptableWord(babbler)

		responseObject = getResponse(&word)

		printDefinition(responseObject)
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

func isAcceptableWord(word string) bool {
	return !strings.ContainsAny(word, "-_'") && internal.IsLower(&word) && isDefined(getResponse(&word))
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

type Welcome []WelcomeElement

func UnmarshalWelcome(data []byte) (Welcome, error) {
	var r Welcome
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Welcome) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type WelcomeElement struct {
	Word       string     `json:"word"`
	Phonetic   string     `json:"phonetic"`
	Phonetics  []Phonetic `json:"phonetics"`
	Meanings   []Meaning  `json:"meanings"`
	License    License    `json:"license"`
	SourceUrls []string   `json:"sourceUrls"`
}

type License struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Meaning struct {
	PartOfSpeech string       `json:"partOfSpeech"`
	Definitions  []Definition `json:"definitions"`
	Synonyms     []any        `json:"synonyms"`
	Antonyms     []any        `json:"antonyms"`
}

type Definition struct {
	Definition string  `json:"definition"`
	Synonyms   []any   `json:"synonyms"`
	Antonyms   []any   `json:"antonyms"`
	Example    *string `json:"example,omitempty"`
}

type Phonetic struct {
	Text      string   `json:"text"`
	Audio     string   `json:"audio"`
	SourceURL *string  `json:"sourceUrl,omitempty"`
	License   *License `json:"license,omitempty"`
}

func getResponse(word *string) Welcome {
	response, err := http.Get("https://api.dictionaryapi.dev/api/v2/entries/en/" + *word)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	responseObject, err := UnmarshalWelcome(responseData)

	json.Unmarshal(responseData, &responseObject)

	return responseObject
}

// isDefined checks if a word is defined in the dictionary
func isDefined(responseObject Welcome) bool {
	if len(responseObject) == 0 {
		return false
	}
	if len(responseObject[0].Meanings) == 0 {
		return false
	} else {
		return true
	}
}

func printDefinition(responseObject Welcome) {
	if len(responseObject) == 0 {
		fmt.Println("No definition found.")
		return
	}
	fmt.Println(responseObject[0].Meanings[0].Definitions[0].Definition)
}
