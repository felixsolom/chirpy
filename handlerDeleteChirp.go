package main

import (
	"net/http"

	"github.com/felixsolom/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpIdStr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIdStr)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid chirp id: %w", err)
		return
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
	chirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp wasn't found in db: %w", err)
		return
	}

	if chirp.UserID.UUID != userID {
		respondWithError(w, http.StatusForbidden, "the user is not the owner of the chirp", nil)
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't delete chirp from db: %w", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}
