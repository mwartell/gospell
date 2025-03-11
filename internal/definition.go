package internal

import (
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

	json.Unmarshal(responseData, &responseObject)

	return responseObject
}

// isDefined checks if a word is defined in the dictionary
func IsDefined(responseObject Welcome) bool {
	if len(responseObject) == 0 {
		return false
	}
	if len(responseObject[0].Meanings) == 0 {
		return false
	} else {
		return true
	}
}

func PrintDefinition(responseObject Welcome) {
	if len(responseObject) == 0 {
		fmt.Println("No definition found.")
		return
	}
	fmt.Println(responseObject[0].Meanings[0].Definitions[0].Definition)
}
