package main

import (
	"slices"
	"testing"
)

func TestCleanInput(t *testing.T) {
	// Use Table driven tests
	testCases := []struct {
		name      string
		textInput string
		expected  []string
	}{
		// Different test cases
		{name: "Simple", textInput: "This ain't It ", expected: []string{"this", "ain't", "it"}},
		{name: "Too Much Space", textInput: " hello world  ", expected: []string{"hello", "world"}},
		{name: "Complex", textInput: "Come Back, This  IS where YOU belong", expected: []string{"come", "back,", "this", "is", "where", "you", "belong"}},
		{name: "Empty", textInput: "", expected: []string{}},
	}

	for _, test := range testCases {
		// Subtests
		t.Run(test.name, func(t *testing.T) {
			actual := cleanInput(test.textInput)
			if !slices.Equal(actual, test.expected) {
				t.Errorf("cleanInput(%q):\nwanted: %v\ngot: %v", test.textInput, test.expected, actual)
			}
		})
	}
}
