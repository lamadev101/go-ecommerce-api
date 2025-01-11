package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Banner struct {
	Id             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name           string             `json:"name" bson:"name"`
	Description    string             `json:"description" bson:"description"`
	BtnLabel       string             `json:"btn_label" bson:"btn_label"`
	ProductLink    string             `json:"product_link" bson:"product_link"`
	BannerImageUrl string             `json:"banner_image_url" bson:"banner_image_url"`
	BannerType     string             `json:"banner_type" bson:"banner_type"`
	CreatedAt      int64              `json:"created_at" bson:"created_at"`
	UpdatedAt      int64              `json:"updated_at" bson:"updated_at"`
}

type BannerClient struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	BtnLabel       string `json:"btn_label"`
	ProductLink    string `json:"product_link"`
	BannerImageUrl string `json:"banner_image_url"`
	BannerType     string `json:"banner_type"`
}
