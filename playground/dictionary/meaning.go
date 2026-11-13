package dictionary

//"errors"

type Dictionary map[string]string

// Allows errors to be immatubale and reusable
type DictionaryErr string

const (
	ErrWordNotFound      = DictionaryErr("unable to find the word you are looking for")
	ErrWordExists        = DictionaryErr("cannot add would because it exists")
	ErrEmpty             = DictionaryErr("word or meaning cannot be empty")
	ErrWordDoesNotExists = DictionaryErr("cannot update a word that does not exists")
)

func (e DictionaryErr) Error() string {
	return string(e)
}

// Search: Accepts two arguments the dictionary and search term.
// Returns the meaning of the search (word) term.
func (d Dictionary) Search(word string) (string, error) {
	definition, ok := d[word]
	if !ok {
		return "", ErrWordNotFound
	}

	return definition, nil
}

// Add a word and its meaning to a dictionary
func (d Dictionary) Add(word, meaning string) error {
	if word == "" || meaning == "" {
		return ErrEmpty
	}
	_, err := d.Search(word)

	switch err {
	// If word is not found add it
	case ErrWordNotFound:
		d[word] = meaning
	case nil:
		return ErrWordExists
	default:
		return err
	}
	return nil
}

func (d Dictionary) Update(word, definition string) error {
	if word == "" {
		return ErrEmpty
	}

	if _, ok := d[word]; ok {
		d[word] = definition
		return nil
	}

	return ErrWordDoesNotExists
}

func (d Dictionary) Delete(word string) error {
	if word == "" {
		return ErrEmpty
	}

	if _, ok := d[word]; ok {
		delete(d, word)
		return nil
	}

	return ErrWordDoesNotExists
}
