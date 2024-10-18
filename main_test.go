package main

import (
	"log"
	"os"
	"products-api/database"
	"products-api/models"
	"products-api/routes"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testDB *gorm.DB
var testRouter *gin.Engine

func TestMain(m *testing.M) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Setup Test Database
	var err error
	dsn := "host=localhost user=testuser password=testpass dbname=testdb port=5433 sslmode=disable"
	testDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to test database:", err)
	}

	// Migrate the schema
	err = testDB.AutoMigrate(&models.Product{})
	if err != nil {
		log.Fatal("Failed to migrate test database:", err)
	}

	// Set the global DB variable to our test DB
	database.DB = testDB

	// Setup the router
	testRouter = setupRouter()

	// Run the tests
	code := m.Run()

	// Exit
	os.Exit(code)
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	routes.SetupRoutes(r)
	return r
}
