package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

// Config holds the database configuration
type Config struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     string
}

// ConnectDB initializes the database connection
func ConnectDB() {
	// Load database configuration from environment variables
	config := Config{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_PORT"),
	}

	// Validate essential fields
	if config.Host == "" || config.User == "" || config.Password == "" || config.Name == "" || config.Port == "" {
		log.Fatal("Database configuration is incomplete.")
	}

	// Create the Data Source Name (DSN)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.Host, config.User, config.Password, config.Name, config.Port)

	// Open the database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = db
}
