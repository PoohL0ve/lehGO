// Package main lives in the root directory and ties together CLI execution.
//
// File: command_map.go
// Purpose: Implements the 'map' and 'mapb' command handlers. These functions handle
// user interaction, updating pagination state pointers in config, and displaying location names.
package main

import (
	"fmt"
)

// commandMap fetches and displays the next 20 location areas.
func commandMap(cfg *config, args []string) error {
	// 1. Fetch location areas using the URL pointer currently saved in config
	locationResp, err := cfg.PokeapiClient.ListLocations(cfg.NextLocationUrl)
	if err != nil {
		return fmt.Errorf("List locations: %w", err)
	}

	// 2. Save updated pagination URL pointers back to config
	if locationResp.Next != nil {
		cfg.NextLocationUrl = locationResp.Next
	} else {
		cfg.NextLocationUrl = nil
	}

	if locationResp.Previous != nil {
		cfg.PrevLocationUrl = locationResp.Previous
	} else {
		cfg.PrevLocationUrl = nil
	}

	// 3. Print the names of the locations returned
	for _, location := range locationResp.Results {
		fmt.Println(location.Name)
	}

	return nil
}

// commandMapb fetches and displays the previous 20 location areas.
func commandMapb(cfg *config, arg []string) error {
	// Guard Check: If prevLocationsURL is nil, we are on the first page
	if cfg.PrevLocationUrl == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	// 1. Fetch location areas using the previous URL pointer from config
	locationResp, err := cfg.PokeapiClient.ListLocations(cfg.PrevLocationUrl)
	if err != nil {
		return err
	}

	// 2. Save updated pagination URL pointers back to config
	cfg.NextLocationUrl = locationResp.Next
	cfg.PrevLocationUrl = locationResp.Previous

	// 3. Print the names of the locations returned
	for _, location := range locationResp.Results {
		fmt.Println(location.Name)
	}

	return nil
}
