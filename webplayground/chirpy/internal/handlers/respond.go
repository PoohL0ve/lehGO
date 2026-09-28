package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

// Represents a standard json error payload
type failure struct {
	Error string `json:"error"`
}

// Converts a struct to JSON and writes a response
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// Convert payload (struct to JSON bytes)
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(bodyBytes)

}

// Calls RespondWithJSON to send an error response
func RespondWithError(w http.ResponseWriter, code int, msg string) {
	RespondWithJSON(w, code, failure{
		Error: msg,
	})
}
