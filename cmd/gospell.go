package main

import (
	"context"
	"flag"
	"fmt"
	"gospell/internal/api"
	"gospell/internal/definition"
	"gospell/internal/tts"
	"log"
	"strings"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"github.com/tjarratt/babble"
	"google.golang.org/api/option"
)

var wordsWithoutDefinitions = make(map[string]struct{})

func main() {
	// this will be a program to help study for spelling tests.
	// it will read a list of words from a file or maybe a db, and then quiz the user on them.
	// the user will be able to choose the difficulty level, and the program will
	// generate a quiz based on that level.
	// the program will also keep track of the user's progress, and provide feedback
	// on their performance.
	fmt.Println("gospell - a spelling quiz program")

	credentialFlag := flag.String("credentials", "", "Path to Google Cloud credentials JSON file")

	flag.Parse()

	var babbler = babble.NewBabbler() // babbler gets a random word from /usr/share/dict/words
	babbler.Count = 1

	// start the tts client
	ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
	client := &texttospeech.Client{}

    // load the cache
    wordsWithoutDefinitions = definition.LoadCache("/Users/25jaso/.cache/gospell/cache")

	if *credentialFlag != "" {
		var err error

		client, err = texttospeech.NewClient(ctx, option.WithCredentialsFile(*credentialFlag))
		defer client.Close()
		if err != nil {
			log.Fatal("Bad credentials file")
		}
		fmt.Println("Using credentials file:", *credentialFlag)
	} else {
		log.Fatal("No credentials file provided")
	}

	for { // main loop
		word := getAcceptableWord(babbler)
		responseObject := definition.GetResponse(&word)
        definition.SaveCache("/Users/25jaso/.cache/gospell/cache", &wordsWithoutDefinitions)

		go definition.PrintDefinition(responseObject)
		go tts.SayWord(ctx, *client, word)

		handleUserInput(ctx, client, word)
	}
}

func handleUserInput(ctx context.Context, client *texttospeech.Client, word string) {
	var userInput string

	fmt.Scan(&userInput)

	if userInput == "/r" {
		tts.SayWord(ctx, *client, word)
		handleUserInput(ctx, client, word)
		return
	} else if userInput == "/c" {
		for key, value := range wordsWithoutDefinitions {
			fmt.Println("Key:", key, "Value:", value)
		}
		return
	}

	if userInput == word {
		fmt.Println("Correct!")
	} else {
		fmt.Println("Incorrect.")
		// this will be a function that will provide the correct spelling with error highlighting
		// for now, we will just print the correct spelling
		fmt.Printf("You typed: %s\n", userInput)
		fmt.Printf("The correct spelling is: %s\n", word)
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

// Acceptable words are lowercase and contain no special characters & are defined in the dictionary
func isAcceptableWord(word string) bool {
	return !strings.ContainsAny(word, "-_'") && api.IsLower(&word) && definition.IsDefined(word, &wordsWithoutDefinitions)
}
