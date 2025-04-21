package store

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	
	"github.com/hustler/trading-bot/pkg/execution"
)

// Logger handles database operations and logging
type Logger struct {
	db *sql.DB
}

// NewLogger creates a new Logger
func NewLogger(host string, port int, dbname, user, password string) (*Logger, error) {
	connStr := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		host, port, dbname, user, password)
	
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	
	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	
	return &Logger{db: db}, nil
}

// InitDB initializes the database schema
func (l *Logger) InitDB() error {
	// Create trades table
	_, err := l.db.Exec(`
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
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create trades table: %w", err)
	}
	
	// Create trade_logs table
	_, err = l.db.Exec(`
		CREATE TABLE IF NOT EXISTS trade_logs (
			id SERIAL PRIMARY KEY,
			trade_id VARCHAR(255) REFERENCES trades(id),
			event_type VARCHAR(50) NOT NULL,
			event_data JSONB,
			created_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create trade_logs table: %w", err)
	}
	
	// Create indicators table
	_, err = l.db.Exec(`
		CREATE TABLE IF NOT EXISTS indicators (
			id SERIAL PRIMARY KEY,
			symbol VARCHAR(50) NOT NULL,
			indicator_name VARCHAR(50) NOT NULL,
			value DECIMAL(10, 4) NOT NULL,
			timestamp TIMESTAMP NOT NULL,
			UNIQUE(symbol, indicator_name, timestamp)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create indicators table: %w", err)
	}
	
	// Create app_state table
	_, err = l.db.Exec(`
		CREATE TABLE IF NOT EXISTS app_state (
			key VARCHAR(255) PRIMARY KEY,
			value JSONB NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create app_state table: %w", err)
	}
	
	return nil
}

// LogTrade logs a trade to the database
func (l *Logger) LogTrade(trade *execution.Trade) error {
	// Begin transaction
	tx, err := l.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	
	// Insert into trades table
	_, err = tx.Exec(`
		INSERT INTO trades (id, symbol, quantity, price, type, status, created_at, updated_at, reason)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO UPDATE SET
			status = EXCLUDED.status,
			updated_at = EXCLUDED.updated_at
	`, trade.ID, trade.Symbol, trade.Quantity, trade.Price, trade.Type, trade.Status,
		trade.CreatedAt, trade.UpdatedAt, trade.Reason)
	if err != nil {
		return fmt.Errorf("failed to insert trade: %w", err)
	}
	
	// Insert into trade_logs table
	_, err = tx.Exec(`
		INSERT INTO trade_logs (trade_id, event_type, event_data, created_at)
		VALUES ($1, $2, $3, $4)
	`, trade.ID, trade.Status, fmt.Sprintf(`{"price": %.2f, "quantity": %d}`, trade.Price, trade.Quantity), time.Now())
	if err != nil {
		return fmt.Errorf("failed to insert trade log: %w", err)
	}
	
	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	return nil
}

// LogIndicator logs an indicator value to the database
func (l *Logger) LogIndicator(symbol, indicatorName string, value float64) error {
	_, err := l.db.Exec(`
		INSERT INTO indicators (symbol, indicator_name, value, timestamp)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (symbol, indicator_name, timestamp) DO UPDATE SET
			value = EXCLUDED.value
	`, symbol, indicatorName, value, time.Now())
	if err != nil {
		return fmt.Errorf("failed to insert indicator: %w", err)
	}
	
	return nil
}

// SaveAppState saves application state to the database
func (l *Logger) SaveAppState(key string, value []byte) error {
	_, err := l.db.Exec(`
		INSERT INTO app_state (key, value, updated_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (key) DO UPDATE SET
			value = EXCLUDED.value,
			updated_at = EXCLUDED.updated_at
	`, key, value, time.Now())
	if err != nil {
		return fmt.Errorf("failed to save app state: %w", err)
	}
	
	return nil
}

// LoadAppState loads application state from the database
func (l *Logger) LoadAppState(key string) ([]byte, error) {
	var value []byte
	err := l.db.QueryRow(`
		SELECT value FROM app_state WHERE key = $1
	`, key).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to load app state: %w", err)
	}
	
	return value, nil
}

// GetTradeHistory gets trade history for a symbol
func (l *Logger) GetTradeHistory(symbol string) ([]*execution.Trade, error) {
	rows, err := l.db.Query(`
		SELECT id, symbol, quantity, price, type, status, created_at, updated_at, reason
		FROM trades
		WHERE symbol = $1
		ORDER BY created_at DESC
	`, symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to query trades: %w", err)
	}
	defer rows.Close()
	
	trades := make([]*execution.Trade, 0)
	for rows.Next() {
		trade := &execution.Trade{}
		err := rows.Scan(
			&trade.ID,
			&trade.Symbol,
			&trade.Quantity,
			&trade.Price,
			&trade.Type,
			&trade.Status,
			&trade.CreatedAt,
			&trade.UpdatedAt,
			&trade.Reason,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan trade: %w", err)
		}
		trades = append(trades, trade)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating trades: %w", err)
	}
	
	return trades, nil
}

// ExportDailyReport exports a daily report
func (l *Logger) ExportDailyReport(date time.Time) (string, error) {
	// Format date for SQL query
	dateStr := date.Format("2006-01-02")
	
	// Query trades for the day
	rows, err := l.db.Query(`
		SELECT id, symbol, quantity, price, type, status, created_at, updated_at, reason
		FROM trades
		WHERE DATE(created_at) = $1
		ORDER BY created_at
	`, dateStr)
	if err != nil {
		return "", fmt.Errorf("failed to query trades: %w", err)
	}
	defer rows.Close()
	
	// Build CSV report
	report := "ID,Symbol,Quantity,Price,Type,Status,CreatedAt,UpdatedAt,Reason\n"
	
	for rows.Next() {
		var id, symbol, typeStr, status, reason string
		var quantity int
		var price float64
		var createdAt, updatedAt time.Time
		
		err := rows.Scan(&id, &symbol, &quantity, &price, &typeStr, &status, &createdAt, &updatedAt, &reason)
		if err != nil {
			return "", fmt.Errorf("failed to scan trade: %w", err)
		}
		
		report += fmt.Sprintf("%s,%s,%d,%.2f,%s,%s,%s,%s,%s\n",
			id, symbol, quantity, price, typeStr, status,
			createdAt.Format(time.RFC3339),
			updatedAt.Format(time.RFC3339),
			reason)
	}
	
	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("error iterating trades: %w", err)
	}
	
	return report, nil
}

// Close closes the database connection
func (l *Logger) Close() error {
	return l.db.Close()
}
