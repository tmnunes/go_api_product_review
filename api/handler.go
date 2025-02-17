package api

import (
	"net/http"
	"go_api_product_review/models"
	"go_api_product_review/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

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

func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	createdProduct, err := service.CreateProduct(&product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdProduct)
}

func ListProducts(c *gin.Context) {
	products, err := service.ListProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

func GetProductByID(c *gin.Context) {
	// Get the product ID from URL parameters and convert it to uint
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr) // Convert the string to an integer
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Convert the integer to uint
	productID := uint(id)

	// Call the service to get the product by ID
	product, err := service.GetProductByID(productID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Respond with the product
	c.JSON(http.StatusOK, product)
}

func UpdateProduct(c *gin.Context) {
	// Get the product ID from URL parameters and convert to uint
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr) // Convert the string to an integer
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Convert the integer to uint
	productID := uint(id)

	// Bind the request body to a product model
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service layer to update the product
	updatedProduct, err := service.UpdateProduct(productID, &product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with the updated product
	c.JSON(http.StatusOK, updatedProduct)
}

func DeleteProduct(c *gin.Context) {
	// Get the product ID from URL parameters and convert it to uint
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr) // Convert the string to an integer
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Convert the integer to uint
	productID := uint(id)

	// Call the service to delete the product by ID
	err = service.DeleteProduct(productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with a success message
	c.JSON(http.StatusNoContent, nil)
}

func RegisterReviewRoutes(router *gin.Engine) {
	reviewGroup := router.Group("/reviews")
	{
		reviewGroup.POST("/", CreateReview)
		reviewGroup.PUT("/:id", UpdateReview)
		reviewGroup.DELETE("/:id", DeleteReview)
	}
}

func CreateReview(c *gin.Context) {
	var review models.Review
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdReview, err := service.CreateReview(&review)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdReview)
}

func UpdateReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	var review models.Review
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedReview, err := service.UpdateReview(uint(id), &review)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedReview)
}

func DeleteReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	err = service.DeleteReview(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
