package api
// PACKAGE TO BE DEPRECATED

import (
	"strings"
	"unicode"

	"github.com/jharlan-hash/gospell/internal/definition"

	"github.com/tjarratt/babble"
)

func GetAcceptableWord(babbler babble.Babbler, cache map[string]struct{}) string {
	for {
		word := babbler.Babble()
		if isAcceptableWord(word, cache) {
			return word
		}
	}
}

// Acceptable words are lowercase and contain no special characters & are defined in the dictionary
func isAcceptableWord(word string, cache map[string]struct{}) bool {
	return !strings.ContainsAny(word, "-_'") && isLower(&word) && definition.IsDefined(word, cache)
}

func isLower(s *string) bool {
	for _, r := range *s {
		if !unicode.IsLower(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

