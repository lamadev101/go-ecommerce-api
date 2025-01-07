package utils

import (
	"errors"

	"github.com/lamadev101/ecommerce-api/types"
)

func RegisterUserValidation(user types.UserClient) error {
	if user.Name == "" {
		return errors.New("name can't be empty")
	}
	if user.Email == "" {
		return errors.New("email can't be empty")
	}
	if user.Password == "" {
		return errors.New("password can't be empty")
	}
	return nil
}

func ChangePasswordValidation(password types.ChangePassword) error {
	if password.OldPassword == "" {
		return errors.New("old password can't be empty")
	}
	if password.NewPassword == "" {
		return errors.New("new password can't be empty")
	}
	if password.ConfirmPassword == "" {
		return errors.New("confirm password can't be empty")
	}
	return nil
}
