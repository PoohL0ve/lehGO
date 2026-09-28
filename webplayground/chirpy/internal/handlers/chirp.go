package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"chirpy/internal/database"
	"chirpy/internal/profane"

	"github.com/google/uuid"
)

// JSON body to decode: Expected incoming body
type parameters struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

// Defines outgoing JSON response
type returnVals struct {
	CleanedBody string `json:"cleaned_body"`
}

// Chirp: JSON data to return
type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

// Chirp length and filters profane words
func (cfg *ApiConfig) HandlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	const maxChirpLength = 140
	params := parameters{}

	// Decode message
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	if len(params.Body) > maxChirpLength {
		RespondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedBody := profane.SubstituteWords(params.Body)

	// Create a chirp
	dbChirp, err := cfg.DB.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: params.UserID,
	})

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Could not create chirp")
		return
	}

	// Object to return
	response := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	RespondWithJSON(w, http.StatusCreated, response)
}

// Retrieves all the chirps
func (cfg *ApiConfig) HandlerGetChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Unable to retrieve chirps")
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		})
	}

	RespondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *ApiConfig) HandlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid chirp ID format")
		return
	}

	dbChirp, err := cfg.DB.GetChirp(r.Context(), chirpID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}

	RespondWithJSON(w, http.StatusOK, Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	})
}
