package models

import "github.com/jinzhu/gorm"

// Review represents a review for a product in the database
// @Description Represents a review for a specific product, including the reviewer's name, review text, and rating.
// @Schema
type Review struct {
	gorm.Model
	// First name of the reviewer
	// @example "Miguel"
	FirstName string `json:"first_name"`
	// Last name of the reviewer
	// @example "Filip"
	LastName string `json:"last_name"`
	// Text content of the review
	// @example "This bananas are amazing!"
	ReviewText string `json:"review_text"`
	// Rating given by the reviewer (1-5)
	// @example 4
	Rating int `json:"rating"`
	// ProductID is the foreign key that links to the product being reviewed
	// @example 999
	ProductID uint `json:"product_id"`
}
