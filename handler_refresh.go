package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gcancel/chirpy/internal/auth"
)

func (cfg *apiConfig) handleRefresh(w http.ResponseWriter, req *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(req.Header)
	fmt.Println(token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "bearer token not found", err)
		return
	}
	storedToken, err := cfg.dbQueries.GetRefreshToken(req.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "refresh token not found", err)
		return
	}

	// confusing, but will check if the value stored in RevokedAt is not a null
	if storedToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "token revoked", err)
		return
	}

	if time.Now().After(storedToken.ExpiresAt) {
		respondWithError(w, http.StatusUnauthorized, "expired token", err)
		return
	}

	expiresIn := time.Hour
	newToken, err := auth.MakeJWT(storedToken.UserID, cfg.secretToken, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error refreshing token", err)
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		Token: newToken,
	})

}

func (cfg *apiConfig) handleRevoke(w http.ResponseWriter, req *http.Request) {
	type response struct {
		//
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "bearer token not found", err)
		return
	}
	revoked, err := cfg.dbQueries.RevokeRefreshToken(req.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error revoking refresh token", err)
		return
	}
	fmt.Printf("Revoked Token: %v", revoked.Token)

	respondWithJSON(w, http.StatusNoContent, response{})
}
