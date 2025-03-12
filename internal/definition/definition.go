package definition

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Welcome []WelcomeElement

func UnmarshalWelcome(data []byte) (Welcome, error) {
	var r Welcome
	err := json.Unmarshal(data, &r)
	return r, err
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

func GetResponse(word *string) Welcome {
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

	return responseObject
}

// isDefined checks if a word is defined in the dictionary
func IsDefined(word string, wordsWithoutDefinitions *map[string]struct{}) bool {
	// first, search the cache
	// if found, return false
	// if not found, proceed to API call

	if wordInCache(word, wordsWithoutDefinitions) {
        fmt.Println("WORD WAS IN CACHE!! finally my work paid off for about 100 ms of timesave!! the word was: ", word)
		return false
	}

	responseObject := GetResponse(&word)

	if len(responseObject) == 0 {
        addToCache(word, wordsWithoutDefinitions)
		return false
	}
	if len(responseObject[0].Meanings) == 0 {
        fmt.Println("the other one called")
		return false
	} else {
		return true
	}
}

func addToCache(word string, wordsWithoutDefinitions *map[string]struct{}) {
	(*wordsWithoutDefinitions)[word] = struct{}{}
}

func wordInCache(word string, wordsWithoutDefinitions *map[string]struct{}) bool {
	if _, exists := (*wordsWithoutDefinitions)[word]; exists {
		return true
	} else {
		return false
	}
}

func SaveCache(file string, cache *map[string]struct{}) {
	f, err := os.Create(file)
	if err != nil {
		panic("cant open file")
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	if err := enc.Encode(*cache); err != nil {
		panic("cant encode")
	}
}

func LoadCache(file string) (cache map[string]struct{}) {
	f, err := os.Open(file)
	if err != nil {
		panic("cant open file")
	}
	defer f.Close()

	enc := gob.NewDecoder(f)
	if err := enc.Decode(&cache); err != nil {
		panic("cant decode")
	}

	return cache
}

// TODO: let the number of definitions printed be parameterized
func PrintDefinition(responseObject Welcome) {
	if len(responseObject) == 0 {
		fmt.Println("No definition found.")
		return
	}

	// fmt.Println("Definitions:")
	// for _, meaning := range responseObject[0].Meanings {
	// 	for _, definitions := range meaning.Definitions {
	// 		fmt.Print(meaning.PartOfSpeech + ": ")
	// 		fmt.Println(definitions.Definition)
	// 	}
	// }

	// uncomment to print just the first definition
	fmt.Println(responseObject[0].Meanings[0].Definitions[0].Definition)
}
