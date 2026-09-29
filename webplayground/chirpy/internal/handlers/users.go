package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"chirpy/internal/auth"
	"chirpy/internal/database"

	"github.com/google/uuid"
)

// User Response struct matching JSON tags
type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

func (cfg *ApiConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	// Data to pass to sql query
	type paramerters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	params := paramerters{}

	// Decode incoming data
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Request is invalid")
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Unable to hash password")
		return
	}
	// Execute db query
	// Extract the context
	ctx := r.Context()
	dbUser, err := cfg.DB.CreateUser(ctx, database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
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

func (cfg *ApiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	req := request{}
	err := decoder.Decode(&req)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	dbUser, err := cfg.DB.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	match, err := auth.CheckPasswordHash(req.Password, dbUser.HashedPassword)
	if err != nil || !match {
		RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	// Generate access token
	accessToken, err := auth.MakeJWT(dbUser.ID, cfg.JWTSecret, time.Hour)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Unable to create access token")
		return
	}

	// Generate refresh token
	refreshTokenString, err := auth.MakeRefreshToken()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Unable to create refresh token")
		return
	}

	// Save refresh token string to db
	refreshToken, err := cfg.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshTokenString,
		UserID:    dbUser.ID,
		ExpiresAt: time.Now().UTC().Add(60 * 24 * time.Hour), // 60 days
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Could not save refresh token")
		return
	}

	RespondWithJSON(w, http.StatusOK, User{
		ID:           dbUser.ID,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
		Email:        dbUser.Email,
		Token:        accessToken,
		RefreshToken: refreshToken.Token,
	})
}
