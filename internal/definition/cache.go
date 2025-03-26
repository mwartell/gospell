package definition

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"os"
)

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
