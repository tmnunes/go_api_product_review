package api

import (
	"go_api_product_review/db"
	"go_api_product_review/models"
	"go_api_product_review/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Bearer token authentication. Example: 'Bearer token'

// RegisterProductRoutes initializes the routes for products
// @Summary Register product routes
// @Description Initializes the API endpoints for managing products
// @Tags products
// @Security ApiKeyAuth
func RegisterProductRoutes(router *gin.Engine) {
	productGroup := router.Group("/products")
	{
		productGroup.POST("/", CreateProduct)
		productGroup.GET("/", ListProducts)
		productGroup.GET("/:id", GetProductByID)
		productGroup.PUT("/:id", UpdateProduct)
		productGroup.DELETE("/:id", DeleteProduct)
	}
}

// CreateProduct creates a new product
// @Summary Create a new product
// @Description Creates a new product in the catalog
// @Tags products
// @Accept json
// @Produce json
// @Param product body models.Product true "Product details"
// @Success 201 {object} models.Product "Successfully created product"
// @Failure 400 {object} models.ErrorResponse "Invalid product data"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /products [post]
func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid product data",
			Details: err.Error(),
		})
		return
	}

	createdProduct, err := service.CreateProduct(db.GetDB(), &product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Failed to create product",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, createdProduct)
}

// ListProducts lists all products
// @Summary List all products
// @Description Fetches all products from the catalog
// @Tags products
// @Produce json
// @Success 200 {array} models.Product "List of products"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /products [get]
func ListProducts(c *gin.Context) {
	products, err := service.ListProducts(db.GetDB())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Failed to list products",
			Details: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, products)
}

// GetProductByID retrieves a product by its ID
// @Summary Get product by ID
// @Description Fetches a product by its unique ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} models.Product "Product found"
// @Failure 400 {object} models.ErrorResponse "Invalid product id"
// @Failure 500 {object} models.ErrorResponse "Failed to retrieve product"
// @Router /products/{id} [get]
func GetProductByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid product ID",
			Details: err.Error(),
		})
		return
	}

	productID := uint(id)
	product, err := service.GetProductByID(db.GetDB(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Failed to get product",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, product)
}

// UpdateProduct updates an existing product
// @Summary Update product
// @Description Updates an existing product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body models.Product true "Updated product details"
// @Success 200 {object} models.Product "Updated product"
// @Failure 400 {object} models.ErrorResponse "Invalid product data"
// @Failure 400 {object} models.ErrorResponse "Invalid product ID"
// @Failure 500 {object} models.ErrorResponse "Failed to update product"
// @Router /products/{id} [put]
func UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid product ID",
			Details: err.Error(),
		})
		return
	}

	productID := uint(id)
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid product data",
			Details: err.Error(),
		})
		return
	}

	updatedProduct, err := service.UpdateProduct(db.GetDB(), productID, &product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Failed to update product",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, updatedProduct)
}

// DeleteProduct deletes a product by its ID
// @Summary Delete product
// @Description Deletes a product by its ID
// @Tags products
// @Param id path int true "Product ID"
// @Success 204 "Product deleted"
// @Failure 400 {object} models.ErrorResponse "Invalid product ID"
// @Failure 500 {object} models.ErrorResponse "Failed to delete product"
// @Router /products/{id} [delete]
func DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid product ID",
			Details: err.Error(),
		})
		return
	}
	productID := uint(id)
	err = service.DeleteProduct(db.GetDB(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Failed to delete product",
			Details: err.Error(),
		})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// RegisterReviewRoutes initializes the routes for reviews
// @Summary Register review routes
// @Description Initializes the API endpoints for managing reviews
// @Tags reviews
// @Security ApiKeyAuth
func RegisterReviewRoutes(router *gin.Engine) {
	reviewGroup := router.Group("/reviews")
	{
		reviewGroup.POST("/", CreateReview)
		reviewGroup.PUT("/:id", UpdateReview)
		reviewGroup.DELETE("/:id", DeleteReview)
	}
}

// CreateReview creates a new review
// @Summary Create a new review
// @Description Creates a new review for a product
// @Tags reviews
// @Accept json
// @Produce json
// @Param review body models.Review true "Review details"
// @Success 201 {object} models.Review "Successfully created review"
// @Failure 400 {object} models.ErrorResponse "Invalid review data"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /reviews [post]
func CreateReview(c *gin.Context) {
	var review models.Review
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid review data",
			Details: err.Error(),
		})
		return
	}

	// Validate using the model's method
	if err := review.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid review data",
			Details: err.Error(),
		})
		return
	}

	createdReview, err := service.CreateReview(db.GetDB(), &review)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Failed creating review",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, createdReview)
}

// UpdateReview updates an existing review
// @Summary Update review
// @Description Updates an existing review by its ID
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "Review ID"
// @Param review body models.Review true "Updated review details"
// @Success 200 {object} models.Review "Updated review"
// @Failure 400 {object} models.ErrorResponse "Invalid review ID"
// @Failure 400 {object} models.ErrorResponse "Invalid review Data"
// @Failure 500 {object} models.ErrorResponse "Failed to update review"
// @Router /reviews/{id} [put]
func UpdateReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid review ID",
			Details: err.Error(),
		})
		return
	}

	var review models.Review
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid review data",
			Details: err.Error(),
		})
		return
	}

	updatedReview, err := service.UpdateReview(db.GetDB(), uint(id), &review)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Failed to update review",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, updatedReview)
}

// DeleteReview deletes a review by its ID
// @Summary Delete review
// @Description Deletes a review by its ID
// @Tags reviews
// @Param id path int true "Review ID"
// @Success 204 "Review deleted"
// @Failure 400 {object} models.ErrorResponse "Invalid review ID"
// @Failure 500 {object} models.ErrorResponse "Failed to delete review"
// @Router /reviews/{id} [delete]
func DeleteReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid review ID",
			Details: err.Error(),
		})
		return
	}

	err = service.DeleteReview(db.GetDB(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Failed to delete review",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
