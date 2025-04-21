#!/bin/bash

# Create a script to build and run the application locally
set -e

echo "Building and deploying Hustler Trading Bot locally..."

# Build the backend
echo "Building backend..."
cd /Users/dagmfekaduyenealem/Documents/ai/hustler-trading-bot
go build -o bin/hustler ./cmd/hustler

# Build the frontend
echo "Building frontend..."
cd /Users/dagmfekaduyenealem/Documents/ai/hustler-trading-bot/web/frontend
npm install
npm run build

# Start PostgreSQL with Docker
echo "Starting PostgreSQL..."
docker run -d --name hustler-postgres \
  -e POSTGRES_PASSWORD=hustlerpass \
  -e POSTGRES_USER=hustler \
  -e POSTGRES_DB=hustler \
  -p 5432:5432 \
  postgres:14

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
sleep 10

# Initialize the database
echo "Initializing database..."
docker exec -i hustler-postgres psql -U hustler -d hustler < /Users/dagmfekaduyenealem/Documents/ai/hustler-trading-bot/db/schema.sql

# Start the backend in the background
echo "Starting backend..."
cd /Users/dagmfekaduyenealem/Documents/ai/hustler-trading-bot
./bin/hustler &
BACKEND_PID=$!

# Start the frontend in the background
echo "Starting frontend..."
cd /Users/dagmfekaduyenealem/Documents/ai/hustler-trading-bot/web/frontend
npx serve -s build -l 3000 &
FRONTEND_PID=$!

echo "Application is running!"
echo "Frontend: http://localhost:3000"
echo "Backend: http://localhost:8080"
echo ""
echo "Press Ctrl+C to stop the application"

# Function to clean up on exit
function cleanup {
  echo "Stopping application..."
  kill $FRONTEND_PID
  kill $BACKEND_PID
  docker stop hustler-postgres
  docker rm hustler-postgres
  echo "Application stopped"
}

# Register the cleanup function to be called on exit
trap cleanup EXIT

# Wait for user to press Ctrl+C
wait $FRONTEND_PID
