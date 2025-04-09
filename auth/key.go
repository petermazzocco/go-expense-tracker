package auth

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/argon2"
)

// GenerateSalt creates a new random salt (16 bytes)
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

// DeriveKey generates a cryptographic key from a password + salt
func DeriveKey(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
}

// EncodeToBase64 converts bytes to Base64
func EncodeToBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func DecodeFromBase64(key string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(key)
}
