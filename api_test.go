package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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

	cleanupProducts(t)
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

	cleanupProducts(t)
}

func TestGetProductsPagination(t *testing.T) {
	// Create a larger set of test products
	var testProducts []models.Product
	for i := 1; i <= 25; i++ {
		testProducts = append(testProducts, models.Product{
			Name:        fmt.Sprintf("Test Product %d", i),
			Description: fmt.Sprintf("Description %d", i),
			Price:       float64(i) * 10.0,
		})
	}

	createTestProducts(t, testProducts)

	// Test cases for different pagination scenarios
	testCases := []struct {
		name          string
		page          int
		limit         int
		expectedCount int
	}{
		{"First page with 10 items", 1, 10, 10},
		{"Second page with 10 items", 2, 10, 10},
		{"Third page with 10 items (only 5 left)", 3, 10, 5},
		{"First page with 20 items", 1, 20, 20},
		{"Second page with 20 items (only 5 left)", 2, 20, 5},
		{"Page beyond data range", 4, 10, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/products?page=%d&limit=%d", tc.page, tc.limit)
			req, _ := http.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()
			testRouter.ServeHTTP(w, req)

			// Check response
			assert.Equal(t, http.StatusOK, w.Code)

			var response GetProductsResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedCount, len(response.Data), "Unexpected number of products returned")
			assert.Equal(t, tc.page, response.Page, "Unexpected page number")
			assert.Equal(t, tc.limit, response.Limit, "Unexpected limit")
			assert.Equal(t, len(testProducts), response.Total, "Unexpected total count")

			// Check if the returned products are the correct ones for this page
			if len(response.Data) > 0 {
				firstProductID := response.Data[0].ID
				expectedFirstProductID := uint((tc.page-1)*tc.limit + 1)
				assert.Equal(t, expectedFirstProductID, firstProductID, "Unexpected first product ID for this page")
			}
		})
	}
	cleanupProducts(t)
}

func TestDeleteProduct(t *testing.T) {
	// Create a test product
	testProduct := models.Product{
		Name:        "Test Delete Product",
		Description: "This product will be deleted",
		Price:       99.99,
	}

	createdProductIDs := createTestProducts(t, []models.Product{testProduct})
	assert.Len(t, createdProductIDs, 1)
	productID := createdProductIDs[0]

	// Make delete request
	url := fmt.Sprintf("/products/%d", productID)
	req, _ := http.NewRequest("DELETE", url, nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	// Check response is 200
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Product deleted successfully", response["message"])

	// Verify product is deleted from database
	var deletedProduct models.Product
	result := database.DB.First(&deletedProduct, productID)
	assert.Error(t, result.Error)
	assert.Equal(t, "record not found", result.Error.Error())

	// Try to delete the same product again (should return not found)
	req, _ = http.NewRequest("DELETE", url, nil)
	w = httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	// Check response is 404
	assert.Equal(t, http.StatusNotFound, w.Code)
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Product not found", response["error"])

	cleanupProducts(t)
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

func cleanupProducts(t *testing.T) {
	result := database.DB.Exec("TRUNCATE TABLE products RESTART IDENTITY")
	assert.NoError(t, result.Error, "Failed to truncate products table")
}
