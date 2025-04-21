# Hustler Trading Bot - User Guide

## Introduction

The Hustler Trading Bot is an intraday trading signal generator that identifies volatility patterns in stock prices to capture short-term profit opportunities. This guide will help you install, configure, and use the system effectively.

## Installation

### Prerequisites

- Kubernetes cluster (Minikube for local development)
- Go 1.18 or higher
- Docker
- kubectl command-line tool

### Installation Steps

1. **Clone the repository**

```bash
git clone https://github.com/hustler/trading-bot.git
cd trading-bot
```

2. **Build the application**

```bash
go mod tidy
go build -o hustler ./cmd/hustler
```

3. **Deploy to Kubernetes**

```bash
# Start Minikube (if using local development)
minikube start --cpus 4 --memory 4096

# Apply Kubernetes manifests
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secrets.yaml
kubectl apply -f k8s/postgres.yaml
kubectl apply -f k8s/backend.yaml
kubectl apply -f k8s/frontend.yaml
kubectl apply -f k8s/ingress.yaml

# Verify deployment
kubectl get pods
```

## Configuration

### Admin Configuration

Access the admin dashboard at `http://<your-cluster-ip>/admin` to configure:

1. **Stock Watchlist**
   - Add/remove stocks to monitor
   - Set maximum number of stocks to track

2. **Signal Parameters**
   - Volatility threshold (%)
   - Minimum expected ROI (%)
   - Stop-loss percentage
   - Confidence threshold

3. **Trading Hours**
   - Set market hours (e.g., 9:30 AM - 3:30 PM EST)
   - Configure days of operation

4. **LLM Settings**
   - Choose provider (OpenAI or DeepSeek)
   - Set API keys
   - Configure prompt templates

5. **Telegram Settings**
   - Set bot token
   - Configure channel/group ID

### Configuration File

Alternatively, you can edit the configuration file directly:

```bash
# Create default config
./hustler -generate-config > config.json

# Edit the configuration
nano config.json

# Run with custom config
./hustler -config config.json
```

## Using the Telegram Bot

### User Commands

- `/start` - Subscribe to trading signals
- `/help` - Display help information
- `/stop` - Unsubscribe from signals

### Signal Format

When a trading signal is generated, you'll receive a message like this:

```
🚨 BUY SIGNAL: AAPL 🚨

💰 Entry Price: $150.25
🎯 Target Price: $155.50
🛑 Stop Loss: $148.00
📈 Expected ROI: +3.49%
🔍 Confidence: 85%
⏱ Time Frame: 1-3 hours

📝 Rationale:
This BUY signal for AAPL is based on a clear volatility pattern indicating potential upward movement in the short term. The price has shown increased volatility with a bullish bias, technical indicators suggest the stock is currently undervalued, and volume has increased significantly, confirming buying interest.

⏰ Generated at: 2025-04-20 10:15:00
```

## Monitoring Performance

### Performance Dashboard

Access the performance dashboard at `http://<your-cluster-ip>/admin/performance` to view:

- Overall success rate
- Average ROI per signal
- Performance by stock symbol
- Daily performance metrics
- Historical signals and outcomes

### Performance Reports

Generate performance reports using:

```bash
./hustler -report daily > daily_report.txt
./hustler -report weekly > weekly_report.txt
./hustler -report monthly > monthly_report.txt
```

## Troubleshooting

### Common Issues

1. **No signals being generated**
   - Check if market is open
   - Verify stock watchlist is not empty
   - Ensure volatility thresholds are not too high

2. **Telegram bot not sending messages**
   - Verify bot token is correct
   - Check if bot has permission to post in the channel
   - Ensure users have started the bot with /start

3. **LLM integration not working**
   - Verify API keys are correct
   - Check network connectivity to LLM provider
   - Try switching to the alternative LLM provider

### Logs

View application logs:

```bash
kubectl logs -f deployment/hustler-backend
```

## Best Practices

1. **Signal Configuration**
   - Start with conservative volatility thresholds (15-20%)
   - Set reasonable ROI expectations (2-5% for intraday)
   - Use tight stop-losses (1-2% below entry)

2. **Stock Selection**
   - Focus on liquid stocks with high trading volume
   - Select stocks from different sectors for diversification
   - Avoid stocks with upcoming earnings or major announcements

3. **Trading Hours**
   - Avoid the first 15 minutes after market open
   - Be cautious during lunch hours (low volume)
   - Consider ending trading 30 minutes before market close

4. **Risk Management**
   - Never risk more than 1-2% of your capital on a single trade
   - Don't chase signals that have already moved significantly
   - Always use the provided stop-loss levels

## Support and Feedback

For support or to provide feedback:

- Create an issue on GitHub
- Contact the development team at support@hustler-trading-bot.com
- Join our Telegram community group: @HustlerTradingCommunity

## Disclaimer

The Hustler Trading Bot provides trading signals based on technical analysis and volatility patterns. These signals are not financial advice. Always conduct your own research and consider your risk tolerance before making trading decisions. Past performance is not indicative of future results.
