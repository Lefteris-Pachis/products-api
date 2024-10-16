package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"products-api/database"
	"products-api/models"
	"products-api/routes"
)

func main() {
	router := gin.Default()

	// Initialize database connection
	database.ConnectDB()

	// Perform migration
	migrateDatabase()

	router.Handle("GET", "/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Set up routes
	routes.SetupRoutes(router)

	// Start server on default port
	err := router.Run()
	if err != nil {
		log.Fatal("Failed to start the server:", err)
	}
}

// migrateDatabase performs database migrations
func migrateDatabase() {
	err := database.DB.AutoMigrate(&models.Product{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}
