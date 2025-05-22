package auth

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	pwd := []byte(password)
	hashedPass, errHash := bcrypt.GenerateFromPassword(pwd, 3)
	if errHash != nil {
		log.Printf("error hashing password: %v", errHash)
		return "", errHash
	}
	return string(hashedPass), nil
}

func CheckPasswordHash(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
