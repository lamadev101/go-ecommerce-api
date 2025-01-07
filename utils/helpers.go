package utils

import "golang.org/x/crypto/bcrypt"

func GenerateHashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return ""
	}
	return string(bytes)
}
