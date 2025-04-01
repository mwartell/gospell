package definition

import (
	"embed"
	"encoding/gob"
	"fmt"
	"log"
)

//go:embed wordmap.gob
var fs embed.FS

// LoadCache loads the wordmap cache from a file.
// It returns a Dictionary, which is a map of words to an Entry slice.
func LoadCache() Dictionary {
	f, err := fs.Open("wordmap.gob")
	if err != nil {
		log.Fatal(err)
	}

	enc := gob.NewDecoder(f)
	cache := make(Dictionary)

	if err := enc.Decode(&cache); err != nil {
		fmt.Println("Error decoding cache:", err)
	}

	return cache
}
