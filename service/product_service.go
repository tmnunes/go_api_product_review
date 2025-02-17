package service

import (
	"go_api_product_review/db"
	"go_api_product_review/models"
)

func CreateProduct(product *models.Product) (*models.Product, error) {
	result := db.GetDB().Create(product)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}

func UpdateProduct(id uint, updatedProduct *models.Product) (*models.Product, error) {
	var product models.Product
	result := db.GetDB().First(&product, id)
	if result.Error != nil {
		return nil, result.Error
	}

	product.Name = updatedProduct.Name
	product.Description = updatedProduct.Description
	product.Price = updatedProduct.Price
	db.GetDB().Save(&product)
	return &product, nil
}

func DeleteProduct(id uint) error {
	result := db.GetDB().Delete(&models.Product{}, id)
	return result.Error
}

func GetProductByID(id uint) (*models.Product, error) {
	var product models.Product
	result := db.GetDB().Preload("Reviews").First(&product, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}

func ListProducts() ([]models.Product, error) {
	var products []models.Product
	result := db.GetDB().Preload("Reviews").Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

func UpdateProductAverageRating(productID uint) error {
	var product models.Product
	var reviews []models.Review
	db.GetDB().First(&product, productID)
	db.GetDB().Where("product_id = ?", productID).Find(&reviews)

	var totalRating int
	for _, review := range reviews {
		totalRating += review.Rating
	}

	averageRating := float64(totalRating) / float64(len(reviews))

	// Update the product with the new average rating
	product.AverageRating = averageRating
	db.GetDB().Save(&product)

	return nil
}

// CreateReview creates a new review for a product
func CreateReview(review *models.Review) (*models.Review, error) {
	result := db.GetDB().Create(review)
	if result.Error != nil {
		return nil, result.Error
	}

	// Update the average rating of the product after review creation
	err := UpdateProductAverageRating(review.ProductID)
	if err != nil {
		return nil, err
	}

	return review, nil
}

// UpdateReview updates an existing review
func UpdateReview(id uint, updatedReview *models.Review) (*models.Review, error) {
	var review models.Review
	result := db.GetDB().First(&review, id)
	if result.Error != nil {
		return nil, result.Error
	}

	review.FirstName = updatedReview.FirstName
	review.LastName = updatedReview.LastName
	review.ReviewText = updatedReview.ReviewText
	review.Rating = updatedReview.Rating

	db.GetDB().Save(&review)
	// Recalculate the product's average rating after review update
	err := UpdateProductAverageRating(review.ProductID)
	if err != nil {
		return nil, err
	}

	return &review, nil
}

// DeleteReview deletes a review and updates the average rating of the product
func DeleteReview(id uint) error {
	var review models.Review
	result := db.GetDB().First(&review, id)
	if result.Error != nil {
		return result.Error
	}

	// Delete the review
	db.GetDB().Delete(&review)

	// Recalculate the product's average rating after review deletion
	err := UpdateProductAverageRating(review.ProductID)
	if err != nil {
		return err
	}

	return nil
}
