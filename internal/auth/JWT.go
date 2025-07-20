package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error) {

	now := time.Now().UTC()
	hour := time.Hour

	claims := &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(hour)),
		Issuer:    "chirpy",
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("token couldn't be produced: %w", err)
	}
	return ss, err
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	type UserClaims struct {
		jwt.RegisteredClaims
	}

	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("couldn't parse JWT: %w", err)
	} else if claims, ok := token.Claims.(*UserClaims); ok {
		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			return uuid.Nil, fmt.Errorf("subject is not a valid uuid: %w", err)
		}
		return userID, nil
	}
	return uuid.Nil, fmt.Errorf("unknown claims type, can't proceed")
}

func GetBearerToken(headers http.Header) (string, error) {
	bearerToken := headers.Get("Authorization")
	if bearerToken == "" {
		return "", fmt.Errorf("no authentication header was found")
	}
	bearerToken = strings.TrimSpace(strings.TrimPrefix(bearerToken, "Bearer"))
	return bearerToken, nil
}

func GetAPIKey(headers http.Header) (string, error) {
	apiKey := headers.Get("Authorization")
	if apiKey == "" {
		return "", fmt.Errorf("no authentication header was found")
	}
	apiKey = strings.TrimSpace(strings.TrimPrefix(apiKey, "ApiKey"))
	return apiKey, nil
}
