package models

import "github.com/jinzhu/gorm"

type Review struct {
	gorm.Model
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	ReviewText string `json:"review_text"`
	Rating     int    `json:"rating"`
	ProductID  uint   `json:"product_id"`
}
