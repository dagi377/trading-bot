#!/bin/bash

# Script to verify Minikube environment and deploy the Hustler Trading Bot

# Check if Minikube is running
echo "Checking Minikube status..."
MINIKUBE_STATUS=$(minikube status -f '{{.Host}}')

if [ "$MINIKUBE_STATUS" != "Running" ]; then
  echo "Minikube is not running. Starting Minikube..."
  minikube start --cpus=4 --memory=4096
else
  echo "Minikube is already running. Checking resources..."
  
  # Check allocated resources
  MINIKUBE_CPUS=$(minikube config view | grep cpus | awk '{print $3}')
  MINIKUBE_MEMORY=$(minikube config view | grep memory | awk '{print $3}')
  
  if [[ -z "$MINIKUBE_CPUS" || "$MINIKUBE_CPUS" -lt 4 ]]; then
    echo "Warning: Minikube CPU allocation may be insufficient. Recommended: 4 CPUs"
  else
    echo "Minikube CPU allocation: $MINIKUBE_CPUS CPUs (OK)"
  fi
  
  if [[ -z "$MINIKUBE_MEMORY" || "$MINIKUBE_MEMORY" -lt 4096 ]]; then
    echo "Warning: Minikube memory allocation may be insufficient. Recommended: 4096MB"
  else
    echo "Minikube memory allocation: $MINIKUBE_MEMORY MB (OK)"
  fi
fi

# Enable required addons
echo "Enabling required Minikube addons..."
minikube addons enable storage-provisioner
minikube addons enable default-storageclass

# Verify kubectl connectivity
echo "Verifying kubectl connectivity..."
kubectl cluster-info

# Check Kubernetes version
echo "Checking Kubernetes version..."
kubectl version --short

# Check available resources
echo "Checking available resources in the cluster..."
kubectl get nodes -o=custom-columns=NAME:.metadata.name,CPU:.status.capacity.cpu,MEMORY:.status.capacity.memory

echo "Minikube environment verification completed."
