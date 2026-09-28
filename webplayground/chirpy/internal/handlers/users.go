package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// User Response struct matching JSON tags
type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *ApiConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	// Data to pass to sql query
	type paramerters struct {
		Email string `json:"email"`
	}
	params := paramerters{}

	// Decode incoming data
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Request is invalid")
	}

	// Execute db query
	// Extract the context
	ctx := r.Context()
	dbUser, err := cfg.DB.CreateUser(ctx, params.Email)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "User could not be created")
	}

	// Map model to response struct
	responseUser := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	RespondWithJSON(w, http.StatusCreated, responseUser)
}
