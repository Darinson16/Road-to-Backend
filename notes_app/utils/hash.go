package utils

import "golang.org/x/crypto/bcrypt"

const bcryptCost = 12 // or 10 if you want very fast login

func HashPassword(password string) (string, error) {
	newPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	return string(newPassword), err
}

func CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
