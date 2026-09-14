package main

import (
	"time"

	"github.com/PoohL0ve/lehGo/workspace/pokedexcli/internal/pokeapi"
)

type config struct {
	PokeapiClient   pokeapi.Client
	NextLocationUrl *string
	PrevLocationUrl *string
	CaughtPokemon   map[string]pokeapi.Pokemon
}

func initialConfig() *config {
	return &config{
		PokeapiClient: pokeapi.NewClient(5*time.Second, 5*time.Minute),
		CaughtPokemon: make(map[string]pokeapi.Pokemon), // 👈 Initialize map
	}
}
