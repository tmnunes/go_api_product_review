package service_test_aux

import (
	"context"
	"go_api_product_review/cache"
	"go_api_product_review/models"
	"go_api_product_review/service"
	"strconv"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // Importar o SQLite para GORM
	"github.com/stretchr/testify/assert"
)

// MockRedisClient simulates Redis operations in memory
type MockRedisClient struct {
	data map[string]interface{}
}

// NewMockRedisClient initializes a new MockRedisClient
func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{
		data: make(map[string]interface{}),
	}
}

// Ping simulates a Redis ping to check connection status
func (m *MockRedisClient) Ping(ctx context.Context) *redis.StatusCmd {
	return redis.NewStatusResult("PONG", nil) // Simulate success response
}

// Set simulates setting a value in Redis with expiration time
func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	m.data[key] = value                     // Store value in the mock data map
	return redis.NewStatusResult("OK", nil) // Simulate success response
}

// Get simulates retrieving a value from Redis by key
func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	value, exists := m.data[key]
	if !exists {
		return redis.NewStringResult("", redis.Nil) // Simulate cache miss
	}
	// Return the value as a string
	return redis.NewStringResult(strconv.FormatFloat(value.(float64), 'f', -1, 64), nil)
}

// setupTestDB sets up an in-memory SQLite database for testing
func setupTestDB() (*gorm.DB, error) {
	// Open an in-memory SQLite database
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	// Auto-migrate models to create tables
	db.AutoMigrate(&models.Product{})
	db.AutoMigrate(&models.Review{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// TestGetProductByID tests the functionality of retrieving a product by ID and calculating the average rating
func TestGetProductByID(t *testing.T) {
	// Set up the test database and Redis client
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("failed to setup test db: %v", err)
	}
	defer db.Close()

	// Create mock Redis client and initialize cache
	mockClient := NewMockRedisClient()
	cache.InitRedis(mockClient)

	// Add a product and reviews to the database for testing
	product := models.Product{Name: "Test Product"}
	if err := db.Create(&product).Error; err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	// Add reviews for the product
	review1 := models.Review{ProductID: product.ID, Rating: 4}
	review2 := models.Review{ProductID: product.ID, Rating: 5}
	if err := db.Create(&review1).Error; err != nil {
		t.Fatalf("failed to create review1: %v", err)
	}
	if err := db.Create(&review2).Error; err != nil {
		t.Fatalf("failed to create review2: %v", err)
	}

	// Call the function under test (calculating average rating)
	avgRating, err := service.GetProductByID(db, product.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Calculate the expected average rating manually
	expectedAverageRating := float64(review1.Rating+review2.Rating) / float64(2)

	// Check if the calculated average rating is correct
	if avgRating != expectedAverageRating {
		t.Errorf("expected average rating %v, got %v", expectedAverageRating, avgRating)
	}

	// Verify if the average rating is correctly cached in Redis
	cachedRating, err := mockClient.Get(cache.Ctx, "product:"+strconv.Itoa(int(product.ID))+":average_rating").Result()
	if err != nil {
		t.Fatalf("expected no error getting cached rating, got %v", err)
	}
	cachedFloat, _ := strconv.ParseFloat(cachedRating, 64)
	if cachedFloat != expectedAverageRating {
		t.Errorf("expected cached average rating %v, got %v", expectedAverageRating, cachedFloat)
	}
}

// TestUpdateProductAverageRating tests the updating of a product's average rating
func TestUpdateProductAverageRating(t *testing.T) {
	// Set up the test database and Redis client
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("failed to setup test db: %v", err)
	}
	defer db.Close()

	// Create a mock Redis client and initialize cache
	mockClient := NewMockRedisClient()
	cache.InitRedis(mockClient)

	// Add a product to the database
	product := &models.Product{Name: "Test Product"}
	if err := db.Create(product).Error; err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	// Add reviews for the product
	review1 := &models.Review{ProductID: product.ID, FirstName: "BN", LastName: "SK", ReviewText: "Good product!", Rating: 4.0}
	review2 := &models.Review{ProductID: product.ID, FirstName: "JO", LastName: "LIV", ReviewText: "Excellent product!", Rating: 5.0}
	if err := db.Create(review1).Error; err != nil {
		t.Fatalf("failed to create review1: %v", err)
	}
	if err := db.Create(review2).Error; err != nil {
		t.Fatalf("failed to create review2: %v", err)
	}

	// Call the function under test (updating the product's average rating)
	err = service.UpdateProductAverageRating(db, product.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check if the average rating was updated correctly in the database
	var updatedProduct models.Product
	err = db.First(&updatedProduct, product.ID).Error
	if err != nil {
		t.Fatalf("failed to find updated product: %v", err)
	}

	// Calculate the expected average rating
	expectedAverageRating := (float64(review1.Rating) + float64(review2.Rating)) / float64(2)

	// Verify that the product's average rating was updated correctly
	if updatedProduct.AverageRating != expectedAverageRating {
		t.Errorf("expected average rating %v, got %v", expectedAverageRating, updatedProduct.AverageRating)
	}

	// Verify that the cache was updated with the new average rating
	cachedRating, err := mockClient.Get(cache.Ctx, "product:"+strconv.Itoa(int(product.ID))+":average_rating").Result()
	if err != nil {
		t.Fatalf("expected no error getting cached rating, got %v", err)
	}
	cachedFloat, _ := strconv.ParseFloat(cachedRating, 64)
	if cachedFloat != expectedAverageRating {
		t.Errorf("expected cached average rating %v, got %v", expectedAverageRating, cachedFloat)
	}
}

// TestListProducts tests the listing of all products and their associated reviews
func TestListProducts(t *testing.T) {
	// Set up the test database and Redis client
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("failed to setup test db: %v", err)
	}
	defer db.Close()

	// Create a mock Redis client and initialize cache
	mockClient := NewMockRedisClient()
	cache.InitRedis(mockClient)

	// Add a product and a review to the database for testing
	product := &models.Product{Name: "Test Product"}
	if err := db.Create(product).Error; err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	// Add a review for the product
	review := &models.Review{ProductID: product.ID, FirstName: "Rob", LastName: "Me", ReviewText: "Great product!", Rating: 5.0}
	if err := db.Create(review).Error; err != nil {
		t.Fatalf("failed to create review: %v", err)
	}

	// Call the function under test (listing products)
	products, err := service.ListProducts(db)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify that the product was listed and has the expected reviews
	assert.NotNil(t, products)
	assert.Len(t, products, 1) // One product should be returned
	assert.Equal(t, product.ID, products[0].ID)
	assert.Equal(t, "Test Product", products[0].Name)
	assert.NotNil(t, products[0].Reviews)
	assert.Len(t, products[0].Reviews, 1) // One review should be associated with the product
	assert.Equal(t, review.Rating, products[0].Reviews[0].Rating)
}

// TestUpdateReview tests the review update functionality.
func TestUpdateReview(t *testing.T) {
	// Setup test database
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("failed to setup test db: %v", err)
	}

	// Initialize mock Redis client
	mockClient := NewMockRedisClient()
	cache.InitRedis(mockClient)

	defer db.Close()

	// Create initial review for testing
	initialReview := &models.Review{
		ProductID:  2,
		FirstName:  "Bob",
		LastName:   "Boby",
		ReviewText: "Great",
		Rating:     5.0,
	}
	if err := db.Create(initialReview).Error; err != nil {
		t.Fatalf("failed to create initial review: %v", err)
	}

	// Prepare updated review data
	updatedReview := &models.Review{
		FirstName:  "Michael",
		LastName:   "Nunes",
		ReviewText: "Good product",
		Rating:     4.0,
	}

	// Call the UpdateReview function to update the review
	updated, err := service.UpdateReview(db, initialReview.ID, updatedReview)

	// Ensure no error occurred during the update
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check if the review was updated correctly
	assert.NotNil(t, updated)
	assert.Equal(t, updatedReview.FirstName, updated.FirstName)
	assert.Equal(t, updatedReview.LastName, updated.LastName)
	assert.Equal(t, updatedReview.ReviewText, updated.ReviewText)
	assert.Equal(t, updatedReview.Rating, updated.Rating)

	// Verify that the updated review is correctly saved in the database
	var savedReview models.Review
	err = db.First(&savedReview, initialReview.ID).Error
	if err != nil {
		t.Fatalf("expected review to be found, got error: %v", err)
	}
	assert.Equal(t, updatedReview.FirstName, savedReview.FirstName)
	assert.Equal(t, updatedReview.LastName, savedReview.LastName)
	assert.Equal(t, updatedReview.ReviewText, savedReview.ReviewText)
	assert.Equal(t, updatedReview.Rating, savedReview.Rating)
}

// TestCreateReview tests the functionality of creating a new review.
func TestCreateReview(t *testing.T) {
	// Setup test database
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("failed to setup test db: %v", err)
	}

	// Initialize mock Redis client
	mockClient := NewMockRedisClient()
	cache.InitRedis(mockClient)

	defer db.Close()

	// Initialize a new review object
	review := models.Review{ProductID: 2, Rating: 5.0}

	// Call CreateReview function to create the review
	createdReview, err := service.CreateReview(db, &review)

	// Ensure no error occurred during creation
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify the review was created successfully
	assert.NotNil(t, createdReview)
	assert.Equal(t, review.ProductID, createdReview.ProductID)
	assert.Equal(t, review.Rating, createdReview.Rating)

	// Check if the review was saved correctly in the database
	var savedReview models.Review
	err = db.First(&savedReview, createdReview.ID).Error
	if err != nil {
		t.Fatalf("expected review to be found, got error: %v", err)
	}
	assert.Equal(t, review.Rating, savedReview.Rating)
}

