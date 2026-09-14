package pokeapi

// LocationAreaResponse stores detailed information about a single location area,
// including which Pokemon spawn there.
type LocationAreaResponse struct {
	Name              string `json:"name"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}
