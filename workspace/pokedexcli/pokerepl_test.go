package main

import (
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
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			actual := cleanInput(test.textInput)
			if len(actual) != len(test.expected) {
				t.Errorf("Wanted: %v but got %v", test.expected, actual)
			}
			for i := range actual {
				word := actual[i]
				expectedWord := test.expected[i]

				if word != expectedWord {
					t.Fatalf("Needed: %s but Got %s", expectedWord, word)
				}
			}
		})
	}
}
