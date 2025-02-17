package models

import "github.com/jinzhu/gorm"

type Product struct {
	gorm.Model
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Price         float64  `json:"price"`
	AverageRating float64  `json:"average_rating"`
	Reviews       []Review `gorm:"foreignkey:ProductID"`
}
