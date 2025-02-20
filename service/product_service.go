package service

import (
	"encoding/json"
	"errors"
	"go_api_product_review/cache"
	"go_api_product_review/models"
	"math"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Import the PostgreSQL dialect for GORM
)

// CreateProduct creates a new product in the database.
// Pass a mock or real DB instance as a parameter for testing.
func CreateProduct(db *gorm.DB, product *models.Product) (*models.Product, error) {
	result := db.Create(product)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}

// UpdateProduct updates an existing product in the database by ID.
// It accepts an ID and an updated product object.
// Returns the updated product or an error if the operation fails.
func UpdateProduct(db *gorm.DB, id uint, updatedProduct *models.Product) (*models.Product, error) {
	var product models.Product

	result := db.First(&product, id)

	if result.Error != nil {
		return nil, result.Error
	}

	// Update product fields
	product.Name = updatedProduct.Name
	product.Description = updatedProduct.Description
	product.Price = updatedProduct.Price
	db.Save(&product)
	return &product, nil
}

// DeleteProduct deletes a product by its ID.
// It accepts a product ID and removes the product from the database.
// Returns an error if the deletion fails.
func DeleteProduct(db *gorm.DB, id uint) error {
	result := db.Delete(&models.Product{}, id)
	key := "product:" + strconv.Itoa(int(id)) + ":average_rating"
	err := cache.Rdb.Del(cache.Ctx, key).Err()
	if err != nil {
		return err
	}

	return result.Error
}

// GetProductByID retrieves a product's average rating by its ID.
// First, it checks the cache for the average rating.
// If not found, it calculates the rating, updates the cache, and returns the value.
func GetProductByID(db *gorm.DB, id uint) (float64, error) {
	// Attempt to get the cached average rating
	avgRating, err := GetCachedProductAverageRating(id)
	if err != nil {
		return 0, err
	}

	// If no cached value, calculate the average rating from reviews
	if avgRating == 0 {
		var product models.Product
		result := db.First(&product, id)
		if result.Error != nil {
			return 0, result.Error
		}

		// Recalculate and update the cached average rating
		err := UpdateProductAverageRating(db, id)
		if err != nil {
			return 0, err
		}

		// Fetch the newly cached average rating
		avgRating, err = GetCachedProductAverageRating(id)
		if err != nil {
			return avgRating, err
		}
	}

	// Handle case where the rating is not available
	if math.IsNaN(avgRating) {
		return 0, errors.New("Error: Product average rating is not available")
	}

	return avgRating, nil
}

// ListProducts retrieves all products along with their reviews.
// It returns a list of products or an error if the operation fails.
func ListProducts(db *gorm.DB) ([]models.Product, error) {
	var products []models.Product
	result := db.Preload("Reviews").Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	// Ensure AverageRating is not NaN
	for i := range products {
		if math.IsNaN(products[i].AverageRating) {
			products[i].AverageRating = 0 // Set to 0 if NaN
		}
	}
	return products, nil
}

// UpdateProductAverageRating recalculates the average rating of a product
// based on its reviews and updates both the product record and the cache.
func UpdateProductAverageRating(db *gorm.DB, productID uint) error {
	var product models.Product
	var reviews []models.Review

	db.First(&product, productID)
	db.Where("product_id = ?", productID).Find(&reviews)

	// Calculate the total rating from reviews
	var totalRating int
	for _, review := range reviews {
		totalRating += review.Rating
	}

	averageRating := 0.0
	// Calculate the average rating
	if totalRating != 0 && len(reviews) != 0 {
		averageRating = float64(totalRating) / float64(len(reviews))
	}

	// Cache the updated average rating
	err := CacheProductAverageRating(productID, averageRating)
	if err != nil {
		return err
	}

	// Update the product record with the new average rating
	product.AverageRating = averageRating
	db.Save(&product)

	return nil
}

// CreateReview creates a new review for a product and updates the product's average rating.
// It accepts a review object, saves it to the database, and recalculates the product's average rating.
func CreateReview(db *gorm.DB, review *models.Review) (*models.Review, error) {
	result := db.Create(review)
	if result.Error != nil {
		return nil, result.Error
	}

	// Cache the new review
	err := CacheReview(review)
	if err != nil {
		return nil, err
	}

	// Update the product's average rating after the review is added
	err = UpdateProductAverageRating(db, review.ProductID)
	if err != nil {
		return nil, err
	}

	return review, nil
}

