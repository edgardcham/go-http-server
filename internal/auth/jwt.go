package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	// Step 1: Create the claims
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	}

	// step 2: Create token with claims and signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Step 3: sign the token with secret
	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	// ParseWithClaims needs:
	// 1. The token string
	// 2. Empty claims struct (to fill with data)
	// 3. A function that returns the secret key

	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{}, // Empty struct to fill
		// This function is called by the JWT library to retrieve the secret key
		// The library needs a function (not just the secret) so it can validate
		// the signing method and provide the correct key for token verification
		func(token *jwt.Token) (any, error) {
			return []byte(tokenSecret), nil
		},
	)

	if err != nil {
		return uuid.Nil, err
	}

	// Extract the user ID from claims
	// we're trying a type assertion first to ensure it's of type RegisteredClaims
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid claims")
	}

	// Convert Subject (string) back to UUID
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("no bearer set")
	}
	authHeaderArr := strings.Split(authHeader, " ")
	if len(authHeaderArr) != 2 {
		return "", fmt.Errorf("incorrect bearer format")
	}
	return authHeaderArr[1], nil
}

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("no api key set")
	}
	authHeaderArr := strings.Split(authHeader, " ")
	if len(authHeaderArr) != 2 {
		return "", fmt.Errorf("incorrect bearer format")
	}
	return authHeaderArr[1], nil
}
