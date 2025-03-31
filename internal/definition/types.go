package definition

// Dictionary represents the entire JSON structure
// The JSON is a map where keys are words and values are arrays of entries
type Dictionary map[string][]Entry

// Entry represents a single word definition
type Entry struct {
	Word            string `json:"word"`
	DefinitionIndex int64  `json:"definition_index"`
	NumDefinitions  int64  `json:"num_definitions"`
	PartOfSpeech    string `json:"part_of_speech"`
	Definition      string `json:"definition"`
}
