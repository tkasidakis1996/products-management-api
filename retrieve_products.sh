#!/bin/bash

# Base URL for the retrieve products endpoint
BASE_URL="http://localhost:8000/retrieve-products"

# Total number of products
TOTAL_PRODUCTS=30
# Number of products per page
LIMIT=10
# Calculate total pages
TOTAL_PAGES=$((TOTAL_PRODUCTS / LIMIT))

# Loop to retrieve 10 products at a time
for page in $(seq 1 $TOTAL_PAGES)
do
  echo "Retrieving page $page..."
  
  # Make the HTTP GET request to retrieve the paginated products
  curl "$BASE_URL?page=$page&limit=$LIMIT"

  echo "Retrieved products on page $page"
done
