#!/bin/bash

# Base URL for the insert product endpoint
BASE_URL="http://localhost:8000/insert-product"

# Loop to insert 30 products
for i in {1..30}
do
  # Define the product name and price
  PRODUCT_NAME="Laptop_$i"
  PRODUCT_PRICE=$((1500 + i * 10))

  # Make the HTTP POST request to insert the product
  curl -X POST "$BASE_URL" \
    -H "Content-Type: application/json" \
    -d "{\"name\": \"$PRODUCT_NAME\", \"price\": $PRODUCT_PRICE}"
  
  echo "Inserted: $PRODUCT_NAME with price $PRODUCT_PRICE"
done
