package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashPassword), nil
}

//check if provided password is corrrect
func ValidatePassword(password string, hashedPassord string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassord), []byte(password))
}