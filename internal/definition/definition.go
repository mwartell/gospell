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

var responseObject Welcome

func GetResponseObject(word string) Welcome {
	if len(responseObject) == 0 || responseObject[0].Word != word {
		getResponse(word)
	}

	return responseObject
}

func GetFirstDefinition(res Welcome) string {
	slice := GetDefinitionSlice(res)
	slice[0] = fmt.Sprintf("(1 of %d) %s", len(slice), slice[0])
	return slice[0]
}

// func NextDefinition(definition *string, word string, index *int) {
// 	if len(responseObject) == 0 || responseObject[0].Word != word {
// 		getResponse(word)
// 	}
//
// 	definitionSlice := GetDefinitionSlice(responseObject)
//
// 	*index++
//
// 	if *index >= len(definitionSlice) {
// 		*definition = fmt.Sprintf("(%d of %d) %s", *index, len(definitionSlice), definitionSlice[0])
// 		*index--
// 	} else {
// 		*definition = fmt.Sprintf("(%d of %d) %s", *index + 1, len(definitionSlice), definitionSlice[*index])
// 	}
// }

func NextDefinition(definition *string, index *int) {
	if len(responseObject) == 0 {
		panic("everything has gone wrong and response object is empty")
	}

	definitionSlice := GetDefinitionSlice(responseObject)

	if *index+1 >= len(definitionSlice) { // if user requests something past the end of the definition list
		*definition = fmt.Sprintf(
			"(%d of %d) %s",
			len(definitionSlice),
			len(definitionSlice),
			definitionSlice[len(definitionSlice)-1],
		)
		return
	} else { // increment index & change definition
		*index++
		*definition = fmt.Sprintf(
			"(%d of %d) %s",
			*index+1,
			len(definitionSlice),
			definitionSlice[*index],
		)
		return
	}

}

func PrevDefinition(definition *string, word string, index *int) {
	if len(responseObject) == 0 {
		panic("everything has gone wrong and response object is empty")
	}

	definitionSlice := GetDefinitionSlice(responseObject)

	if *index <= 0 { // if the user requests something before the start of the list
		*index = 0
		*definition = fmt.Sprintf(
			"(%d of %d) %s",
			1,
			len(definitionSlice),
			definitionSlice[0],
		)
	} else {
		*index--
		*definition = fmt.Sprintf(
			"(%d of %d) %s",
			*index+1,
			len(definitionSlice),
			definitionSlice[*index],
		)
	}
}

func getResponse(word string) {
	response, err := http.Get("https://api.dictionaryapi.dev/api/v2/entries/en/" + word)
	if err != nil { // TODO: handle this better
		log.Fatal(err)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil { // TODO: handle this better
		log.Fatal(err)
	}
	responseObject, _ = UnmarshalWelcome(responseData)
}

// isDefined checks if a word is defined in the dictionary
func IsDefined(word string, wordsWithoutDefinitions map[string]struct{}) bool {
	// first, search the cache
	// if found, return false
	// if not found, proceed to API call

	if wordInCache(word, wordsWithoutDefinitions) {
		return false
	}

	if len(responseObject) == 0 || responseObject[0].Word != word {
		getResponse(word)
	}

	if len(responseObject) == 0 {
		addToCache(word, wordsWithoutDefinitions)
		return false
	} else {
		return true
	}
}

func addToCache(word string, wordsWithoutDefinitions map[string]struct{}) {
	(wordsWithoutDefinitions)[word] = struct{}{}
}

func wordInCache(word string, wordsWithoutDefinitions map[string]struct{}) bool {
	if _, exists := (wordsWithoutDefinitions)[word]; exists {
		return true
	} else {
		return false
	}
}

func SaveCache(cache *map[string]struct{}) {
	cacheDir, _ := os.UserCacheDir()
	os.MkdirAll(cacheDir+"/gospell", os.ModePerm)
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
	os.MkdirAll(cacheDir+"/gospell", os.ModePerm)
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

func GetDefinitionSlice(resposeObject Welcome) []string {
	definitionSlice := make([]string, 0)

	for _, meaning := range resposeObject[0].Meanings {
		for _, definitions := range meaning.Definitions {
			definitionString := fmt.Sprintf("%s: %s", meaning.PartOfSpeech, definitions.Definition)
			definitionSlice = append(definitionSlice, definitionString)
		}
	}
	return definitionSlice
}

func GetDefinition(word string) string {
	getResponse(word)

	return GetFirstDefinition(responseObject)
}
