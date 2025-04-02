package api

import (
	_ "embed"
	"hash/maphash"
	"math/rand"
	"strings"
)

//go:embed wordlist.txt
var fileString string
var file []string = splitWords(fileString)
var rng = NewRand()

// RandomWord returns a random word from the wordlist.
func RandomWord() string {
	randomNumber := rng.Intn(len(file))
	return file[randomNumber]
}

// Rand64 returns a pseudo-random uint64. It can be used concurrently and is lock-free.
// Effectively, it calls runtime.fastrand.
func Rand64() uint64 {
	return new(maphash.Hash).Sum64()
}

// NewRand returns a properly seeded *rand.Rand. It has *slightly* higher overhead than
// Rand64 (as it has to allocate), but the resulting PRNG can be re-used to offset that cost.
func NewRand() *rand.Rand {
	return rand.New(rand.NewSource(int64(Rand64())))
}

func splitWords(wordlist string) []string {
	return strings.Split(wordlist, "\n")
}
