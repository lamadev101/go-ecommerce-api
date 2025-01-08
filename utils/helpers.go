package utils

import (
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func GenerateHashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func ConverStringIntoInt(str string) (int, error) {
	integer, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return integer, err
}
