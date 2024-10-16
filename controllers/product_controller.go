package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"products-api/database"
	"products-api/models"
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