// TestDeleteReview tests the functionality of deleting a review.
func TestDeleteReview(t *testing.T) {
	// Setup test database
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("failed to setup test db: %v", err)
	}

	// Initialize mock Redis client
	mockClient := NewMockRedisClient()
	cache.InitRedis(mockClient)

	defer db.Close()

	// Create a new review to test deletion
	review := models.Review{ProductID: 2}
	if err := db.Create(&review).Error; err != nil {
		t.Fatalf("failed to create review: %v", err)
	}

	// Call DeleteReview function to delete the review
	err = service.DeleteReview(db, review.ID)

	// Ensure no error occurred during deletion
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check that the review has been deleted
	var deletedReview models.Review
	err = db.First(&deletedReview, review.ID).Error
	if err == nil {
		t.Fatalf("expected review to be deleted, but it was found")
	}
}

// Test storing product average rating in cache
func TestCacheProductAverageRating(t *testing.T) {
	// Create a mocked Redis client
	mockClient := NewMockRedisClient()
	cache.InitRedis(mockClient) // Initialize cache with the mock client

	productID := uint(1) // Product ID to be stored in the cache
	rating := 4.5        // Average rating of the product

	// Call the function to store the average rating in the cache
	err := service.CacheProductAverageRating(productID, rating)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify if the value was correctly stored in the mocked Redis
	key := "product:" + strconv.Itoa(int(productID)) + ":average_rating"
	if value, exists := mockClient.data[key]; !exists || value != rating {
		t.Fatalf("expected value %f for key %s, got %v", rating, key, value)
	}
}

