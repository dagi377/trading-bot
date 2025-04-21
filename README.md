# Hustler Trading Bot - Deployment Guide

## Overview

The Hustler Trading Bot is a Go-based intraday stock trading application that uses the Questrade API and LLM-driven analysis to make trading decisions. This guide will walk you through the process of deploying the application to Kubernetes using Minikube on your Mac.

## Prerequisites

- macOS with Minikube installed (v1.35.0+)
- Docker installed and running
- kubectl command-line tool
- At least 4 CPU cores and 4GB RAM available for Minikube

## Configuration

Before deploying the application, you need to configure the following:

1. **API Credentials**: Update the `k8s/secrets.yaml` file with your Questrade API credentials and LLM API key:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: hustler-secrets
type: Opaque
stringData:
  db-password: hustlerpass
  questrade-client-id: "YOUR_QUESTRADE_CLIENT_ID"
  questrade-refresh-token: "YOUR_QUESTRADE_REFRESH_TOKEN"
  llm-api-key: "YOUR_LLM_API_KEY"
```

2. **Trading Parameters**: Update the `k8s/configmap.yaml` file to customize your trading parameters:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: hustler-config
data:
  config.json: |
    {
      "watch_list": ["NVDA", "SHOP", "TD"],
      "max_loss_per_trade": 30.0,
      "max_daily_loss": 400.0,
      "capital_per_stock": 300.0,
      "poll_interval": 5,
      "log_level": "info",
      "ui_port": 8080
    }
```

## Deployment Steps

1. **Clone the Repository**:

```bash
git clone <repository-url>
cd hustler-trading-bot
```

2. **Verify Minikube Environment**:

```bash
chmod +x ./scripts/verify-minikube.sh
./scripts/verify-minikube.sh
```

This script checks if Minikube is running with sufficient resources (4 CPU cores and 4GB RAM).

3. **Deploy the Application**:

```bash
chmod +x ./scripts/deploy.sh
./scripts/deploy.sh
```

This script:
- Builds the Docker image
- Applies Kubernetes manifests
- Deploys PostgreSQL database
- Deploys the Trading Bot application
- Provides the URL to access the UI

4. **Validate the Deployment**:

```bash
chmod +x ./scripts/validate.sh
./scripts/validate.sh
```

This script checks if all components are running correctly.

## Accessing the Application

After successful deployment, you can access the Trading Bot UI using the URL provided by the deployment script. The UI allows you to:

- View the current watchlist
- Add stocks to the watchlist
- View trade history
- Monitor risk metrics

## API Endpoints

The Trading Bot exposes the following API endpoints:

- `GET /api/stocks` - Get all stocks in the watchlist
- `GET /api/stock/{symbol}` - Get details for a specific stock
- `POST /api/stock/add` - Add a stock to the watchlist
- `GET /api/trades` - Get all trades
- `GET /api/risk` - Get risk report
- `GET /api/indicators/{symbol}` - Get indicators for a specific stock

## Troubleshooting

### Common Issues

1. **Minikube Not Starting**:
   - Ensure Docker is running
   - Check if you have sufficient resources (4 CPU cores and 4GB RAM)
   - Try running `minikube delete` and then `minikube start --cpus=4 --memory=4096`

2. **Pods Not Starting**:
   - Check pod status: `kubectl get pods`
   - Check pod logs: `kubectl logs <pod-name>`
   - Check pod events: `kubectl describe pod <pod-name>`

3. **Database Connection Issues**:
   - Ensure PostgreSQL pod is running: `kubectl get pods -l app=hustler-postgres`
   - Check PostgreSQL logs: `kubectl logs -l app=hustler-postgres`
   - Verify database credentials in secrets and config

4. **API Connection Issues**:
   - Verify Questrade API credentials in secrets
   - Check Trading Bot logs: `kubectl logs -l app=hustler-trading-bot`

## Cleanup

To remove the application from your Minikube cluster:

```bash
kubectl delete -f k8s/trading-bot.yaml
kubectl delete -f k8s/postgres.yaml
kubectl delete -f k8s/secrets.yaml
kubectl delete -f k8s/configmap.yaml
```

## Architecture

The Hustler Trading Bot consists of the following components:

1. **Trading Bot Application**: A Go application that handles trading logic, LLM integration, and UI
2. **PostgreSQL Database**: Stores trade history, indicators, and application state

The application is structured with the following modules:

- `auth`: Handles Questrade API authentication
- `data`: Manages real-time market data
- `indicators`: Calculates technical indicators
- `strategy`: Integrates with LLM for trade decisions
- `execution`: Manages trade execution
- `store`: Handles database operations
- `monitor`: Enforces risk management
- `ui`: Provides web interface and API endpoints

## Resource Requirements

The application is configured with the following resource requirements:

- Trading Bot: 1-4 CPU cores, 1-4GB RAM
- PostgreSQL: 0.5-1 CPU core, 512MB-1GB RAM

These can be adjusted in the Kubernetes deployment files based on your needs.

## Maintenance

### Updating the Application

To update the application:

1. Make changes to the code
2. Rebuild the Docker image: `docker build -t hustler-trading-bot:latest .`
3. Restart the deployment: `kubectl rollout restart deployment hustler-trading-bot`

### Backing Up the Database

To back up the PostgreSQL database:

```bash
kubectl exec -it $(kubectl get pod -l app=hustler-postgres -o jsonpath="{.items[0].metadata.name}") -- pg_dump -U hustler hustler > backup.sql
```

### Restoring the Database

To restore the PostgreSQL database:

```bash
cat backup.sql | kubectl exec -i $(kubectl get pod -l app=hustler-postgres -o jsonpath="{.items[0].metadata.name}") -- psql -U hustler -d hustler
```
