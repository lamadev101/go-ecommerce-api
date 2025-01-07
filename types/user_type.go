package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Email       string             `json:"email" bson:"email"`
	Password    string             `json:"password" bson:"password"`
	UserType    string             `json:"user_type" bson:"user_type"`
	Is2FAEnable bool               `json:"is_2fa_enable" bson:"is_2fa_enable"`
	CreatedAt   int64              `json:"created_at" bson:"created_at"`
	UpdatedAt   int64              `json:"updated_at" bson:"updated_at"`
}

type UserClient struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Is2FAEnable bool   `json:"is_2fa_enable"`
}

type DeliveryAddress struct {
	Id       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId   primitive.ObjectID `json:"user_id" bson:"user_id"`
	State    string             `json:"state" bson:"state"`
	District string             `json:"district" bson:"district"`
	City     string             `json:"city" bson:"city"`
	Address  string             `json:"address" bson:"address"`
}

type UserVerification struct {
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email     string             `json:"email" bson:"email"`
	Otp       int64              `json:"otp" bson:"otp"`
	Status    bool               `json:"status" bson:"status"`
	CreatedAt int64              `json:"created_at" bson:"created_at"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ChangePassword struct {
	Email           string `json:"email"`
	OldPassword     string `json:"old_password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}
