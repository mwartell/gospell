package api

import (
	_ "embed"
	"math/rand"
	"strings"
)

//go:embed wordlist.txt
var wordlist string

var words = splitWords(wordlist)

func splitWords(wordlist string) []string {
	return strings.Split(wordlist, "\n")
}

// RandomWord returns a random word from the embedded word list
func RandomWord() string {
	i := rand.Intn(len(words))
	return words[i]
}
