package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/felixsolom/chirpy/internal/auth"
	"github.com/felixsolom/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type Chirp struct {
		ID        uuid.UUID     `json:"id"`
		CreatedAt time.Time     `json:"created_at"`
		UpdatedAt time.Time     `json:"updated_at"`
		Body      string        `json:"body"`
		UserID    uuid.NullUUID `json:"user_id"`
	}

	var in struct {
		Body   string     `json:"body"`
		UserID *uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&in)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	//authorization
	var userID uuid.NullUUID
	if tokenStr, err := auth.GetBearerToken(r.Header); err == nil {
		if id, err := auth.ValidateJWT(tokenStr, cfg.secret); err == nil {
			userID = uuid.NullUUID{UUID: id, Valid: true}
		}
	}

	//allowed length of chirp
	if len(in.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	cleaned := wordCleanUp(in.Body)
	params := database.CreateChirpParams{
		Body:   cleaned,
		UserID: userID,
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create new chirp in db", err)
	}
	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
