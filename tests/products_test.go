package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"simpler-go-home-test/api"
	"simpler-go-home-test/data_layer"
	"testing"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"fmt"
)

// Setup function for initializing Fiber app
func SetupApp() (*fiber.App) {
	
	app := fiber.New()

	app.Post("/insert-product", api.InsertProduct)

	app.Delete("/delete-product/:id", api.DeleteProduct)

	app.Put("/update-product-name", api.UpdateProductName)

	app.Put("/update-product-price", api.UpdateProductPrice)
	
	app.Get("/retrieve-product/:id", api.RetrieveProduct)

	app.Get("/retrieve-products", api.RetrieveProductsWithPagination)

	return app
}



func TestInsertProduct_HappyPath(t *testing.T) {
	// Arrange
	app := SetupApp()

	defer data_layer.DestroyProductsDB()

	product := map[string]interface{}{
		"name":  "Laptop_Test",
		"price": 1999.99,
	}
	
	body, _ := json.Marshal(product)

	// Act
	req := httptest.NewRequest(http.MethodPost, "/insert-product", bytes.NewReader(body))
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, _ := app.Test(req, -1)

	// Assert
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Check the response data
	var responseData map[string]interface{}
	
	json.NewDecoder(resp.Body).Decode(&responseData)
	
	assert.Equal(t, "Product inserted successfully to the products database", responseData["message"])
	
	assert.NotNil(t, responseData["product_id"])
}

func TestInsertProduct_InvalidBody(t *testing.T) {
	// Arrange
	app := SetupApp()
	
	defer data_layer.DestroyProductsDB()

	// Invalid product data (missing name field)
	invalidProduct := map[string]interface{}{
		"price": 1999.99,
	}
	
	body, _ := json.Marshal(invalidProduct)

	// Act
	req := httptest.NewRequest(http.MethodPost, "/insert-product", bytes.NewReader(body))
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, _ := app.Test(req, -1)

	// Assert
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Check the response data
	var responseData map[string]interface{}
	
	json.NewDecoder(resp.Body).Decode(&responseData)
	
	assert.Equal(t,"Invalid product data for insertion: name must be non-empty, and price must be greater than zero", responseData["Error"])
}


func TestInsertProduct_InvalidPrice(t *testing.T) {
	
	// Arrange
	app := SetupApp()

	defer data_layer.DestroyProductsDB()

	// Invalid product data (negative price)
	invalidProduct := map[string]interface{}{"name":  "Laptop_Test_Invalid", "price": -100,}

	body, _ := json.Marshal(invalidProduct)

	// Act
	req := httptest.NewRequest(http.MethodPost, "/insert-product", bytes.NewReader(body))
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, _ := app.Test(req, -1)

	
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	
	var responseData map[string]interface{}
	
	json.NewDecoder(resp.Body).Decode(&responseData)

	// Assert
	assert.Equal(t,"Invalid product data for insertion: name must be non-empty, and price must be greater than zero", responseData["Error"])
}

// Test the happy path (valid product retrieval)
func TestRetrieveProduct_HappyPath(t *testing.T) {
	// Arrange
	app := SetupApp()

	defer data_layer.DestroyProductsDB()

	// Insert a product first, so that we can retrieve it
	product := map[string]interface{}{
		"name":  "Laptop_Test_Retrieve",
		"price": 1999.99,
	}

	body, _ := json.Marshal(product)
	req := httptest.NewRequest(http.MethodPost, "/insert-product", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)

	// Extract product ID from insertion response
	var insertResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&insertResponse)
	productID := int(insertResponse["product_id"].(float64))

	// Act - Now, retrieve the product by its ID
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/retrieve-product/%d", productID), nil)
	resp, _ = app.Test(req, -1)

	// Assert
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseData map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&responseData)

	assert.Equal(t, "Laptop_Test_Retrieve", responseData["name"])
	assert.Equal(t, 1999.99, responseData["price"])
	assert.NotNil(t, responseData["created_at"])
	assert.NotNil(t, responseData["updated_at"])
}


// Test invalid ID format (non-integer ID)
func TestRetrieveProduct_InvalidIDFormat(t *testing.T) {
	// Arrange
	app := SetupApp()

	defer data_layer.DestroyProductsDB()

	// Act - Try to retrieve a product with an invalid ID (e.g., "abc")
	req := httptest.NewRequest(http.MethodGet, "/retrieve-product/abc", nil)
	resp, _ := app.Test(req, -1)

	// Assert
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var responseData map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&responseData)

	assert.Equal(t, "Invalid product ID. Please provide a valid ID", responseData["Error"])
}


