package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HashPassword(password string) string {
	const bcryptCost = 14
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	return string(bytes)
}
