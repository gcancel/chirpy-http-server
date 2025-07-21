package main

import (
	"encoding/json"
	"net/http"

	"github.com/gcancel/chirpy/internal/auth"
	"github.com/gcancel/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleUpdateUser(w http.ResponseWriter, req *http.Request) {

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		Email  string    `json:"email"`
		UserID uuid.UUID `json:"user_id"`
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "bearer token not found", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secretToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error validating JWT", err)
		return
	}

	var loginInfo parameters
	decoder := json.NewDecoder(req.Body)
	decoder.Decode(&loginInfo)

	hashed_password, err := auth.HashPassword(loginInfo.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error hashing password", err)
		return
	}

	updatedUser, err := cfg.dbQueries.UpdateUser(req.Context(), database.UpdateUserParams{
		Email:          loginInfo.Email,
		HashedPassword: hashed_password,
		ID:             userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error updating user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Email:  updatedUser.Email,
		UserID: updatedUser.ID,
	})

}
