#!/bin/bash

# Script to deploy the Hustler Trading Bot to Minikube

# Set variables
APP_NAME="hustler-trading-bot"
NAMESPACE="default"

# Make scripts executable
chmod +x ./scripts/verify-minikube.sh

# Verify Minikube environment
echo "Verifying Minikube environment..."
./scripts/verify-minikube.sh

# Build Docker image
echo "Building Docker image..."
eval $(minikube docker-env)
docker build -t ${APP_NAME}:latest .

# Apply Kubernetes manifests
echo "Applying Kubernetes manifests..."

# Create ConfigMap and Secrets
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secrets.yaml

# Deploy PostgreSQL
echo "Deploying PostgreSQL..."
kubectl apply -f k8s/postgres.yaml

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
kubectl wait --for=condition=ready pod -l app=hustler-postgres --timeout=120s

# Deploy Trading Bot
echo "Deploying Trading Bot..."
kubectl apply -f k8s/trading-bot.yaml

# Wait for Trading Bot to be ready
echo "Waiting for Trading Bot to be ready..."
kubectl wait --for=condition=ready pod -l app=hustler-trading-bot --timeout=120s

# Get service URL
echo "Getting service URL..."
minikube service hustler-trading-bot --url

echo "Deployment completed successfully!"
echo "You can access the Trading Bot UI using the URL above."
echo "To view logs, run: kubectl logs -f deployment/hustler-trading-bot"
