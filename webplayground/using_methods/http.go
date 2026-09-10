package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func updateUser(baseURL, id, apiKey string, data User) (User, error) {
	fullURL := baseURL + "/" + id

	// Encode user data
	userBytes, err := json.Marshal(data)
	if err != nil {
		return User{}, err
	}

	// Create a new request
	request, reqErr := http.NewRequest("PUT", fullURL, bytes.NewBuffer(userBytes))
	if reqErr != nil {
		return User{}, err
	}

	// Modify Header
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-API-Key", apiKey)

	// Make the request
	client := &http.Client{}
	response, resErr := client.Do(request)
	if resErr != nil {
		return User{}, err
	}
	defer response.Body.Close()

	var user User
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&user); err != nil {
		return User{}, err
	}

	return user, nil

}

func getUserById(baseURL, id, apiKey string) (User, error) {
	fullURL := baseURL + "/" + id

	newReq, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return User{}, err
	}

	// Modify Headers
	newReq.Header.Set("X-API-Key", apiKey)

	// Make request
	client := &http.Client{}
	res, err := client.Do(newReq)
	if err != nil {
		return User{}, err
	}

	defer res.Body.Close()

	// Decode the response
	var newUser User
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&newUser); err != nil {
		return User{}, err
	}
	return newUser, nil
}
