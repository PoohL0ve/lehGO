package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"chirpy/internal/auth"

	"github.com/google/uuid"
)

func (cfg *ApiConfig) HandlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	// Before doing anything ensure resquest is from the correct client/service
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil || apiKey != cfg.PolkaKey {
		RespondWithError(w, http.StatusUnauthorized, "Invalid API key")
		return
	}

	// Polka request hook
	type PolkaWebhook struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"` // Helps Go map to data field
	}

	// Decode request body
	polka := PolkaWebhook{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&polka)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid Request")
		return
	}

	// Filter for user.upgraded
	if polka.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Parse user id
	userPolkaID, err := uuid.Parse(polka.Data.UserID)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid user id")
		return
	}

	// Update db and user status
	_, err = cfg.DB.UpdateToChirpyRed(r.Context(), userPolkaID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			RespondWithError(w, http.StatusNotFound, "User does not exist")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "Unable to update user")
		return
	}

	// Successful response with no content
	w.WriteHeader(http.StatusNoContent)
}
