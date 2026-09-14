package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	cfg := initialConfig() // Initialize shared config state
	scanner := bufio.NewScanner(os.Stdin)

	// Print initial prompt before waiting for input
	fmt.Print("Pokedex > ")

	for scanner.Scan() {
		text := scanner.Text()
		cleanText := cleanInput(text)

		// Check slice length to prevent index out of range panic
		if len(cleanText) == 0 {
			continue
		}

		// Access the first word cleanText[0] safely
		inform := cleanText[0]
		args := []string{}
		if len(cleanText) > 1 {
			args = cleanText[1:] // Extract arguments after command
		}

		// Check for existing command
		command, exists := getCommands(cfg)[inform]
		if exists {
			err := command.callback(cfg, args)
			if err != nil {
				fmt.Println(err)
			}
			continue
		} else {
			fmt.Println("Unknown command")
			continue
		}
	}

	// Check scanner error once AFTER the loop exits
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	// To use cleanInput() run go .
}
