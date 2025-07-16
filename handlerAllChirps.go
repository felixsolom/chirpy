package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerAllChirps(w http.ResponseWriter, r *http.Request) {
	type Chirp struct {
		ID        uuid.UUID     `json:"id"`
		CreatedAt time.Time     `json:"created_at"`
		UpdatedAt time.Time     `json:"updated_at"`
		Body      string        `json:"body"`
		UserID    uuid.NullUUID `json:"user_id"`
	}

	chirps, err := cfg.db.RetrieveAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve all chirps", err)
		return
	}

	out := make([]Chirp, 0, len(chirps))
	for _, chirp := range chirps {
		out = append(out, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, out)
}
