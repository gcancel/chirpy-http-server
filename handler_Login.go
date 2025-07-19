package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gcancel/chirpy/internal/auth"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type parameters struct {
		Email            string
		Password         string
		ExpiresInSeconds int `json:"expires_in_seconds"`
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	var loginInfo parameters
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&loginInfo)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding json", err)
		return
	}
	user, err := cfg.dbQueries.GetUser(req.Context(), loginInfo.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error retrieving user from database", err)
		return
	}
	err = auth.CheckPasswordHash(loginInfo.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid password", err)
		return
	}

	// setting default expiration time for token
	expiresIn := time.Hour
	if loginInfo.ExpiresInSeconds > 0 && loginInfo.ExpiresInSeconds < 3600 {
		expiresIn = time.Duration(loginInfo.ExpiresInSeconds) * time.Second
	}

	token, err := auth.MakeJWT(user.ID, cfg.secretToken, time.Duration(expiresIn))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating JWT token", err)
		return
	}

	fmt.Println(loginInfo)
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			Id:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token: token,
	})
	fmt.Printf("Login Successful. User: %v\n", user.Email)

}
