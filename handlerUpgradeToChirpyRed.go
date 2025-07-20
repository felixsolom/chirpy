package main

import (
	"encoding/json"
	"net/http"

	"github.com/felixsolom/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUpgradeToChirpyRed(w http.ResponseWriter, r *http.Request) {
	type webhook struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Authentication required: malformed api key", err)
		return
	}

	if apiKey != cfg.polkaKey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var params webhook
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode JSON: %v", err)
		return
	}

	if params.Event != "user.upgraded" {
		respondWithError(w, http.StatusNoContent, "event doesn't match end point", nil)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid user id format", err)
		return
	}

	_, err = cfg.db.UpgradeToChirpyRed(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User wasn't found", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
