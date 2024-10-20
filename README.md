#Simpler Go Home Test - Product Management API

Welcome to the Simpler Go Home Test project! This repository contains a Go-based RESTful API that allows for managing products. The API is built using Fiber, GORM, and SQLite, and includes complete CRUD (Create, Read, Update, Delete) functionality along with pagination support.

Features
Fiber Framework: Simple and fast web framework in Go.
GORM ORM: For seamless interaction with SQLite.
SQLite: Lightweight and embedded database used for local testing.
Dockerized: Run everything inside a Docker container.
Unit Tested: Extensive unit testing to ensure everything works as expected.
How to Build the Docker Image

Clone the repository:
git clone https://github.com/tkasidakis1996/simpler-go-home-test.git
cd simpler-go-home-test

Build the Docker image:
sudo docker build -t simpler-go-home-test .

How to Run the API
Once the Docker image is built, you can run the API by creating a container named products-api:

sudo docker run -d --name products-api -p 8000:8000 simpler-go-home-test
This will start the API server on port 8000. You should see a message indicating that the server is up and running:
Products API is running on port 8000

Accessing the Docker Container & Running Tests
If you'd like to enter the Docker container and run the tests:

Access the container:
sudo docker exec -it products-api /bin/bash

Run the tests: Inside the container, you can run the tests like this:

go test ./... -v
This will run all the tests inside the container, including CRUD operations and pagination.

Inserting Multiple Products
If you want to quickly populate the database with 30 products, use the provided shell script insert_products.sh:

./insert_products.sh
This script will automatically call the /insert-product API 30 times, adding sample products to the database.

Example API Requests with Curl

Once the API is up and running, you can interact with it using curl. Below are some example commands:

Insert a Product:

curl -X POST http://localhost:8000/insert-product \
-H "Content-Type: application/json" \
-d '{"name": "Laptop", "price": 1500.50}'

Retrieve a Product by ID:
curl http://localhost:8000/retrieve-product/1

Update a Product Name:
curl -X PUT http://localhost:8000/update-product-name \
-H "Content-Type: application/json" \
-d '{"id": 1, "name": "Updated Laptop Name"}'


Update a Product Price:
curl -X PUT http://localhost:8000/update-product-price \
-H "Content-Type: application/json" \
-d '{"id": 1, "price": 1600.00}'

Delete a Product:
curl -X DELETE http://localhost:8000/delete-product/1

Retrieve Products with Pagination:
curl "http://localhost:8000/retrieve-products?page=1&limit=10"
