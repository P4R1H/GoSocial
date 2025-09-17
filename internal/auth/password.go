package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword takes a plaintext password and returns a bcrypt hash
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedBytes), nil
}

// ComparePassword checks if a plaintext password matches a hashed password
func ComparePassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		// bcrypt.ErrMismatchedHashAndPassword means wrong password
		return fmt.Errorf("invalid password: %w", err)
	}

	return nil
}
