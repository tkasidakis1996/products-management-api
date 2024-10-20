package api

import (
	"github.com/gofiber/fiber/v2"
	"simpler-go-home-test/data_layer"
	"log"
	"strconv"
	"time"
	"errors"
	"gorm.io/gorm"
)

type UpdateProductNameRequest struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UpdateProductPriceRequest struct {
	ID    int     `json:"id"`
	Price float64 `json:"price"`
}

type ProductResponse struct {
	ID        uint      `json:"id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}


func InsertProduct(c *fiber.Ctx) error {

	products_db, err := data_layer.ProductsDB()

	defer func() {
		sqlDB, _ := products_db.DB()
		sqlDB.Close()
	}()
	
	if err != nil {
		
		log.Printf("Failed to connect to the products database: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error": "Failed to connect to the products database",})
	}

	product := data_layer.Product{} 

	err = c.BodyParser(&product); 

	if(err != nil) {
		
		log.Printf("Cannot parse JSON: %v", err)
		
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "Cannot parse JSON",})
	}

	if product.Name == "" || product.Price <= 0 {

		log.Printf("Invalid product data for insertion: name must be non-empty, and price must be greater than zero")

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "Invalid product data for insertion: name must be non-empty, and price must be greater than zero",})
	}

	log.Println("Inserting product (name : ", product.Name, ", price : ", product.Price, ") to the products database")

	productID, err := data_layer.InsertProduct(products_db, product.Name, product.Price)
	
	if err != nil {
		
		log.Printf("Failed to insert product at the products database %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to insert product at the products database",})
	}

	log.Println("Product (name : ", product.Name, ", price : ", product.Price, ") inserted successfully to the products database. Product ID : ", productID,)

	return c.JSON(fiber.Map{"message": "Product inserted successfully to the products database","product_id": productID,})
}

func DeleteProduct(c *fiber.Ctx) error {

	products_db, err := data_layer.ProductsDB()

	defer func() {
		sqlDB, _ := products_db.DB()
		sqlDB.Close()
	}()
	
	if err != nil {
		
		log.Printf("Failed to connect to the products database: %v", err)
		
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error": "Failed to connect to the products database",})
	}

	idParam := c.Params("id")
	
	productID, err := strconv.Atoi(idParam)
	
	if err != nil {
		
		log.Printf("Invalid product ID: %v", err)
		
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "Invalid product ID. Please provide a valid ID",})
	}

	log.Println("Attempting to delete product with ID:", productID)

	err = data_layer.DeleteProduct(products_db, productID)

	if errors.Is(err, gorm.ErrRecordNotFound) {

		log.Printf("Product with ID %d not found", productID)
		
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"Error": "Product not found",})
	}
	
	if err != nil {
		
		log.Printf("Failed to delete product with ID %d: %v", productID, err)
		
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error": "Failed to delete product from the products database",})
	}

	log.Printf("Product with ID %d deleted successfully from the products database", productID)

	return c.JSON(fiber.Map{"message": "Product deleted successfully from the products database","product_id": productID,})
}

func UpdateProductName(c *fiber.Ctx) error {

	products_db, err := data_layer.ProductsDB()

	defer func() {
		sqlDB, _ := products_db.DB()
		sqlDB.Close()
	}()
	
	if err != nil {
		log.Printf("Failed to connect to the products database: %v", err)
		
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error": "Failed to connect to the products database",})
	}

	requestBody := UpdateProductNameRequest{}
	
	err = c.BodyParser(&requestBody) 

	if(err != nil) {
		
		log.Printf("Cannot parse JSON: %v", err)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "Cannot parse JSON",})
	}

	if (requestBody.ID <= 0 || requestBody.Name == "") {
		
		log.Printf("Invalid product data for update: ID must be positive and name must be non-empty")
		
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "Invalid product data for update: ID must be positive and name must be non-empty",})
	}

	log.Printf("Attempting to update the name of product with ID: %d to '%s'", requestBody.ID, requestBody.Name)

	err = data_layer.UpdateProductName(products_db, requestBody.ID, requestBody.Name)

	if errors.Is(err, gorm.ErrRecordNotFound) {

		log.Printf("Product with ID %d not found", requestBody.ID)
		
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"Error": "Product not found",})
	}
	
	if err != nil {
		
		log.Printf("Failed to update product name for ID %d: %v", requestBody.ID, err)
		
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error": "Failed to update product name in the products database",})
	}

	log.Printf("Product with ID %d updated successfully. New name: '%s'", requestBody.ID, requestBody.Name)

	return c.JSON(fiber.Map{"message": "Product name updated successfully", "product_id": requestBody.ID, "new_name": requestBody.Name,})
}

func UpdateProductPrice(c *fiber.Ctx) error {
	
	products_db, err := data_layer.ProductsDB()

	defer func() {
		sqlDB, _ := products_db.DB()
		sqlDB.Close()
	}()
	
	if err != nil {
		
		log.Printf("Failed to connect to the products database: %v", err)
		
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error": "Failed to connect to the products database",})
	}

	requestBody := UpdateProductPriceRequest{}
	
	err = c.BodyParser(&requestBody); 

	if(err != nil) {
		
		log.Printf("Cannot parse JSON: %v", err)
		
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "Cannot parse JSON",})
	}

	if (requestBody.ID <= 0 || requestBody.Price <= 0) {
		
		log.Printf("Invalid product data for update: ID must be positive and price must be greater than zero")
		
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "Invalid product data for update: ID must be positive and price must be greater than zero",})
	}

	log.Printf("Attempting to update the price of product with ID: %d to %.2f", requestBody.ID, requestBody.Price)

	err = data_layer.UpdateProductPrice(products_db, requestBody.ID, requestBody.Price)

	if errors.Is(err, gorm.ErrRecordNotFound) {

		log.Printf("Product with ID %d not found", requestBody.ID)
		
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"Error": "Product not found",})
	}
	
	if (err != nil) {
		
		log.Printf("Failed to update product price for ID %d: %v", requestBody.ID, err)
		
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error": "Failed to update product price in the products database",})
	}

	log.Printf("Product with ID %d updated successfully. New price: %.2f", requestBody.ID, requestBody.Price)

	return c.JSON(fiber.Map{"message": "Product price updated successfully", "product_id": requestBody.ID,"new_price":  requestBody.Price,})
}

func RetrieveProduct(c *fiber.Ctx) error {

	products_db, err := data_layer.ProductsDB()

	defer func() {
		sqlDB, _ := products_db.DB()
		sqlDB.Close()
	}()
	
	if err != nil {
		
		log.Printf("Failed to connect to the products database: %v", err)
		
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error": "Failed to connect to the products database",})
	}

	idParam := c.Params("id")
	
	productID, err := strconv.Atoi(idParam)
	
	if err != nil {
		
		log.Printf("Invalid product ID: %v", err)
		
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "Invalid product ID. Please provide a valid ID",})
	}

	log.Printf("Attempting to retrieve product with ID: %d", productID)

	product, err := data_layer.RetrieveProduct(products_db, productID)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Handle the case where the product is not found
		log.Printf("Product with ID %d not found", productID)
		
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"Error": "Product not found",})
	}
	
	if err != nil {
		
		log.Printf("Failed to retrieve product from the products database: %v", err)
		
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error": "Failed to retrieve product from the products database",})
	}	

	productResponse := ProductResponse{ID: product.ID, Name: product.Name, Price: product.Price, CreatedAt: product.CreatedAt, UpdatedAt: product.UpdatedAt}

	log.Printf("Product with ID %d retrieved successfully: Name: %s, Price: %.2f", productID, productResponse.Name, productResponse.Price)

	return c.JSON(productResponse)
}

func RetrieveProductsWithPagination(c *fiber.Ctx) error {

	products_db, err := data_layer.ProductsDB()

	defer func() {
		sqlDB, _ := products_db.DB()
		sqlDB.Close()
	}()
	
	if err != nil {
		log.Printf("Failed to connect to the products database: %v", err)
		
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error": "Failed to connect to the products database",})
	}

	pageParam := c.Query("page", "1")    // Default page is 1
	limitParam := c.Query("limit", "10") // Default limit is 10

	page, err := strconv.Atoi(pageParam)
	
	if err != nil || page <= 0 {
		log.Printf("Invalid page number: %v", err)
		
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "Invalid page number. Must be a positive integer",})
	}

	limit, err := strconv.Atoi(limitParam)
	
	if err != nil || limit <= 0 {
		
		log.Printf("Invalid limit number: %v", err)
		
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "Invalid limit number. Must be a positive integer",})
	}

	offset := (page - 1) * limit

	log.Printf("Attempting to retrieve products with pagination: page = %d, limit = %d, offset = %d", page, limit, offset)

	products, err := data_layer.RetrieveProductsWithPagination(products_db, offset, limit)
	
	if err != nil {
		
		log.Printf("Failed to retrieve products with pagination: %v", err)
		
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error": "Failed to retrieve products with pagination",})
	}

	total_number_of_products, err := data_layer.GetTotalNumberOfProducts(products_db)

	if(err != nil) {

		log.Printf("Failed to retrieve total number of products: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error": "Failed to retrieve total number of products",})
		
	}
	
	totalPages := (int(total_number_of_products) + limit - 1) / limit
	
	var paginatedResponse []ProductResponse
	
	for _, product := range products {
		
		paginatedResponse = append(paginatedResponse, ProductResponse{
			ID:        product.ID,
			Name:      product.Name,
			Price:     product.Price,
			CreatedAt: product.CreatedAt,
			UpdatedAt: product.UpdatedAt,
		})
	}

	metadata := fiber.Map{"current_page": page, "total_pages":  totalPages, "total_number_of_products": total_number_of_products, "next_page":    page + 1, "prev_page":    page - 1}

	if page == 1 {
		metadata["prev_page"] = nil
	}
	if page >= totalPages {
		metadata["next_page"] = nil
	}

	log.Printf("Successfully retrieved %d products on page %d with limit %d", len(paginatedResponse), page, limit)

	return c.JSON(fiber.Map{"metadata": metadata, "products": paginatedResponse,})
}