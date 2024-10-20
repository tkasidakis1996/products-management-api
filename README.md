# **Simpler Go Home Test - Product Management API**

Welcome to the **Simpler Go Home Test** project! This repository contains a Go-based RESTful API that allows for managing products. The API is built using **Fiber**, **GORM**, and **SQLite**, and includes complete **CRUD** (Create, Read, Update, Delete) functionality along with pagination support.

## **Features**

- **Fiber Framework**: Simple and fast web framework in Go.
- **GORM ORM**: For seamless interaction with SQLite.
- **SQLite**: Lightweight and embedded database used for local testing.
- **Dockerized**: Run everything inside a Docker container.
- **Unit Tested**: Extensive unit testing to ensure everything works as expected.

---

## **How to Build the Docker Image**

1. **Clone the repository**:
   ```bash
   git clone https://github.com/tkasidakis1996/simpler-go-home-test.git
   cd simpler-go-home-test

2. **Build the Docker image**:
   ```bash
   sudo docker build -t simpler-go-home-test .

## **How to Run the API**

Once the Docker image is built, you can run the API by creating a container named `products-api`:

  ```bash
  sudo docker run -d --name products-api -p 8000:8000 simpler-go-home-test
  ```


This will start the API server on port 8000. You should see a message indicating that the server is up and running:

  ```bash
  Products API is running on port 8000
  ```

## **Accessing the Docker Container & Running Tests**

If you'd like to enter the Docker container and run the tests:

  ```bash
  sudo docker exec -it products-api /bin/bash
  ```

Inside the container, you can run the tests like this:

  ```bash
  go test ./... -v
  ```

## **Inserting Multiple Products**

If you want to quickly populate the database with 30 products, use the provided shell script `insert_products.sh`:

  ```bash
  ./insert_products.sh
  ```
## **Example API Requests with Curl**

Once the API is up and running, you can interact with it using `curl`. Below are some example commands:

### **Insert a Product**
  ```bash
  curl -X POST http://localhost:8000/insert-product \
  -H "Content-Type: application/json" \
  -d '{"name": "Laptop", "price": 1500.50}'
  ```
### **Retrieve a Product by ID**
  ```bash
  curl http://localhost:8000/retrieve-product/1
  ```
### **Update a Product Name**
  ```bash
  curl -X PUT http://localhost:8000/update-product-name \
 -H "Content-Type: application/json" \
 -d '{"id": 1, "name": "Updated Laptop Name"}'
  ```
### **Update a Product Price**
  ```bash
  curl -X PUT http://localhost:8000/update-product-price \
 -H "Content-Type: application/json" \
 -d '{"id": 1, "price": 1600.00}'
  ```
### **Delete a Product**
  ```bash
  curl -X DELETE http://localhost:8000/delete-product/1
  ```

### **Retrieve Products with Pagination**
  ```bash
  curl http://localhost:8000/retrieve-products?page=1&limit=10
  ```




