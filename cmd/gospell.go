package main

import (
	"context"
	"fmt"
	"gospell/internal/api"
	"gospell/internal/definition"
	"gospell/internal/tts"
	"strings"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
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


	for {
        word := getAcceptableWord(babbler)
		responseObject := definition.GetResponse(&word)

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
        return;
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
	return !strings.ContainsAny(word, "-_'") && api.IsLower(&word) && definition.IsDefined(definition.GetResponse(&word))
}

