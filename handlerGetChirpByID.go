package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirp(w http.ResponseWriter, r *http.Request) {
	type Chirp struct {
		ID        uuid.UUID     `json:"id"`
		CreatedAt time.Time     `json:"created_at"`
		UpdatedAt time.Time     `json:"updated_at"`
		Body      string        `json:"body"`
		UserID    uuid.NullUUID `json:"user_id"`
	}

	path := r.PathValue("chirpID")
	if path == "" {
		respondWithError(w, http.StatusNotFound, "Chirp id wasn't provided", nil)
		return
	}

	u, err := uuid.Parse(path)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid UUID", err)
	}

	chirp, err := cfg.db.GetChirpByID(r.Context(), u)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp wasn't found in database", err)
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
