package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword mengenkripsi password dengan bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	
	return string(hashedPassword), nil
}

// CheckPassword membandingkan password dengan hash
func CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}