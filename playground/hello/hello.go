package main

import "fmt"

const englishHelloPrefix = "Hello, "

func Hello(name string) string {
	if name == "" {
		name = "World"
	}
	return englishHelloPrefix + name
}

// Separation of Concerns: Println may produce side effects
func main() {
	fmt.Println(Hello("Chris"))
}
