package main

import (
	"log"
	"github.com/gofiber/fiber/v2"
	"simpler-go-home-test/api"
)

func main() {

	products_api := fiber.New()

	products_api.Post("/insert-product", api.InsertProduct)

	products_api.Delete("/delete-product/:id", api.DeleteProduct)

	products_api.Put("/update-product-name", api.UpdateProductName)

	products_api.Put("/update-product-price", api.UpdateProductPrice)
	
	products_api.Get("/retrieve-product/:id", api.RetrieveProduct)

	products_api.Get("/retrieve-products", api.RetrieveProductsWithPagination)
	
	log.Println("Products API is running on port 8000")
	
	log.Fatal(products_api.Listen(":8000"))
}