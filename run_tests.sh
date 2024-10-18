#!/bin/bash

# Start the test database
docker-compose -f docker-compose.test.yml up -d

# Wait for the database to be ready
until docker-compose -f docker-compose.test.yml exec -T test-db pg_isready -U testuser -d testdb
do
  echo "Waiting for database connection..."
  sleep 2
done

# Run the tests
go test ./... -v

# Stop the test database
docker-compose -f docker-compose.test.yml down -v