package data_layer

import (
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"log"
)

// Product model definition
type Product struct {
	gorm.Model
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}



func ProductsDB() (*gorm.DB, error) {
	
	products_db, err := gorm.Open(sqlite.Open("simpler_home_test.db"), &gorm.Config{})
	
	if err != nil {

		return nil, err
	}

	err = products_db.AutoMigrate(&Product{})
	
	if err != nil {

		return nil, err
	}

	return products_db, nil
}

func DestroyProductsDB() {

	// Remove the database file (for SQLite)
	err := os.Remove("simpler_home_test.db")
	
	if err != nil {
		
		log.Printf("Error removing database file: %v", err)
	
	} else {
		
		log.Println("Database file removed successfully")
	}
}

func InsertProduct(products_db *gorm.DB, name string, price float64) (uint, error) {
	
	product := Product{Name: name, Price: price}
	
	result := products_db.Create(&product)
	
	if result.Error != nil {

		return 0, result.Error
	}

	return product.ID, nil
}

func DeleteProduct(products_db *gorm.DB, id int) error {

	var product Product

	result := products_db.Unscoped().Delete(&product, id)

	if result.Error != nil {
		
		return result.Error
	
	}

	if result.RowsAffected == 0 {
		
		return gorm.ErrRecordNotFound
	
	}

	return nil
	
}

func UpdateProductName(products_db *gorm.DB, id int, name string) error {

	var product Product

	result := products_db.Model(&product).Where("id = ?", id).Update("name", name)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func UpdateProductPrice(products_db *gorm.DB, id int, price float64) error {

	var product Product

	result := products_db.Model(&product).Where("id = ?", id).Update("price", price)

	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func RetrieveProduct(products_db *gorm.DB, id int) (Product, error) {

	var product Product

	result := products_db.First(&product, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return Product{}, gorm.ErrRecordNotFound
	}

	// Return any other errors
	if result.Error != nil {
		return Product{}, result.Error
	}

	return product, nil
}

func RetrieveProductsWithPagination(products_db *gorm.DB, offset int, limit int) ([]Product, error) {

	var products []Product

	result := products_db.Limit(limit).Offset(offset).Find(&products)

	if result.Error != nil {

		return nil, result.Error
	}

	return products, nil
}

func GetTotalNumberOfProducts(products_db *gorm.DB) (int64, error) {

	var totalRecords int64

	var product Product

	result := products_db.Model(product).Count(&totalRecords)

	if result.Error != nil {

		return -1, result.Error
	}

	return totalRecords, nil
}