// Test product not found (valid but non-existent ID)
func TestRetrieveProduct_IDNotFound(t *testing.T) {
	// Arrange
	app:= SetupApp()

	// Act - Try to retrieve a product with a valid but non-existent ID (e.g., ID 999)
	req := httptest.NewRequest(http.MethodGet, "/retrieve-product/999", nil)
	
	resp, _ := app.Test(req, -1)

	// Assert
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	var responseData map[string]interface{}
	
	json.NewDecoder(resp.Body).Decode(&responseData)

	assert.Equal(t, "Product not found", responseData["Error"])
}

// Test the happy path (valid product deletion)
func TestDeleteProduct_HappyPath(t *testing.T) {
	// Arrange
	app := SetupApp()

	defer data_layer.DestroyProductsDB()

	// Insert a product first, so that we can delete it
	product := map[string]interface{}{
		"name":  "Laptop_Test_Delete",
		"price": 1500.00,
	}

	body, _ := json.Marshal(product)
	req := httptest.NewRequest(http.MethodPost, "/insert-product", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)

	// Extract product ID from insertion response
	var insertResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&insertResponse)
	productID := int(insertResponse["product_id"].(float64))

	// Act - Now, delete the product by its ID
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/delete-product/%d", productID), nil)
	resp, _ = app.Test(req, -1)

	// Assert
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseData map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&responseData)

	assert.Equal(t, "Product deleted successfully from the products database", responseData["message"])
	assert.Equal(t, float64(productID), responseData["product_id"]) // JSON unmarshal converts numbers to float64
}


// Test invalid ID format (non-integer ID)
func TestDeleteProduct_InvalidIDFormat(t *testing.T) {
	// Arrange
	app := SetupApp()

	defer data_layer.DestroyProductsDB()

	// Act - Try to delete a product with an invalid ID (e.g., "abc")
	req := httptest.NewRequest(http.MethodDelete, "/delete-product/abc", nil)
	resp, _ := app.Test(req, -1)

	// Assert
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var responseData map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&responseData)

	assert.Equal(t, "Invalid product ID. Please provide a valid ID", responseData["Error"])
}


// Test product not found (valid but non-existent ID)
func TestDeleteProduct_IDNotFound(t *testing.T) {
	// Arrange
	app := SetupApp()

	defer data_layer.DestroyProductsDB()

	// Act - Try to delete a product with a valid but non-existent ID (e.g., ID 999)
	req := httptest.NewRequest(http.MethodDelete, "/delete-product/999", nil)
	
	resp, _ := app.Test(req, -1)

	// Assert
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	var responseData map[string]interface{}
	
	json.NewDecoder(resp.Body).Decode(&responseData)

	assert.Equal(t, "Product not found", responseData["Error"])
}

// Test the happy path (valid price update)
func TestUpdateProductPrice_HappyPath(t *testing.T) {
	// Arrange
	app := SetupApp()

	defer data_layer.DestroyProductsDB()

	// Insert a product first, so that we can update its price
	product := map[string]interface{}{
		"name":  "Laptop_Test_Update",
		"price": 1500.00,
	}

	body, _ := json.Marshal(product)
	
	req := httptest.NewRequest(http.MethodPost, "/insert-product", bytes.NewReader(body))
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, _ := app.Test(req, -1)

	// Extract product ID from insertion response
	var insertResponse map[string]interface{}
	
	json.NewDecoder(resp.Body).Decode(&insertResponse)
	
	productID := int(insertResponse["product_id"].(float64))

	// Act - Now, update the product's price
	updatePrice := map[string]interface{}{
		"id":    productID,
		"price": 1800.00,
	}
	
	updateBody, _ := json.Marshal(updatePrice)
	
	req = httptest.NewRequest(http.MethodPut, "/update-product-price", bytes.NewReader(updateBody))
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, _ = app.Test(req, -1)

	// Assert
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseData map[string]interface{}
	
	json.NewDecoder(resp.Body).Decode(&responseData)

	assert.Equal(t, "Product price updated successfully", responseData["message"])
	
	assert.Equal(t, float64(productID), responseData["product_id"])
	
	assert.Equal(t, 1800.00, responseData["new_price"])
}



