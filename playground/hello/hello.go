package main

import "fmt"

const (
	spanish = "Spanish"
	french  = "French"

	// Prefixes
	englishHelloPrefix = "Hello, "
	spanishHelloPrefix = "Hola, "
	frenchHelloPrefix  = "Bonjour, "
)

func Hello(name, language string) string {
	if name == "" {
		name = "World"
	}

	return greetingPrefix(language) + name
}

// Refactor
func greetingPrefix(language string) (prefix string) {
	switch language {
	case french:
		prefix = frenchHelloPrefix
	case spanish:
		prefix = spanishHelloPrefix
	default:
		prefix = englishHelloPrefix
	}
	return
}

// Separation of Concerns: Println may produce side effects
func main() {
	fmt.Println(Hello("Chris", "French"))
	fmt.Println(Hello("Elodie", "Spanish"))
}
