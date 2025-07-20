package main

import (
	"net/http"
	"time"

	"github.com/felixsolom/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerAllChirps(w http.ResponseWriter, r *http.Request) {
	type Chirp struct {
		ID        uuid.UUID  `json:"id"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt time.Time  `json:"updated_at"`
		Body      string     `json:"body"`
		UserID    *uuid.UUID `json:"user_id"`
	}

	var chirps []database.Chirp
	var err error

	authorIDStr := r.URL.Query().Get("author_id")

	if authorIDStr != "" {
		id, parseErr := uuid.Parse(authorIDStr)
		if parseErr != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author_id format", err)
			return
		}
		userID := uuid.NullUUID{
			UUID:  id,
			Valid: true,
		}
		chirps, err = cfg.db.GetChirpsByUserID(r.Context(), userID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps for this user", err)
			return
		}

	} else {
		chirps, err = cfg.db.RetrieveAllChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve all chirps", err)
			return
		}
	}

	out := make([]Chirp, len(chirps))
	for i, dbChirp := range chirps {
		out[i] = Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
		}
		if dbChirp.UserID.Valid {
			out[i].UserID = &dbChirp.UserID.UUID
		}
	}
	respondWithJSON(w, http.StatusOK, out)
}
