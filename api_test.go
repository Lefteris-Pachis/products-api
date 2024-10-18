package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"products-api/models"
	"testing"
)

// Define structs to match the responses
type CreateProductResponse struct {
	Message string         `json:"message"`
	Product models.Product `json:"product"`
}

func TestCreateProduct(t *testing.T) {
	// Test data
	product := models.Product{
		Name:        "Test Product",
		Description: "This is a test product",
		Price:       19.99,
	}
	jsonValue, _ := json.Marshal(product)

	// Create request
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(jsonValue))

	// Create response recorder
	w := httptest.NewRecorder()

	// Perform request
	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Parse response body
	var response CreateProductResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Assert response structure and values
	assert.Equal(t, "Product created successfully", response.Message)
	assert.NotZero(t, response.Product.ID)
	assert.Equal(t, product.Name, response.Product.Name)
	assert.Equal(t, product.Description, response.Product.Description)
	assert.Equal(t, product.Price, response.Product.Price)
	assert.NotZero(t, response.Product.CreatedAt)
	assert.NotZero(t, response.Product.UpdatedAt)
}
