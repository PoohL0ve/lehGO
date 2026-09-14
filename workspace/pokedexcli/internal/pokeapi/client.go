// Package pokeapi handles low-level HTTP communication with the PokeAPI REST endpoints.
//
// File: client.go
// Purpose: Defines the core PokeAPI Client struct and constructor function.
// It configures connection options (like HTTP timeouts) and acts as the entry point
// for all PokeAPI requests across the application.
package pokeapi

import (
	"net/http"
	"time"

	"github.com/PoohL0ve/lehGo/workspace/pokedexcli/internal/pokecache"
)

// Client wraps a standard Go http.Client to make API requests to PokeAPI.
type Client struct {
	httpClient http.Client
	cache      pokecache.Cache
}

// NewClient creates and returns a pre-configured PokeAPI Client.
// It accepts a timeout duration to ensure network requests do not hang indefinitely.
func NewClient(timeout, interval time.Duration) Client {
	return Client{
		cache: pokecache.NewCache(interval),
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}
