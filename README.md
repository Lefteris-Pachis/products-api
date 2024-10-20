# Products API

This is a RESTful API for managing products, built with Go, Gin, and GORM.


## Table of Contents

- [How to run](#how-to-run)
- [How to test](#how-to-test)
- [API Endpoints](#api-endpoints)


## How to run
To run the application using Docker Compose:

1. Make sure you have Docker and Docker Compose installed
2. Run the following commands:
```sh
cp .env.example .env
docker-compose up -d --build
```

You can access the endpoints on ``http://localhost:8080``

There is also a postman collection: `products-api.postman_collection.json` that you can import.

## How to test
To run the tests you have to run the following commands:
```sh
chmod +x run_tests.sh
./run_tests.sh
```

## API Endpoints
- `GET /products`: List all products (with pagination)
- `GET /products/:id`: Get a specific product
- `POST /products`: Create a new product
- `PATCH /products/:id`: Update an existing product
- `DELETE /products/:id`: Delete a product
