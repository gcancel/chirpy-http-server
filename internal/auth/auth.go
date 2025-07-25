package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return fmt.Errorf("invalid password")
	}
	fmt.Println("Password valid!")
	return nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	tokenIssuer := "chirpy-access"
	claims := jwt.RegisteredClaims{
		Issuer:    tokenIssuer,
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claimsStruct, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("error parsing JWT token %v", err)
	}

	id, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, fmt.Errorf("error retrieving subject from claim %v", err)
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != string("chirpy-access") {
		return uuid.Nil, fmt.Errorf("invalid issuer")
	}

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID %v", err)
	}
	return parsedUUID, nil

}

func GetBearerToken(headers http.Header) (string, error) {
	bearerHeader := headers.Values("Authorization")
	if len(bearerHeader) < 1 {
		return "", fmt.Errorf("no bearer token in header")
	}
	//fmt.Println(bearerHeader)
	bearerToken := strings.TrimPrefix(bearerHeader[0], "Bearer ")

	return bearerToken, nil
}

func MakeRefreshToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	refreshToken := hex.EncodeToString(token)
	return refreshToken, nil
}

func GetAPIKey(headers http.Header) (string, error) {
	apiKeyHeader := headers.Values("Authorization")
	if len(apiKeyHeader) < 1 {
		return "", fmt.Errorf("no apiKey in header")
	}
	apiKey := strings.TrimPrefix(apiKeyHeader[0], "ApiKey ")

	return apiKey, nil
}