// Test invalid price update (negative price)
func TestUpdateProductPrice_InvalidPrice(t *testing.T) {
	// Arrange
	app := SetupApp()

	defer data_layer.DestroyProductsDB()

	// Insert a product first
	product := map[string]interface{}{
		"name":  "Laptop_Test_Invalid_Update",
		"price": 1500.00,
	}

	body, _ := json.Marshal(product)
	
	req := httptest.NewRequest(http.MethodPost, "/insert-product", bytes.NewReader(body))
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, _ := app.Test(req, -1)

	// Extract product ID
	var insertResponse map[string]interface{}
	
	json.NewDecoder(resp.Body).Decode(&insertResponse)
	
	productID := int(insertResponse["product_id"].(float64))

	// Act - Try to update the product's price with an invalid (negative) price
	updatePrice := map[string]interface{}{
		"id":    productID,
		"price": -500.00,
	}

	updateBody, _ := json.Marshal(updatePrice)
	
	req = httptest.NewRequest(http.MethodPut, "/update-product-price", bytes.NewReader(updateBody))
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, _ = app.Test(req, -1)

	// Assert
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var responseData map[string]interface{}
	
	json.NewDecoder(resp.Body).Decode(&responseData)

	assert.Equal(t, "Invalid product data for update: ID must be positive and price must be greater than zero", responseData["Error"])
}


// Test product not found (valid ID but non-existent product)
func TestUpdateProductPrice_IDNotFound(t *testing.T) {
	// Arrange
	app := SetupApp()

	defer data_layer.DestroyProductsDB()

	// Act - Try to update the price of a product that doesn't exist (e.g., ID 999)
	updatePrice := map[string]interface{}{
		"id":    999,
		"price": 2000.00,
	}
	updateBody, _ := json.Marshal(updatePrice)
	req := httptest.NewRequest(http.MethodPut, "/update-product-price", bytes.NewReader(updateBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)

	// Assert
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	var responseData map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&responseData)

	assert.Equal(t, "Product not found", responseData["Error"])
}


func TestUpdateProductName_HappyPath(t *testing.T) {
	// Arrange
	app := SetupApp()

	defer data_layer.DestroyProductsDB()

	// Insert a product first, so that we can update its name
	product := map[string]interface{}{
		"name":  "Laptop_Test_Update",
		"price": 1500.00,
	}
	
	body, _ := json.Marshal(product)
	
	req := httptest.NewRequest(http.MethodPost, "/insert-product", bytes.NewReader(body))
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, _ := app.Test(req, -1)

	// Extract product ID from insertion response
	var insertResponse map[string]interface{}
	
	json.NewDecoder(resp.Body).Decode(&insertResponse)
	
	productID := int(insertResponse["product_id"].(float64))

	// Act - Now, update the product's name
	updateName := map[string]interface{}{
		"id":   productID,
		"name": "Laptop_Test_Updated",
	}
	
	updateBody, _ := json.Marshal(updateName)
	
	req = httptest.NewRequest(http.MethodPut, "/update-product-name", bytes.NewReader(updateBody))
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, _ = app.Test(req, -1)

	// Assert
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseData map[string]interface{}
	
	json.NewDecoder(resp.Body).Decode(&responseData)

	assert.Equal(t, "Product name updated successfully", responseData["message"])
	assert.Equal(t, float64(productID), responseData["product_id"])
	assert.Equal(t, "Laptop_Test_Updated", responseData["new_name"])
}


// Test invalid name update (empty name)
func TestUpdateProductName_InvalidName(t *testing.T) {
	// Arrange
	app := SetupApp()

	defer data_layer.DestroyProductsDB()

	// Insert a product first
	product := map[string]interface{}{
		"name":  "Laptop_Test_Invalid_Update",
		"price": 1500.00,
	}
	body, _ := json.Marshal(product)
	req := httptest.NewRequest(http.MethodPost, "/insert-product", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)

	// Extract product ID
	var insertResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&insertResponse)
	productID := int(insertResponse["product_id"].(float64))

	// Act - Try to update the product's name with an invalid (empty) name
	updateName := map[string]interface{}{
		"id":   productID,
		"name": "",
	}
	updateBody, _ := json.Marshal(updateName)
	req = httptest.NewRequest(http.MethodPut, "/update-product-name", bytes.NewReader(updateBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req, -1)

	// Assert
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var responseData map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&responseData)

	assert.Equal(t, "Invalid product data for update: ID must be positive and name must be non-empty", responseData["Error"])
}


