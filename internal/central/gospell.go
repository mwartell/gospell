package central

import (
	"gospell/internal/api"
	"gospell/internal/definition"
	"strings"

	"github.com/tjarratt/babble"
)

var wordsWithoutDefinitions = make(map[string]struct{})

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
