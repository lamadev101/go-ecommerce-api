package types

import (
	"regexp"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	Id          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Slug        string             `json:"slug" bson:"slug"`
	Description string             `json:"description" bson:"description"`
	Price       float64            `json:"price" bson:"price"`
	Category    string             `json:"category" bson:"category"`
	ImageURL    string             `json:"image_url" bson:"image_url"`
	Stock       int                `json:"stock" bson:"stock"`
	IsAvailable bool               `json:"is_available" bson:"is_available"`
	Tags        []string           `json:"tags" bson:"tags"`
	Variants    []ProductVariant   `json:"variants" bson:"variants"`
	CreatedAt   int64              `json:"created_at" bson:"created_at"`
	UpdatedAt   int64              `json:"updated_at" bson:"updated_at"`
}

type ProductClient struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Price       float64          `json:"price"`
	Category    string           `json:"category"`
	ImageURL    string           `json:"image_url"`
	Stock       int              `json:"stock"`
	IsAvailable bool             `json:"is_available"`
	Tags        []string         `json:"tags"`
	Variants    []ProductVariant `json:"variants"`
}

type ProductVariant struct {
	Id           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ProductID    primitive.ObjectID `json:"product_id" bson:"product_id"`
	VarientName  string             `json:"varient_name" bson:"varient_name"`
	Description  string             `json:"description" bson:"description"`
	VarientPrice float64            `json:"varient_price" bson:"varient_price"`
	ImageURL     string             `json:"image_url" bson:"image_url"`
	Stock        int                `json:"stock" bson:"stock"`
	IsAvailable  bool               `json:"is_available" bson:"is_available"`
	CreatedAt    int64              `json:"created_at" bson:"created_at"`
	UpdatedAt    int64              `json:"updated_at" bson:"updated_at"`
}

func (b *Product) BeforeCreate() (err error) {
	if b.Slug == "" { // Generate slug only if it is empty
		// b.Slug = utils.GenerateSlug(b.Name)
		slug := strings.ToLower(b.Name)

		// Remove special characters using a regex
		re := regexp.MustCompile(`[^a-z0-9\s-]`)
		slug = re.ReplaceAllString(slug, "")

		// Replace spaces and multiple hyphens with a single hyphen
		slug = strings.ReplaceAll(slug, " ", "-")
		re = regexp.MustCompile(`-{2,}`)
		slug = re.ReplaceAllString(slug, "-")

		// Trim hyphens from start and end
		b.Slug = strings.Trim(slug, "-")
	}
	return
}
