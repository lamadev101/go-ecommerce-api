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

func ProductRequestValidation(product types.ProductClient) error {
	if product.Name == "" {
		return errors.New("product name must be at least 3 characters long")
	}
	if product.Description == "" {
		return errors.New("description can't be empty")
	}
	if product.Price <= 0 {
		return errors.New("price must be greater than 0")
	}
	if product.ImageURL == "" {
		return errors.New("image URL can't be empty")
	}
	if product.Stock < 0 {
		return errors.New("stock can't be negative")
	}
	return nil
}

func BannerRequestValidation(banner types.BannerClient) error {
	if banner.Name == "" {
		return errors.New("banner name must be at least 3 characters long")
	}
	if banner.Description == "" {
		return errors.New("description can't be empty")
	}
	if banner.BtnLabel == "" {
		return errors.New("btn lable cant be empty")
	}
	if banner.BannerImageUrl == "" {
		return errors.New("banner image URL can't be empty")
	}
	if banner.ProductLink == "" {
		return errors.New("product link can't be empty")
	}
	if banner.BannerType == "" {
		return errors.New("banner type can't be empty")
	}
	return nil
}
