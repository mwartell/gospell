package api

import "testing"

func TestGetRandomLineFromWordlist(t *testing.T) {
	word := RandomWord()
	t.Logf("Random word: %v", word)

	if len(word) == 0 {
		t.Errorf("Expected a non-empty word, got empty string")
	}
	if len(word) > 20 {
		t.Errorf("Expected a word of length <= 20, got %d", len(word))
	}
	if word == "a" || word == "I" {
		t.Errorf("Expected a word other than 'a' or 'I', got %s", word)
	}
}

func BenchmarkRandomWord(b *testing.B) {
	var word string
	for b.Loop() {
		word = RandomWord()
	}
	b.Logf("Random word: %v", word)
}

func BenchmarkSplitWords(b *testing.B) {
	for b.Loop() {
		_ = splitWords(wordlist)
	}
}