// Test product not found (valid ID but non-existent product)
func TestUpdateProductName_IDNotFound(t *testing.T) {
	// Arrange
	app := SetupApp()

	defer data_layer.DestroyProductsDB()

	// Act - Try to update the name of a product that doesn't exist (e.g., ID 999)
	updateName := map[string]interface{}{
		"id":   999,
		"name": "NonExistentProduct",
	}
	updateBody, _ := json.Marshal(updateName)
	req := httptest.NewRequest(http.MethodPut, "/update-product-name", bytes.NewReader(updateBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)

	// Assert
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	var responseData map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&responseData)

	assert.Equal(t, "Product not found", responseData["Error"])
}


func TestRetrieveProductsWithPagination_HappyPath(t *testing.T) {
	// Arrange
	app := SetupApp()

	defer data_layer.DestroyProductsDB()

	// Insert multiple products (e.g., 30 products)
	for i := 1; i <= 30; i++ {
		product := map[string]interface{}{
			"name":  fmt.Sprintf("Product_%d", i),
			"price": float64(i * 100),
		}
		body, _ := json.Marshal(product)
		req := httptest.NewRequest(http.MethodPost, "/insert-product", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		app.Test(req, -1)
	}

	// Act - Request page 2 with a limit of 10 products per page
	req := httptest.NewRequest(http.MethodGet, "/retrieve-products?page=2&limit=10", nil)
	resp, _ := app.Test(req, -1)

	// Assert
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseData map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&responseData)

	// Check the pagination metadata
	metadata := responseData["metadata"].(map[string]interface{})
	assert.Equal(t, float64(2), metadata["current_page"])
	assert.Equal(t, float64(3), metadata["total_pages"])
	assert.Equal(t, float64(30), metadata["total_number_of_products"])

	// Check the products returned (should be products 11 to 20)
	products := responseData["products"].([]interface{})
	assert.Equal(t, 10, len(products))

	// Validate that the first product on page 2 is "Product_11"
	firstProduct := products[0].(map[string]interface{})
	assert.Equal(t, "Product_11", firstProduct["name"])
}


// Test invalid pagination (negative page or limit)
func TestRetrieveProductsWithPagination_InvalidPagination(t *testing.T) {
	// Arrange
	app := SetupApp()

	defer data_layer.DestroyProductsDB()

	// Act - Request an invalid page (-1) and limit (0)
	req := httptest.NewRequest(http.MethodGet, "/retrieve-products?page=-1&limit=0", nil)
	resp, _ := app.Test(req, -1)

	// Assert
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var responseData map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&responseData)

	// Check error message for invalid page and limit
	assert.Equal(t, "Invalid page number. Must be a positive integer", responseData["Error"])
}

func TestRetrieveProductsWithPagination_PageExceedsData(t *testing.T) {
	// Arrange
	app := SetupApp()

	defer data_layer.DestroyProductsDB()

	// Insert 5 products (fewer than 10)
	for i := 1; i <= 5; i++ {
		product := map[string]interface{}{
			"name":  fmt.Sprintf("Product_%d", i),
			"price": float64(i * 100),
		}
		body, _ := json.Marshal(product)
		req := httptest.NewRequest(http.MethodPost, "/insert-product", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		app.Test(req, -1)
	}

	// Act - Request page 2 with a limit of 10 (no data on page 2)
	req := httptest.NewRequest(http.MethodGet, "/retrieve-products?page=2&limit=10", nil)
	resp, _ := app.Test(req, -1)

	// Assert
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseData map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&responseData)

	// Check the metadata for page 2
	metadata := responseData["metadata"].(map[string]interface{})
	assert.Equal(t, float64(2), metadata["current_page"])
	assert.Equal(t, float64(1), metadata["total_pages"]) // Only 1 page of data
	assert.Equal(t, float64(5), metadata["total_number_of_products"])

	// Check for the 'products' field in the response
	products, ok := responseData["products"]
	if !ok || products == nil {
		// If 'products' field is missing or nil, assert that this is correct (no products on this page)
		assert.Nil(t, products, "Expected no products on this page, but found some")
	} else {
		// If 'products' field exists, check that it is an array and its length is 0
		productArray, ok := products.([]interface{})
		if ok {
			assert.Equal(t, 0, len(productArray), "Expected 0 products on this page, but found some")
		} else {
			t.Fatalf("Expected 'products' to be an array, but got a different type")
		}
	}
}