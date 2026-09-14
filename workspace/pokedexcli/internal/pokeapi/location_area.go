package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GetLocationArea fetches detailed information about a specific location area,
// including which Pokemon can be found there.
func (c *Client) GetLocationArea(areaName string) (LocationAreaResponse, error) {
	url := "https://pokeapi.co/api/v2/location-area/" + areaName

	// 1. Check Cache
	if val, ok := c.cache.Get(url); ok {
		areaResp := LocationAreaResponse{}
		err := json.Unmarshal(val, &areaResp)
		if err != nil {
			return LocationAreaResponse{}, err
		}
		return areaResp, nil
	}

	// 2. Fetch from API on Cache Miss
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return LocationAreaResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return LocationAreaResponse{}, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	// 3. Save to Cache
	c.cache.Add(url, data)

	// 4. Unmarshal
	areaResp := LocationAreaResponse{}
	err = json.Unmarshal(data, &areaResp)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	return areaResp, nil
}
