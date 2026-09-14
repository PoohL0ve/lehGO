// Package main lives in the root directory and ties together CLI execution.
//
// File: command_pokedex.go
// Purpose: Implements the 'pokedex' command, displaying a list of all
// Pokemon stored in local memory.
package main

import (
	"fmt"
)

func commandPokedex(cfg *config, args []string) error {
	if len(cfg.CaughtPokemon) == 0 {
		fmt.Println("Your Pokedex is empty. Go catch some Pokemon!")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for _, pokemon := range cfg.CaughtPokemon {
		fmt.Printf(" - %s\n", pokemon.Name)
	}

	return nil
}
