package profane

import (
	"fmt"
	"testing"
)

func TestSubstituteWords(t *testing.T) {
	// Using Table Driven Testing
	wordTests := []struct {
		name     string
		text     string
		expected string
	}{
		{name: "Simple", text: "This is a kerfuffle opinion I need to share with the world", expected: "This is a **** opinion I need to share with the world"},
		{name: "Capital", text: "Check out this Sharbert", expected: "Check out this ****"},
		{name: "Empty", text: "", expected: ""},
	}

	for _, tt := range wordTests {
		t.Run(tt.name, func(t *testing.T) {
			actual := SubstituteWords(tt.text)
			if actual != tt.expected {
				t.Errorf("SubstituteWords(%q) = %q, want %q", tt.text, actual, tt.expected)
			}
		})
	}
}

func ExampleSubstituteWords() {
	phrase := SubstituteWords("What the Fornax")
	fmt.Println(phrase)
	// Output: What the ****
}
