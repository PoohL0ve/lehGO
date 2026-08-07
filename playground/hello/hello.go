package main

import "fmt"

func Hello(name string) string {
	return "Hello, " + name
}

// Separation of Concerns: Println may produce side effects
func main() {
	fmt.Println(Hello("Chris"))
}
