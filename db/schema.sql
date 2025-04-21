-- Database schema for Hustler Trading Bot

-- Create trades table
CREATE TABLE IF NOT EXISTS trades (
    id VARCHAR(255) PRIMARY KEY,
    symbol VARCHAR(50) NOT NULL,
    quantity INT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    type VARCHAR(10) NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    reason TEXT
);

-- Create trade_logs table
CREATE TABLE IF NOT EXISTS trade_logs (
    id SERIAL PRIMARY KEY,
    trade_id VARCHAR(255) REFERENCES trades(id),
    event_type VARCHAR(50) NOT NULL,
    event_data JSONB,
    created_at TIMESTAMP NOT NULL
);

-- Create indicators table
CREATE TABLE IF NOT EXISTS indicators (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(50) NOT NULL,
    indicator_name VARCHAR(50) NOT NULL,
    value DECIMAL(10, 4) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    UNIQUE(symbol, indicator_name, timestamp)
);

-- Create app_state table
CREATE TABLE IF NOT EXISTS app_state (
    key VARCHAR(255) PRIMARY KEY,
    value JSONB NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Create daily_summary table
CREATE TABLE IF NOT EXISTS daily_summary (
    date DATE PRIMARY KEY,
    total_trades INT NOT NULL DEFAULT 0,
    profitable_trades INT NOT NULL DEFAULT 0,
    total_profit DECIMAL(10, 2) NOT NULL DEFAULT 0,
    total_loss DECIMAL(10, 2) NOT NULL DEFAULT 0,
    net_pnl DECIMAL(10, 2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Create stock_watchlist table
CREATE TABLE IF NOT EXISTS stock_watchlist (
    symbol VARCHAR(50) PRIMARY KEY,
    added_at TIMESTAMP NOT NULL,
    last_traded_at TIMESTAMP,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_trades_symbol ON trades(symbol);
CREATE INDEX IF NOT EXISTS idx_trades_created_at ON trades(created_at);
CREATE INDEX IF NOT EXISTS idx_trade_logs_trade_id ON trade_logs(trade_id);
CREATE INDEX IF NOT EXISTS idx_indicators_symbol ON indicators(symbol);
CREATE INDEX IF NOT EXISTS idx_indicators_timestamp ON indicators(timestamp);
