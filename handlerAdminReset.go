package main

import "net/http"

func (cfg *apiConfig) handlerAdminReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Permission denied", nil)
	}

	cfg.fileserverHits.Store(0)

	err := cfg.db.DeleteAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete all chirps", err)
		return
	}

	err = cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete users", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All users and chirps were deleted successfully"))
}
