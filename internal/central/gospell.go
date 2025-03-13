package central

import (
	"context"
	"flag"
	"fmt"
	"gospell/internal/api"
	"gospell/internal/definition"
	"gospell/internal/tts"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

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
    fs := flag.NewFlagSet("gospell", flag.ExitOnError)
	credentialFlag := fs.String("credentials", "", "Path to Google Cloud credentials JSON file")
    numDefinitionsFlag := fs.Int("definitions", 1, "Number of definitions to display")

	fs.Parse(os.Args[1:])

	var babbler = babble.NewBabbler() // babbler gets a random word from /usr/share/dict/words
	babbler.Count = 1

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := &texttospeech.Client{}

	// load the cache
	wordsWithoutDefinitions = definition.LoadCache()

	go func() { // signal handler
		sig := <-sigs
		fmt.Printf("\nReceived %s signal:\n", sig)

		// Save the cache before exiting
		fmt.Println("Saving cache before exit...")
		definition.SaveCache(&wordsWithoutDefinitions)

		// Cancel the context to signal other goroutines to clean up
		cancel()

		// Exit cleanly
		os.Exit(0)
	}()

	if *credentialFlag != "" {
		var err error

		client, err = texttospeech.NewClient(ctx, option.WithCredentialsFile(*credentialFlag))
		defer client.Close()
		if err != nil {
			log.Fatal("Bad credentials file")
		}
	} else {
		log.Fatal("No credentials file provided")
	}

	for { // main loop
		word := GetAcceptableWord(babbler)
		responseObject := definition.GetResponse(word)

		go fmt.Print(definition.GetDefinition(responseObject, *numDefinitionsFlag))
		go tts.SayWord(ctx, *client, word)

		handleUserInput(ctx, client, word)
	}
}

func handleUserInput(ctx context.Context, client *texttospeech.Client, word string) {
    for {
        var userInput string
        fmt.Scan(&userInput)
        
        switch userInput {
        case "/r":
            tts.SayWord(ctx, *client, word)
            continue
        case "/q":
            definition.SaveCache(&wordsWithoutDefinitions)
            os.Exit(0)
        }
        
        // Process answer
        if userInput == word {
            fmt.Println("Correct!")
        } else {
            fmt.Println("Incorrect.")
            fmt.Printf("You typed: %s\n", userInput)
            fmt.Printf("The correct spelling is: %s\n", word)
        }
        return
    }
}

func GetAcceptableWord(babbler babble.Babbler) string {
    for {
        word := babbler.Babble()
        if isAcceptableWord(word) {
            return word
        }
    }
}

// Acceptable words are lowercase and contain no special characters & are defined in the dictionary
func isAcceptableWord(word string) bool {
	return !strings.ContainsAny(word, "-_'") && api.IsLower(&word) && definition.IsDefined(word, &wordsWithoutDefinitions)
}
