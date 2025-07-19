package main

import (
	"net/http"
	"time"

	"github.com/felixsolom/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type NewlyMintedAccessToken struct {
		Token string `json:"token"`
	}

	refreshTokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "No valid refresh token provided: %w", err)
		return
	}
	refreshToken, err := cfg.db.GetRefreshToken(r.Context(), refreshTokenStr)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "token not in db: %w", err)
		return
	}
	if refreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "refresh token has been revoked", nil)
		return
	}
	if refreshToken.ExpiresAt.Before(time.Now().UTC()) {
		respondWithError(w, http.StatusUnauthorized, "refresh token has expired", nil)
		return
	}

	newAccessToken, err := auth.MakeJWT(refreshToken.UserID, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "coudln't create new access token: %w", err)
		return
	}
	respondWithJSON(w, http.StatusOK, NewlyMintedAccessToken{
		Token: newAccessToken,
	})
}
