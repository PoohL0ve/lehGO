package handlers

import (
	"chirpy/internal/auth"
	"net/http"
	"time"
)

func (cfg *ApiConfig) HandlerRefresh(w http.ResponseWriter, r *http.Request) {
	// Response struct to return
	type response struct {
		Token string `json:"token"`
	}

	// Extract token from header
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Could not get token")
		return
	}

	// Locate user by refresh token
	user, err := cfg.DB.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Unable to locate user")
		return
	}

	// Check if refresh token has expired
	if time.Now().UTC().After(user.ExpiresAt) {
		RespondWithError(w, http.StatusUnauthorized, "Refresh token has expired")
		return
	}

	// Check if the token has been revoked
	if user.RevokedAt.Valid {
		RespondWithError(w, http.StatusUnauthorized, "Refresh token has been revoked")
		return
	}

	// If all checks are cleared issue a new access token
	accessToken, err := auth.MakeJWT(user.ID, cfg.JWTSecret, time.Hour)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Unable to create an access token")
		return
	}

	RespondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})

}

func (cfg *ApiConfig) HandlerRevoke(w http.ResponseWriter, r *http.Request) {
	// Extract refresh token
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Unable to locate refresh token")
		return
	}

	// Revoke refresh token
	_, err = cfg.DB.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Unable to revoke token")
		return
	}

	// Respond with 204 No Content (no body)
	w.WriteHeader(http.StatusNoContent)
}
