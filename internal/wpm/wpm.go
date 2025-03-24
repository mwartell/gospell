package wpm

import (
	"time"
	"unicode/utf8"
)

var wpm int

func CalculateWpm(word string, initialTime time.Time, finalTime time.Time) int {
	// If the word is empty, return the last calulated wpm
	if word == "" {
		return wpm
	}

	timeElapsed := finalTime.Sub(initialTime)

	// Convert to minutes, and handle very small time intervals
	elapsedMinutes := timeElapsed.Minutes()
	if elapsedMinutes < 0.001 { // Avoid division by very small numbers
		return 0
	}

	numWords := float64(utf8.RuneCountInString(word)) / 5.0

	wpm = int(numWords / elapsedMinutes)
	return wpm
}