// UpdateReview updates an existing review for a product by its ID.
// It accepts a review ID and the updated review data, saves the updated review,
// and recalculates the product's average rating.
func UpdateReview(db *gorm.DB, id uint, updatedReview *models.Review) (*models.Review, error) {
	var review models.Review
	result := db.First(&review, id)
	if result.Error != nil {
		return nil, result.Error
	}

	// Update review fields
	review.FirstName = updatedReview.FirstName
	review.LastName = updatedReview.LastName
	review.ReviewText = updatedReview.ReviewText
	review.Rating = updatedReview.Rating

	db.Save(&review)

	err := CacheReview(&review)
	if err != nil {
		return nil, err
	}

	// Recalculate the product's average rating after the review is updated
	err = UpdateProductAverageRating(db, review.ProductID)
	if err != nil {
		return nil, err
	}

	return &review, nil
}

// DeleteReview deletes an existing review by its ID and recalculates the product's average rating.
// It accepts a review ID, deletes the review, and updates the product's rating accordingly.
func DeleteReview(db *gorm.DB, id uint) error {
	var review models.Review
	result := db.First(&review, id)
	if result.Error != nil {
		return result.Error
	}

	// Delete the review from the database
	db.Delete(&review)

	// Remove the review from the Redis cache
	err := cache.Rdb.Del(cache.Ctx, "review:"+strconv.Itoa(int(id))).Err()
	if err != nil {
		return err
	}

	err = UpdateProductAverageRating(db, review.ProductID)
	if err != nil {
		return err
	}

	return nil
}

// CacheProductAverageRating caches the average rating of a product in Redis.
// It stores the average rating under a unique key based on the product ID.
func CacheProductAverageRating(productID uint, rating float64) error {
	key := "product:" + strconv.Itoa(int(productID)) + ":average_rating"
	err := cache.Rdb.Set(cache.Ctx, key, rating, 10*time.Minute).Err()
	return err
}

// GetCachedProductAverageRating retrieves the cached average rating of a product from Redis.
// It returns the cached rating if available or zero if it's a cache miss.
func GetCachedProductAverageRating(productID uint) (float64, error) {
	key := "product:" + strconv.Itoa(int(productID)) + ":average_rating"
	cachedRating, err := cache.Rdb.Get(cache.Ctx, key).Result()
	if err == redis.Nil {
		return 0, nil // Cache miss, rating not found
	}
	if err != nil {
		return 0, err
	}

	// Convert cached rating to float64
	return strconv.ParseFloat(cachedRating, 64)
}

// GetReview retrieves a review from the cache or database if not found in cache.
func GetReview(db *gorm.DB, id uint) (*models.Review, error) {
	// Check the Redis cache first
	reviewKey := "review:" + strconv.Itoa(int(id))
	reviewJSON, err := cache.Rdb.Get(cache.Ctx, reviewKey).Result()

	if err == redis.Nil {
		// Review not found in cache, query the database
		var review models.Review
		result := db.First(&review, id)
		if result.Error != nil {
			return nil, result.Error // Return error if not found
		}

		// Cache the new review
		reviewJSON, err := json.Marshal(review)
		if err != nil {
			return nil, err
		}

		// Store the JSON string in Redis
		err = cache.Rdb.Set(cache.Ctx, reviewKey, reviewJSON, 10*time.Minute).Err()
		if err != nil {
			return nil, err // Return error if Redis store fails
		}
	} else if err != nil {
		return nil, err // Return error if there is a Redis issue
	}

	// Unmarshal JSON back to Review struct
	var review models.Review
	err = json.Unmarshal([]byte(reviewJSON), &review) // Unmarshal to struct
	if err != nil {
		return nil, err // Return error if unmarshaling fails
	}

	return &review, nil
}

// CacheReview updates the Redis cache with the given review.
func CacheReview(review *models.Review) error {
	// Marshal the review to a JSON byte slice
	reviewJSON, err := json.Marshal(review)
	if err != nil {
		return err
	}

	// Update the cache with the review
	err = cache.Rdb.Set(cache.Ctx, "review:"+strconv.Itoa(int(review.ID)), reviewJSON, 10*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}
