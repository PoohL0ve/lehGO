// Package main lives in the root directory and ties together CLI execution.
//
// File: command_catch.go
// Purpose: Implements the 'catch' command, fetching Pokemon stats and using
// a random threshold based on base_experience to determine catch success.
package main

import (
	"errors"
	"fmt"
	"math/rand"
)

func commandCatch(cfg *config, args []string) error {
	if len(args) < 1 {
		return errors.New("you must provide a pokemon name")
	}

	pokemonName := args[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	pokemon, err := cfg.PokeapiClient.GetPokemon(pokemonName)
	if err != nil {
		return err
	}

	// Calculate catch probability based on BaseExperience.
	// Higher BaseExperience means lower probability of catching.
	res := rand.Intn(pokemon.BaseExperience)

	// Threshold check: if random roll is low enough, catch succeeds
	if res > 40 {
		fmt.Printf("%s escaped!\n", pokemonName)
		return nil
	}

	fmt.Printf("%s was caught!\n", pokemonName)
	cfg.CaughtPokemon[pokemonName] = pokemon
	return nil
}