// Test retrieving cached product average rating
func TestGetCachedProductAverageRating(t *testing.T) {
	// Create a new instance of the mocked Redis client
	mockClient := NewMockRedisClient()
	cache.InitRedis(mockClient) // Initialize cache with the mock client

	productID := uint(1) // Product ID
	rating := 4.5        // Average rating of the product

	// Simulate a cache hit by storing the rating in the mock
	key := "product:" + strconv.Itoa(int(productID)) + ":average_rating"
	mockClient.data[key] = rating

	// Test retrieving the average rating when the key exists
	cachedRating, err := service.GetCachedProductAverageRating(productID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cachedRating != rating {
		t.Fatalf("expected cached rating %f, got %f", rating, cachedRating)
	}

	// Test behavior when the key does not exist (cache miss)
	delete(mockClient.data, key) // Remove the value to simulate a cache miss
	cachedRating, err = service.GetCachedProductAverageRating(productID)
	if err != nil {
		t.Fatalf("expected no error on cache miss, got %v", err)
	}
	if cachedRating != 0 {
		t.Fatalf("expected cached rating 0 on cache miss, got %f", cachedRating)
	}
}

// Test updating a product in the database
func TestUpdateProduct(t *testing.T) {
	// Set up an in-memory SQLite database
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Error connecting to in-memory database: %v", err)
	}
	defer db.Close()

	// Create the necessary table in the database
	db.AutoMigrate(&models.Product{})

	// Create a sample product and save it to the database
	product := &models.Product{
		Name:        "Original Product",
		Description: "Original Description",
		Price:       100.0,
	}
	db.Create(product)

	// Define new data for product update
	updatedProduct := &models.Product{
		Name:        "Updated Product",
		Description: "Updated Description",
		Price:       150.0,
	}

	// Call the update product function
	updatedProductResult, err := service.UpdateProduct(db, product.ID, updatedProduct)

	// Verify that there were no errors and data was updated correctly
	assert.NoError(t, err)
	assert.Equal(t, "Updated Product", updatedProductResult.Name)
	assert.Equal(t, "Updated Description", updatedProductResult.Description)
	assert.Equal(t, 150.0, updatedProductResult.Price)
}

// Test deleting a product from the database
func TestDeleteProduct(t *testing.T) {
	// Set up an in-memory SQLite database
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Error connecting to in-memory database: %v", err)
	}
	defer db.Close()

	// Create the necessary table in the database
	db.AutoMigrate(&models.Product{})

	// Create and save a product in the database to be deleted
	product := &models.Product{
		Name:        "Product to be deleted",
		Description: "This product will be deleted",
		Price:       200.0,
	}
	db.Create(product)

	// Call the delete product function
	err = service.DeleteProduct(db, product.ID)

	// Verify that there were no errors during deletion
	assert.NoError(t, err)

	// Check if the product was actually deleted from the database
	var deletedProduct models.Product
	result := db.First(&deletedProduct, product.ID)
	assert.Error(t, result.Error) // Should result in an error as the product was deleted
}
