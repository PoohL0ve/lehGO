// Package main lives in the root directory and ties together CLI execution.
//
// File: command_inspect.go
// Purpose: Implements the 'inspect' command, retrieving stored Pokemon details
// from local memory if the user has previously caught it.
package main

import (
	"errors"
	"fmt"
)

func commandInspect(cfg *config, args []string) error {
	// Guard Check: Ensure the user provided a pokemon name
	if len(args) < 1 {
		return errors.New("you must provide a pokemon name")
	}

	pokemonName := args[0]

	// 1. Look up Pokemon in local state map
	pokemon, exists := cfg.CaughtPokemon[pokemonName]
	if !exists {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	// 2. Format and print Pokemon details
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)

	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, typeInfo := range pokemon.Types {
		fmt.Printf("  - %s\n", typeInfo.Type.Name)
	}

	return nil
}
