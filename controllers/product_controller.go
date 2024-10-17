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

// Utility function to parse a product ID from the URL parameters
func parseProductID(c *gin.Context) (uint64, error) {
	productIdStr := c.Param("id")
	// Validate that the ID is an unsigned integer
	productId, err := strconv.ParseUint(productIdStr, 10, 0) // set base:10 for decimal and bitSize:0 auto size
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID format"})
	}
	return productId, err
}

// Utility function to bind JSON to a struct and handle errors
func bindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return false
	}
	return true
}

// Utility function to respond with an error when a database operation fails
func handleDBError(c *gin.Context, err error, errorMessage string) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": errorMessage})
	log.Println(err.Error())
}

func CreateProduct(c *gin.Context) {
	var product models.Product
	if !bindJSON(c, &product) {
		return
	}

	// Create product in the database
	if err := database.DB.Create(&product).Error; err != nil {
		handleDBError(c, err, "Failed to create product")
		return
	}

	// Return the created product with a 201 status
	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"product": product,
	})
}

func GetProductById(c *gin.Context) {
	productId, err := parseProductID(c)
	if err != nil {
		return
	}

	var product models.Product
	// Attempt to find the product by ID
	if err := database.DB.First(&product, productId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			// Handle other possible database errors
			handleDBError(c, err, "Could not retrieve product")
		}
		return
	}

	c.JSON(http.StatusOK, product)
}

func GetProducts(c *gin.Context) {
	var products []models.Product // Slice to hold the products array

	// Retrieve all products from the database
	if err := database.DB.Find(&products).Error; err != nil {
		handleDBError(c, err, "Could not retrieve products")
		return
	}

	c.JSON(http.StatusOK, products)
}

func DeleteProduct(c *gin.Context) {
	productId, err := parseProductID(c)
	if err != nil {
		return
	}

	// Attempt to delete the product from the database
	result := database.DB.Delete(&models.Product{}, productId)

	// Check if the product was found and deleted
	if result.Error != nil {
		handleDBError(c, result.Error, "Could not delete product")
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func UpdateProduct(c *gin.Context) {
	productId, err := parseProductID(c)
	if err != nil {
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
	if !bindJSON(c, &input) {
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
			handleDBError(c, err, "Could not update product")
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
