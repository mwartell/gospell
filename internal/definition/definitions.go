package definition

import "fmt"

func GetDefinitionSlice(resposeObject Welcome) []string {
	definitionSlice := make([]string, 0)

	for _, meaning := range resposeObject[0].Meanings {
		for _, definitions := range meaning.Definitions {
			definitionString := fmt.Sprintf("%s: %s", meaning.PartOfSpeech, definitions.Definition)
			definitionSlice = append(definitionSlice, definitionString)
		}
	}
	return definitionSlice
}

func GetFirstDefinition(res Welcome) string {
	slice := GetDefinitionSlice(res)
	slice[0] = fmt.Sprintf("(1 of %d) %s", len(slice), slice[0])
	return slice[0]
}

func NextDefinition(definition *string, index *int) {
	if len(responseObject) == 0 {
		return
	}

	definitionSlice := GetDefinitionSlice(responseObject)

	if *index+1 >= len(definitionSlice) { // if user requests something past the end of the definition list
		*definition = fmt.Sprintf(
			"(%d of %d) %s",
			len(definitionSlice),
			len(definitionSlice),
			definitionSlice[len(definitionSlice)-1],
		)
		return
	} else { // increment index & change definition
		*index++
		*definition = fmt.Sprintf(
			"(%d of %d) %s",
			*index+1,
			len(definitionSlice),
			definitionSlice[*index],
		)
		return
	}
}

func PrevDefinition(definition *string, word string, index *int) {
	if len(responseObject) == 0 {
		return
	}

	definitionSlice := GetDefinitionSlice(responseObject)

	if *index <= 0 { // if the user requests something before the start of the list
		*index = 0
		*definition = fmt.Sprintf(
			"(%d of %d) %s",
			1,
			len(definitionSlice),
			definitionSlice[0],
		)
	} else {
		*index--
		*definition = fmt.Sprintf(
			"(%d of %d) %s",
			*index+1,
			len(definitionSlice),
			definitionSlice[*index],
		)
	}
}

func GetDefinition(word string) string {
	getResponse(word)

	return GetFirstDefinition(responseObject)
}
