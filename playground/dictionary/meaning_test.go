package dictionary

import (
	"errors"
	"testing"
)

func TestSearch(t *testing.T) {
	dictionary := Dictionary{"test": "this is just a test"}

	t.Run("known word", func(t *testing.T) {
		got, _ := dictionary.Search("test")
		want := "this is just a test"

		assertStrings(t, got, want)
	})

	t.Run("unknown word", func(t *testing.T) {
		_, err := dictionary.Search("unknown")

		if err == nil {
			t.Fatal("expected to get an error")
		}

		assertError(t, err, ErrWordNotFound)

	})
}

// Test that checks a word being added
func TestAdd(t *testing.T) {
	// Declare dictionary here
	dictionary := Dictionary{}

	suite := []struct {
		name       string
		word       string
		definition string
		want       error
	}{
		{name: "add word", word: "learn", definition: "to gain knowledge and understanding", want: nil},
		{name: "existing word", word: "learn", definition: "copy a skill", want: ErrWordExists},
		{name: "empty word", word: "", definition: "some meaning", want: ErrEmpty},
		{name: "empty definition", word: "crime", definition: "", want: ErrEmpty},
		{name: "both empty", word: "", definition: "", want: ErrEmpty},
	}

	for _, test := range suite {
		t.Run(test.name, func(t *testing.T) {
			err := dictionary.Add(test.word, test.definition)
			assertError(t, err, test.want)

			if test.want == nil {
				assertDefinition(t, dictionary, test.word, test.definition)
			}

		})
	}
}

func TestUpdate(t *testing.T) {
	// Define initial dictionary state for our tests
	dictionary := Dictionary{"casper": "a friendly ghost"}

	suite := []struct {
		name       string
		word       string
		definition string
		want       error
	}{
		{
			name:       "update existing word",
			word:       "casper",
			definition: "a miserable ghost",
			want:       nil,
		},
		{
			name:       "update non-existent word",
			word:       "missing",
			definition: "some definition",
			want:       ErrWordDoesNotExists,
		},
		{
			name:       "blank search",
			word:       "",
			definition: "some meaning",
			want:       ErrEmpty,
		},
	}

	for _, test := range suite {
		t.Run(test.name, func(t *testing.T) {
			err := dictionary.Update(test.word, test.definition)

			// Assert that the error matches what this specific test case expects
			assertError(t, err, test.want)

			// If the update was expected to succeed, verify the new definition took effect
			if test.want == nil {
				assertDefinition(t, dictionary, test.word, test.definition)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	// All subtests share the same state
	// Define in for loop, for each to have its own
	dictionary := Dictionary{"hope": "a journey of heart break and wishful desire"}

	tests := []struct {
		name string
		word string
		want error
	}{
		{name: "word exists", word: "hope", want: nil},
		{name: "word does not exists", word: "job", want: ErrWordDoesNotExists},
		{name: "empty delete", word: "", want: ErrEmpty},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := dictionary.Delete(test.word)

			assertError(t, err, test.want)

			if test.want == nil {
				_, err := dictionary.Search(test.word)
				assertError(t, err, ErrWordNotFound)
			}
		})
	}
}

func assertDefinition(t testing.TB, dictionary Dictionary, word, definition string) {
	t.Helper()

	got, err := dictionary.Search(word)
	if err != nil {
		t.Fatal("should find added word:", err)
	}
	assertStrings(t, got, definition)
}

func assertStrings(t testing.TB, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got: %s wanted %s", got, want)
	}
}

func assertError(t testing.TB, got, want error) {
	t.Helper()

	if !errors.Is(got, want) {
		t.Errorf("got error: %s, want error: %s", got, want)
	}
}
