package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"products-api/database"
	"products-api/models"
	"testing"
)

// Define structs to match the responses
type CreateProductResponse struct {
	Message string         `json:"message"`
	Product models.Product `json:"product"`
}

type GetProductsResponse struct {
	Data  []models.Product `json:"data"`
	Total int              `json:"total"`
	Page  int              `json:"page"`
	Limit int              `json:"limit"`
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

func TestGetProducts(t *testing.T) {
	testProducts := []models.Product{
		{Name: "Test Product 1", Description: "Description 1", Price: 10.99},
		{Name: "Test Product 2", Description: "Description 2", Price: 20.99},
		{Name: "Test Product 3", Description: "Description 3", Price: 30.99},
	}

	createdProductIDs := createTestProducts(t, testProducts)

	// Make request to get products
	req, _ := http.NewRequest("GET", "/products?page=1&limit=10", nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	//t.Logf("Response body: %s", w.Body.String())

	var response GetProductsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.GreaterOrEqual(t, len(response.Data), len(testProducts))
	assert.GreaterOrEqual(t, response.Total, len(testProducts))
	assert.Equal(t, 1, response.Page)
	assert.Equal(t, 10, response.Limit)

	// Check if all created products are in the response
	foundProducts := 0
	for _, id := range createdProductIDs {
		for _, responseProduct := range response.Data {
			if responseProduct.ID == id {
				foundProducts++
				break
			}
		}
	}
	assert.Equal(t, len(createdProductIDs), foundProducts, "Not all created products were found in the response")
}

func createTestProducts(t *testing.T, products []models.Product) []uint {
	var createdIDs []uint

	for _, p := range products {
		result := database.DB.Unscoped().Create(&p)
		assert.NoError(t, result.Error)
		createdIDs = append(createdIDs, p.ID)
	}
	return createdIDs
}
