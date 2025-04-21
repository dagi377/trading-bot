#!/bin/bash

# Script to validate the Hustler Trading Bot deployment

# Check if all pods are running
echo "Checking pod status..."
kubectl get pods

# Check if all services are running
echo "Checking service status..."
kubectl get services

# Check PostgreSQL deployment
echo "Checking PostgreSQL deployment..."
kubectl describe deployment hustler-postgres

# Check Trading Bot deployment
echo "Checking Trading Bot deployment..."
kubectl describe deployment hustler-trading-bot

# Check logs for PostgreSQL
echo "Checking PostgreSQL logs..."
kubectl logs -l app=hustler-postgres --tail=20

# Check logs for Trading Bot
echo "Checking Trading Bot logs..."
kubectl logs -l app=hustler-trading-bot --tail=20

# Check if UI is accessible
echo "Checking if UI is accessible..."
TRADING_BOT_URL=$(minikube service hustler-trading-bot --url)
curl -s -o /dev/null -w "%{http_code}" $TRADING_BOT_URL

echo "Validation completed. If all checks passed, the deployment is successful."
