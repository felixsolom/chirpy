package main

import (
	"net/http"

	"github.com/felixsolom/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {

	refreshTokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "No valid refresh token provided: %w", err)
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), refreshTokenStr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't revoke refresh token %w", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}
