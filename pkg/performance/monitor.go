package performance

import (
	"sync"
	"time"

	"github.com/hustler/trading-bot/pkg/signal"
)

// Metrics represents performance metrics for the trading bot
type Metrics struct {
	SignalsCount      int                `json:"signals_count"`
	SuccessCount      int                `json:"success_count"`
	FailureCount      int                `json:"failure_count"`
	PendingCount      int                `json:"pending_count"`
	SuccessRate       float64            `json:"success_rate"`
	AverageROI        float64            `json:"average_roi"`
	TotalProfit       float64            `json:"total_profit"`
	SymbolPerformance map[string]SymbolMetrics `json:"symbol_performance"`
	DailyPerformance  map[string]DailyMetrics  `json:"daily_performance"`
	LastUpdated       time.Time          `json:"last_updated"`
}

// SymbolMetrics represents performance metrics for a specific symbol
type SymbolMetrics struct {
	Symbol       string  `json:"symbol"`
	SignalsCount int     `json:"signals_count"`
	SuccessCount int     `json:"success_count"`
	FailureCount int     `json:"failure_count"`
	PendingCount int     `json:"pending_count"`
	SuccessRate  float64 `json:"success_rate"`
	AverageROI   float64 `json:"average_roi"`
	TotalProfit  float64 `json:"total_profit"`
}

// DailyMetrics represents performance metrics for a specific day
type DailyMetrics struct {
	Date         string  `json:"date"`
	SignalsCount int     `json:"signals_count"`
	SuccessCount int     `json:"success_count"`
	FailureCount int     `json:"failure_count"`
	PendingCount int     `json:"pending_count"`
	SuccessRate  float64 `json:"success_rate"`
	TotalProfit  float64 `json:"total_profit"`
}

// SignalStatus represents the status of a signal
type SignalStatus string

const (
	// StatusActive indicates the signal is active
	StatusActive SignalStatus = "ACTIVE"
	// StatusSuccess indicates the signal was successful
	StatusSuccess SignalStatus = "SUCCESS"
	// StatusFailure indicates the signal failed
	StatusFailure SignalStatus = "FAILURE"
	// StatusExpired indicates the signal expired
	StatusExpired SignalStatus = "EXPIRED"
)

// SignalResult represents the result of a signal
type SignalResult struct {
	SignalID    string      `json:"signal_id"`
	Symbol      string      `json:"symbol"`
	Type        string      `json:"type"`
	EntryPrice  float64     `json:"entry_price"`
	ExitPrice   float64     `json:"exit_price"`
	TargetPrice float64     `json:"target_price"`
	StopLoss    float64     `json:"stop_loss"`
	ExpectedROI float64     `json:"expected_roi"`
	ActualROI   float64     `json:"actual_roi"`
	Status      SignalStatus `json:"status"`
	GeneratedAt time.Time   `json:"generated_at"`
	CompletedAt time.Time   `json:"completed_at"`
}

// Monitor tracks and analyzes trading signal performance
type Monitor struct {
	signals      []*signal.Signal
	results      []*SignalResult
	metrics      *Metrics
	mu           sync.RWMutex
}

// NewMonitor creates a new performance monitor
func NewMonitor() *Monitor {
	return &Monitor{
		signals:      []*signal.Signal{},
		results:      []*SignalResult{},
		metrics:      &Metrics{
			SymbolPerformance: make(map[string]SymbolMetrics),
			DailyPerformance:  make(map[string]DailyMetrics),
			LastUpdated:       time.Now(),
		},
		mu:           sync.RWMutex{},
	}
}

// AddSignal adds a new signal to the monitor
func (m *Monitor) AddSignal(s *signal.Signal) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Add signal to list
	m.signals = append(m.signals, s)
	
	// Create result with active status
	result := &SignalResult{
		SignalID:    s.ID,
		Symbol:      s.Symbol,
		Type:        string(s.Type),
		EntryPrice:  s.Price,
		TargetPrice: s.TargetPrice,
		StopLoss:    s.StopLoss,
		ExpectedROI: s.ExpectedROI,
		Status:      StatusActive,
		GeneratedAt: s.GeneratedAt,
	}
	
	m.results = append(m.results, result)
	
	// Update metrics
	m.updateMetrics()
}

// UpdateSignalStatus updates the status of a signal
func (m *Monitor) UpdateSignalStatus(signalID string, status SignalStatus, exitPrice float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Find result
	var result *SignalResult
	for _, r := range m.results {
		if r.SignalID == signalID {
			result = r
			break
		}
	}
	
	if result == nil {
		return
	}
	
	// Update result
	result.Status = status
	result.ExitPrice = exitPrice
	result.CompletedAt = time.Now()
	
	// Calculate actual ROI
	if result.Type == "BUY" {
		result.ActualROI = (exitPrice - result.EntryPrice) / result.EntryPrice * 100
	} else {
		result.ActualROI = (result.EntryPrice - exitPrice) / result.EntryPrice * 100
	}
	
	// Update metrics
	m.updateMetrics()
}

