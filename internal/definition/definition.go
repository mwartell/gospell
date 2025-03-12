package definition

import (
	"encoding/gob"
	"encoding/json"
	"errors"
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

func GetResponse(word string) Welcome {
	response, err := http.Get("https://api.dictionaryapi.dev/api/v2/entries/en/" + word)
	if err != nil {
        log.Fatal(err)
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
		fmt.Println("word was in cache")
		return false
	}

	responseObject := GetResponse(word)

	if len(responseObject) == 0 {
		addToCache(word, wordsWithoutDefinitions)
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

func SaveCache(cache *map[string]struct{}) {
	cacheDir, _ := os.UserCacheDir()
	os.MkdirAll(cacheDir + "/gospell", os.ModePerm)
	file := cacheDir + "/gospell/cache.gob"

    f := new(os.File)
    if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
        f, err = os.Create(file)
        if err != nil {
            log.Fatal("Error creating cache file:", err)
        }
        defer f.Close()
    } else {
        f, err = os.OpenFile(file, os.O_RDWR, 0644)
        if err != nil {
            log.Fatal("Error opening cache file:", err)
        }
        defer f.Close()
    }

	enc := gob.NewEncoder(f)
	if err := enc.Encode(cache); err != nil {
        log.Fatal("Error encoding cache:", err)
    }
}

func LoadCache() map[string]struct{} {
	cacheDir, _ := os.UserCacheDir()
	os.MkdirAll(cacheDir + "/gospell", os.ModePerm)
	file := cacheDir + "/gospell/cache.gob"

	var f = new(os.File)
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		f, err = os.Create(file)
        if err != nil {
            log.Fatal("Error creating cache file:", err)
        }
		defer f.Close()

		return make(map[string]struct{})
	} else {
		f, err = os.Open(file)
        if err != nil {
            log.Fatal("Error opening cache file:", err)
        }
		defer f.Close()
	}

	enc := gob.NewDecoder(f)
    cache := make(map[string]struct{})

    if err := enc.Decode(&cache); err != nil {
        fmt.Println("Error decoding cache:", err)
    }
	return cache
}

func PrintDefinition(responseObject Welcome, numDefinitions int) {
	if len(responseObject) == 0 {
		fmt.Println("No definition found.")
		return
	}

	fmt.Println("Definitions:")
    index := 0
	for _, meaning := range responseObject[0].Meanings {
		for _, definitions := range meaning.Definitions {
            if index == numDefinitions {
                return
            }

            fmt.Print(meaning.PartOfSpeech + ": ")
            fmt.Println(definitions.Definition)
            index++
		}
	}
}
