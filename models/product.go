package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// Product represents a product in the database
// @Description Represents a product in the store or catalog
// @Schema
type Product struct {
	gorm.Model
	// Name of the product
	// @example "Bananas"
	Name string `json:"name" validate:"required"`
	// Description of the product
	// @example "Bananas from Argentina"
	Description string `json:"description"`
	// Price of the product
	// @example 20.00
	Price float64 `json:"price" validate:"required"`
	// Average rating of the product based on reviews
	// @example 4.5
	AverageRating float64 `json:"average_rating"`
	// Reviews associated with this product
	// @example []Review
	// @readOnly
	Reviews []Review `gorm:"foreignkey:ProductID"`
}

// Validate checks if the product fields are valid.
func (p *Product) Validate() error {
	if p.Price <= 0 {
		return errors.New("price must be greater than 0")
	}
	if p.Name == "" {
		return errors.New("name is required")
	}
	return nil
}
