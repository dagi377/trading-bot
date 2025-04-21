#!/bin/bash

# Script to deploy the Hustler Trading Bot to Minikube
set -e

echo "Deploying Hustler Trading Bot to Minikube..."

# Check if Minikube is running
if ! minikube status &>/dev/null; then
  echo "Minikube is not running. Starting Minikube..."
  minikube start --cpus 4 --memory 4096
else
  echo "Minikube is already running"
fi

# Enable ingress addon if not already enabled
if ! minikube addons list | grep -q "ingress: enabled"; then
  echo "Enabling ingress addon..."
  minikube addons enable ingress
fi

# Set docker env to use Minikube's Docker daemon
echo "Setting Docker environment to use Minikube's Docker daemon..."
eval $(minikube docker-env)

# Build the backend Docker image
echo "Building backend Docker image..."
cd /Users/dagmfekaduyenealem/Documents/ai/hustler-trading-bot
docker build -t hustler-trading-bot:latest .

# Build the frontend Docker image
echo "Building frontend Docker image..."
cd /Users/dagmfekaduyenealem/Documents/ai/hustler-trading-bot/
docker build -t hustler-frontend:latest .

# Apply Kubernetes manifests
echo "Applying Kubernetes manifests..."
cd /Users/dagmfekaduyenealem/Documents/ai/hustler-trading-bot
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secrets.yaml
kubectl apply -f k8s/postgres.yaml
kubectl apply -f k8s/backend.yaml
kubectl apply -f k8s/frontend.yaml
kubectl apply -f k8s/ingress.yaml

# Wait for deployments to be ready
echo "Waiting for deployments to be ready..."
kubectl rollout status deployment/hustler-backend
kubectl rollout status deployment/hustler-frontend
kubectl rollout status deployment/postgres

# Get the URL
echo "Getting application URL..."
MINIKUBE_IP=$(minikube ip)
echo "Application is available at: http://$MINIKUBE_IP"
echo "You may need to add an entry to your /etc/hosts file:"
echo "$MINIKUBE_IP hustler.local"
echo "Then you can access the application at: http://hustler.local"

echo "Deployment completed successfully!"
