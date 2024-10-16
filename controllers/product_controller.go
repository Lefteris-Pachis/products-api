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

func GetProducts(c *gin.Context) {
	var products []models.Product // Slice to hold the products array

	// Retrieve all products from the database
	if err := database.DB.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve products"})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, products)
}

func DeleteProduct(c *gin.Context) {
	// Get the product ID from URL parameters
	productIdStr := c.Param("id")

	// Validate that the ID is an unsigned integer
	productId, err := strconv.ParseUint(productIdStr, 10, 0) // set base:10 for decimal and bitSize:0 auto size
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID format"})
		return
	}

	// Attempt to delete the product from the database
	result := database.DB.Delete(&models.Product{}, productId)

	// Check if the product was found and deleted
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete product"})
		log.Println(result.Error.Error())
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func UpdateProduct(c *gin.Context) {
	// Get the product ID from URL parameters
	productIdStr := c.Param("id")

	// Validate that the ID is an unsigned integer
	productId, err := strconv.ParseUint(productIdStr, 10, 0) // set base:10 for decimal and bitSize:0 auto size
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID format"})
		return
	}

	// Find the existing product in the database
	var product models.Product
	if err := database.DB.First(&product, productId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Create a temporary struct to hold the updated values
	var input struct {
		Name        *string  `json:"name"`
		Price       *float64 `json:"price"`
		Description *string  `json:"description"`
	}

	// Bind the incoming JSON to the input struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Track whether any changes were made
	var updated bool

	// Apply updates only if they are provided
	if input.Name != nil && *input.Name != product.Name {
		product.Name = *input.Name
		updated = true
	}
	if input.Price != nil && *input.Price != product.Price {
		product.Price = *input.Price
		updated = true
	}
	if input.Description != nil && *input.Description != product.Description {
		product.Description = *input.Description
		updated = true
	}

	// Only save if there were changes made to the product
	if updated {
		if err := database.DB.Save(&product).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update product"})
			log.Println(err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Product updated successfully",
			"product": product,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "No changes detected, product update not performed",
			"product": product,
		})
	}
}
