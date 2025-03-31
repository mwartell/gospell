package definition

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
)

func LoadCache() Dictionary {
	file := "wordmap.gob"

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
