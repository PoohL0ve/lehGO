// Package pokeapi handles low-level HTTP communication with the PokeAPI REST endpoints.
//
// File: location_list.go
// Purpose: Defines the data structures and API method required to fetch paginated
// location-area data from the PokeAPI (/api/v2/location-area).
package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ListLocations fetches a page of location areas, serving from cache if available.
func (c *Client) ListLocations(pageURL *string) (LocationAreasResponse, error) {
	url := "https://pokeapi.co/api/v2/location-area"
	if pageURL != nil {
		url = *pageURL
	}

	// 1. Check Cache first (Cache Hit)
	if val, ok := c.cache.Get(url); ok {
		locationAreasResp := LocationAreasResponse{}
		err := json.Unmarshal(val, &locationAreasResp)
		if err != nil {
			return LocationAreasResponse{}, err
		}
		return locationAreasResp, nil
	}

	// 2. Cache Miss: Make HTTP Request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LocationAreasResponse{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return LocationAreasResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return LocationAreasResponse{}, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationAreasResponse{}, err
	}

	// 3. Save raw bytes to Cache for next time
	c.cache.Add(url, data)

	// 4. Unmarshal raw bytes into Go struct
	locationAreasResp := LocationAreasResponse{}
	err = json.Unmarshal(data, &locationAreasResp)
	if err != nil {
		return LocationAreasResponse{}, err
	}

	return locationAreasResp, nil
}
