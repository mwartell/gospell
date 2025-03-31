package definition

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
)

// LoadCache loads the wordmap cache from a file.
// It returns a Dictionary, which is a map of words to their definitions.
func LoadCache() Dictionary {
	file := "texts/wordmap.gob"

    f, err := os.Open(file)
	if err != nil {
		log.Fatal("Error opening cache file:", err)
	}
	defer f.Close()

	enc := gob.NewDecoder(f)
	cache := make(Dictionary)

	if err := enc.Decode(&cache); err != nil {
		fmt.Println("Error decoding cache:", err)
	}

	return cache
}
