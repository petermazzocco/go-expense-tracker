package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Create a new JWT token
func GenerateJWT(userID string, key string) (string, error) {
	// Get secret key
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return "", errors.New("SECRET_KEY is not set in the environment")
	}

	// Create token with claims
	claims := jwt.MapClaims{
		"sub": userID,                                // Subject (User ID instead of username)
		"key": key,                                   // Encryption key
		"exp": time.Now().Add(24 * time.Hour).Unix(), // Expiration time
		"iat": time.Now().Unix(),                     // Issued At
		"nbf": time.Now().Unix(),                     // Not Before
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT checks the token's validity
func ValidateJWT(tokenString string) (map[string]any, error) {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return nil, errors.New("SECRET_KEY is not set in environment")
	}

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// Extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
