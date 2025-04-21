# Hustler Trading Bot Deployment Guide

This document provides detailed instructions for deploying the Hustler Trading Bot application to Kubernetes using Minikube.

## Prerequisites

- Minikube v1.35.0 or later
- kubectl command-line tool
- Docker
- 4 CPU cores and 4GB RAM available for Minikube

## Deployment Steps

### 1. Start Minikube

If Minikube is not already running, start it with the required resources:

```bash
minikube start --cpus 4 --memory 4096
```

### 2. Enable Ingress Controller

Enable the Nginx Ingress controller in Minikube:

```bash
minikube addons enable ingress
```

### 3. Configure API Credentials

Update the Kubernetes secrets file with your API credentials:

```bash
# Edit the secrets.yaml file
vi k8s/secrets.yaml

# Update the following values:
# - questrade-client-id
# - questrade-refresh-token
# - llm-api-key (OpenAI or Anthropic)
```

### 4. Deploy the Application

Use the provided deployment script to build and deploy the application:

```bash
# Make the script executable
chmod +x scripts/deploy-minikube.sh

# Run the deployment script
./scripts/deploy-minikube.sh
```

The script will:
1. Build the backend Docker image
2. Build the frontend Docker image
3. Apply all Kubernetes manifests
4. Wait for deployments to be ready
5. Display the application URL

### 5. Access the Application

After deployment, you can access the application using the Minikube IP:

```bash
# Get the Minikube IP
minikube ip
```

Add an entry to your `/etc/hosts` file:

```
<minikube-ip> hustler.local
```

Then access the application at: http://hustler.local

## Kubernetes Resources

The application consists of the following Kubernetes resources:

- **ConfigMap**: Application configuration
- **Secrets**: API credentials and database passwords
- **Deployments**:
  - Frontend (React application)
  - Backend (Go API server)
  - PostgreSQL database
- **Services**: For internal communication
- **PersistentVolumeClaims**: For database storage
- **Ingress**: For routing external traffic

## Resource Allocation

The application is configured with the following resource limits:

- **Frontend**: 1 CPU, 1GB RAM
- **Backend**: 4 CPU, 4GB RAM
- **PostgreSQL**: 1 CPU, 1GB RAM

## Monitoring and Maintenance

### Checking Deployment Status

```bash
# Check all resources
kubectl get all

# Check pods status
kubectl get pods

# Check logs for a specific pod
kubectl logs <pod-name>
```

### Updating the Application

To update the application after making changes:

1. Rebuild the Docker images
2. Apply the Kubernetes manifests again

```bash
# Set Docker environment to use Minikube's Docker daemon
eval $(minikube docker-env)

# Rebuild images
docker build -t hustler-trading-bot:latest .
docker build -t hustler-frontend:latest ./web/frontend

# Apply manifests
kubectl apply -f k8s/
```

### Scaling the Application

To scale the application horizontally:

```bash
# Scale the backend
kubectl scale deployment hustler-backend --replicas=3

# Scale the frontend
kubectl scale deployment hustler-frontend --replicas=3
```

## Troubleshooting

### Common Issues

1. **Pods not starting**: Check pod events and logs
   ```bash
   kubectl describe pod <pod-name>
   kubectl logs <pod-name>
   ```

2. **Database connection issues**: Verify PostgreSQL pod is running
   ```bash
   kubectl get pods | grep postgres
   ```

3. **Ingress not working**: Check Ingress controller status
   ```bash
   kubectl get pods -n ingress-nginx
   ```

### Restarting Components

```bash
# Restart a deployment
kubectl rollout restart deployment/<deployment-name>
```

## Cleanup

To remove the application from Minikube:

```bash
kubectl delete -f k8s/
```

To stop Minikube:

```bash
minikube stop
```

To delete the Minikube cluster:

```bash
minikube delete
```