// GetMetrics returns the current performance metrics
func (m *Monitor) GetMetrics() *Metrics {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	// Return a copy to avoid race conditions
	metricsCopy := *m.metrics
	
	// Copy maps
	metricsCopy.SymbolPerformance = make(map[string]SymbolMetrics, len(m.metrics.SymbolPerformance))
	for k, v := range m.metrics.SymbolPerformance {
		metricsCopy.SymbolPerformance[k] = v
	}
	
	metricsCopy.DailyPerformance = make(map[string]DailyMetrics, len(m.metrics.DailyPerformance))
	for k, v := range m.metrics.DailyPerformance {
		metricsCopy.DailyPerformance[k] = v
	}
	
	return &metricsCopy
}

// GetResults returns all signal results
func (m *Monitor) GetResults() []*SignalResult {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	// Return a copy to avoid race conditions
	resultsCopy := make([]*SignalResult, len(m.results))
	for i, r := range m.results {
		resultCopy := *r
		resultsCopy[i] = &resultCopy
	}
	
	return resultsCopy
}

// GetResultsBySymbol returns signal results for a specific symbol
func (m *Monitor) GetResultsBySymbol(symbol string) []*SignalResult {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	var results []*SignalResult
	
	for _, r := range m.results {
		if r.Symbol == symbol {
			resultCopy := *r
			results = append(results, &resultCopy)
		}
	}
	
	return results
}

// GetResultsByDate returns signal results for a specific date
func (m *Monitor) GetResultsByDate(date string) []*SignalResult {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	var results []*SignalResult
	
	for _, r := range m.results {
		if r.GeneratedAt.Format("2006-01-02") == date {
			resultCopy := *r
			results = append(results, &resultCopy)
		}
	}
	
	return results
}

// updateMetrics recalculates performance metrics
func (m *Monitor) updateMetrics() {
	// Reset counts
	m.metrics.SignalsCount = len(m.results)
	m.metrics.SuccessCount = 0
	m.metrics.FailureCount = 0
	m.metrics.PendingCount = 0
	m.metrics.TotalProfit = 0
	
	// Reset symbol performance
	symbolPerformance := make(map[string]SymbolMetrics)
	
	// Reset daily performance
	dailyPerformance := make(map[string]DailyMetrics)
	
	// Calculate metrics
	for _, r := range m.results {
		// Get or create symbol metrics
		symbol := r.Symbol
		metrics, ok := symbolPerformance[symbol]
		if !ok {
			metrics = SymbolMetrics{
				Symbol: symbol,
			}
		}
		
		// Get or create daily metrics
		date := r.GeneratedAt.Format("2006-01-02")
		daily, ok := dailyPerformance[date]
		if !ok {
			daily = DailyMetrics{
				Date: date,
			}
		}
		
		// Update counts
		metrics.SignalsCount++
		daily.SignalsCount++
		
		// Update status counts
		switch r.Status {
		case StatusSuccess:
			m.metrics.SuccessCount++
			metrics.SuccessCount++
			daily.SuccessCount++
			m.metrics.TotalProfit += r.ActualROI
			metrics.TotalProfit += r.ActualROI
			daily.TotalProfit += r.ActualROI
		case StatusFailure:
			m.metrics.FailureCount++
			metrics.FailureCount++
			daily.FailureCount++
			m.metrics.TotalProfit -= r.ActualROI // Negative ROI
			metrics.TotalProfit -= r.ActualROI
			daily.TotalProfit -= r.ActualROI
		case StatusExpired:
			m.metrics.FailureCount++
			metrics.FailureCount++
			daily.FailureCount++
		default: // Active
			m.metrics.PendingCount++
			metrics.PendingCount++
			daily.PendingCount++
		}
		
		// Update symbol metrics
		symbolPerformance[symbol] = metrics
		
		// Update daily metrics
		dailyPerformance[date] = daily
	}
	
	// Calculate success rates and average ROI
	completedCount := m.metrics.SuccessCount + m.metrics.FailureCount
	if completedCount > 0 {
		m.metrics.SuccessRate = float64(m.metrics.SuccessCount) / float64(completedCount) * 100
		m.metrics.AverageROI = m.metrics.TotalProfit / float64(completedCount)
	}
	
	// Calculate symbol success rates and average ROI
	for symbol, metrics := range symbolPerformance {
		completedCount := metrics.SuccessCount + metrics.FailureCount
		if completedCount > 0 {
			metrics.SuccessRate = float64(metrics.SuccessCount) / float64(completedCount) * 100
			metrics.AverageROI = metrics.TotalProfit / float64(completedCount)
		}
		symbolPerformance[symbol] = metrics
	}
	
	// Calculate daily success rates
	for date, metrics := range dailyPerformance {
		completedCount := metrics.SuccessCount + metrics.FailureCount
		if completedCount > 0 {
			metrics.SuccessRate = float64(metrics.SuccessCount) / float64(completedCount) * 100
		}
		dailyPerformance[date] = metrics
	}
	
	// Update metrics
	m.metrics.SymbolPerformance = symbolPerformance
	m.metrics.DailyPerformance = dailyPerformance
	m.metrics.LastUpdated = time.Now()
}
