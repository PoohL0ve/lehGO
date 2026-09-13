package main

import (
	"bufio"
	"fmt"
	"os"
)

// Composite structure to define commands
type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	commands := getCommands()
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}

	fmt.Println()
	return nil
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex Cli",
			callback:    commandExit,
		},
	}
}

func main() {
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

		// Check for existing command
		command, exists := getCommands()[inform]
		if exists {
			err := command.callback()
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
