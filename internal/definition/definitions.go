package definition

import "fmt"

type State struct {
	cache       Dictionary
	word        string
	index       int
	definitions []string
}

// getDefinitionList returns a list of definitions for a given word from the cache.
// It populates the definitions field in the State struct.
// This function is called internally by GetDefinition to initialize the definitions list.
func (s *State) getDefinitionList() {
	definitions := s.cache[s.word] // O(1) lookup which is so cool
	list := make([]string, 0)

	for _, definition := range definitions {
		list = append(list,
			fmt.Sprintf(
				"(%d of %d) %s: %s",
				definition.DefinitionIndex,
				definition.NumDefinitions,
				definition.PartOfSpeech,
				definition.Definition,
			),
		)
	}

	s.definitions = list // store the definitions in the state
}

// NextDefinition retrieves the next definition of a word from the cache.
// If the user requests a definition past the last one, it returns the last definition.
func (s *State) NextDefinition() string {
	if s.index+1 >= len(s.definitions) { // if user requests something past the end of the definition list
		return s.definitions[len(s.definitions)-1] // return the last definition
	} else { // increment index & change definition
		s.index++
		return s.definitions[s.index]
	}
}

// PrevDefinition retrieves the previous definition of a word from the cache.
// If the user requests a definition before the first one, it returns the first definition.
func (s *State) PrevDefinition() string {
	if s.index-1 < 0 { // if user requests something before the beginning of the definition list
		return s.definitions[0] // return the first definition
	} else { // decrement index & change definition
		s.index--
		return s.definitions[s.index]
	}
}

// GetDefinition retrieves the first definition of a word from the cache.
//
// How To Use:
//
// 1. Create a new State instance.
//
// 2. Call GetDefinition with the word you want to look up.
//
// 3. Use NextDefinition() and PrevDefinition() to navigate through the definitions.
//
// Example:
//
//	state := &definition.State{}
//	firstDef := state.GetDefinition("example") // retrieves the first definition
//	fmt.Println(firstDef) // prints the first definition
//	nextDef := state.NextDefinition() // retrieves the next definition
//	fmt.Println(nextDef) // prints the next definition
//	prevDef := state.PrevDefinition() // retrieves the previous definition
//	fmt.Println(prevDef) // prints the previous definition
func (m *State) GetDefinition(word string) string {
	if m.cache == nil { // Only load the cache if it's not already loaded.
		m.cache = LoadCache()
	}

	m.word = word
	m.index = 0
	m.getDefinitionList() // populate the definitions list

	if len(m.definitions) == 0 {
		panic(fmt.Sprintf("No definitions found for the word: %s", word)) // this SHOULD never happen if the cache is loaded correctly
	}

	return m.definitions[m.index] // return the first definition
}
