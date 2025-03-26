package definition

import (
	"io"
	"log"
	"net/http"
)

var responseObject Welcome
var UseRealVoices = false

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

func getURLS() []string {
	audioURLs := make([]string, 0)

	for _, phonetic := range responseObject[0].Phonetics {
		if phonetic.Audio != "" {
			audioURLs = append(audioURLs, phonetic.Audio)
		}
	}

	return audioURLs
}

func GetResponseObject(word string) Welcome {
	if len(responseObject) == 0 || responseObject[0].Word != word {
		getResponse(word)
	}

	return responseObject
}

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

	if UseRealVoices {
		if len(responseObject) == 0 || doesNotContainAudio(responseObject) {
			addToCache(word, wordsWithoutDefinitions)
			return false
		} else {
			return true
		}
	} else {
		if len(responseObject) == 0 {
			addToCache(word, wordsWithoutDefinitions)
			return false
		} else {
			return true
		}
	}
}

func doesNotContainAudio(res Welcome) bool {
	for _, phonetic := range res[0].Phonetics {
		if phonetic.Audio != "" {
			return false
		}
	}
	return true
}
