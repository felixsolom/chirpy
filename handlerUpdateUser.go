package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/felixsolom/chirpy/internal/auth"
	"github.com/felixsolom/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type User struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	tokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Authentication required: malformed token", err)
		return
	}

	userID, err := auth.ValidateJWT(tokenStr, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Authentication required: invalid token", err)
		return
	}

	var params parameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode JSON: %v", err)
		return
	}

	if params.Email == "" || params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Please fill both email and password fields", nil)
		return
	}

	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash the password: %w", err)
		return
	}

	paramsToUpdate := database.UpdateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPass,
		ID:             userID,
	}

	user, err := cfg.db.UpdateUser(r.Context(), paramsToUpdate)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user credentials: %w", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
