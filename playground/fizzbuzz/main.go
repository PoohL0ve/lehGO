package main

import "fmt"

func fizzbuzz() {
	for i := 1; i <= 100; i++ {
		// Calculate modulo checks ONCE per iteration
		isFizz := i%3 == 0
		isBuzz := i%5 == 0

		// Tagless switch checks conditions in order
		switch {
		case isFizz && isBuzz:
			fmt.Println("fizzbuzz")
		case isFizz:
			fmt.Println("fizz")
		case isBuzz:
			fmt.Println("buzz")
		default:
			fmt.Println(i)
		}
	}
}

// don't touch below this line

func main() {
	fizzbuzz()
}
