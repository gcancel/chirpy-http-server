package main

import (
	"net/http"

	"github.com/gcancel/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleChirpDelete(w http.ResponseWriter, req *http.Request) {

	type response struct {
		//
	}

	chirpParam := req.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpParam)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error parsing UUID parameter", err)
		return
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "bearer token not found", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secretToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error validating JWT token", err)
		return
	}

	chirp, err := cfg.dbQueries.GetChirp(req.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error retrieving chirp", err)
		return
	}
	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "user not allowed to delete chirps of other users", err)
		return
	}

	err = cfg.dbQueries.DeleteChirp(req.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, response{})

}
