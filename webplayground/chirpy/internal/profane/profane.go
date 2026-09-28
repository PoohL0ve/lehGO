package profane

import (
	"strings"
)

var badWords []string = []string{
	"kerfuffle",
	"sharbert",
	"fornax",
}

const charString = "****"

func SubstituteWords(words string) string {
	if len(words) == 0 {
		return ""
	}

	wordSlice := strings.Split(words, " ")
	for pos, word := range wordSlice {
		for _, bad := range badWords {
			if bad == strings.ToLower(word) {
				wordSlice[pos] = charString
				break
			}
		}
	}

	return strings.Join(wordSlice, " ")
}
