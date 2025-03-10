package main

import (
	"fmt"
	"gospell/internal"
	"strings"

	"github.com/amitybell/piper"
	alan "github.com/amitybell/piper-voice-alan"
	"github.com/tjarratt/babble"
)

func main() {
	// this will be a program to help study for spelling tests.
	// it will read a list of words from a file or maybe a db, and then quiz the user on them.
	// the user will be able to choose the difficulty level, and the program will
	// generate a quiz based on that level.
	// the program will also keep track of the user's progress, and provide feedback
	// on their performance.

	fmt.Println("gospell - a spelling quiz program")

	var babbler = babble.NewBabbler()
	babbler.Count = 1

    tts, err := piper.New("", alan.Asset)
    if err != nil {
        panic(err)
    }

    var userInput string
    var word string

	for true {
        word = babbler.Babble()

		// this is true if it contains no special characters and is all lowercase
		wordIsAcceptable := !strings.ContainsAny(word, "-_'") && internal.IsLower(word)

		for !wordIsAcceptable {
			if !wordIsAcceptable {
				word = babbler.Babble()

				wordIsAcceptable = !strings.ContainsAny(word, "-_'") && internal.IsLower(word)
			}
		}

		wav, err := tts.Synthesize(word)
		if err != nil {
			panic(err)
		}
        
		if err := internal.PlayWav(wav); err != nil {
			panic(err)
		}

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

		// get the user's input

		// check the user's input against the correct spelling
		// provide feedback on the user's performance
	}
}
