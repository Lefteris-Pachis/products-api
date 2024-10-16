package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"products-api/database"
	"products-api/models"
	"strconv"
)

func CreateProduct(c *gin.Context) {
	var product models.Product
	// Bind JSON to product struct
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Create product in the database
	if err := database.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product", "details": err.Error()})
		return
	}

	// Return the created product with a 201 status
	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"product": product,
	})
}

func GetProductById(c *gin.Context) {
	// Get the product ID from URL parameters
	productIdStr := c.Param("id")

	// Validate that the ID is an unsigned integer
	productId, err := strconv.ParseUint(productIdStr, 10, 0) // set base:10 for decimal and bitSize:0 auto size
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID format"})
		return
	}

	var product models.Product
	// Attempt to find the product by ID
	if err := database.DB.First(&product, productId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			// Handle other possible database errors
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve product"})
			log.Println(err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, product)
}
