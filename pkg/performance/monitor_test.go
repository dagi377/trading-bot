package performance

import (
	"testing"
	"time"

	"github.com/hustler/trading-bot/pkg/signal"
	"github.com/stretchr/testify/assert"
)

func TestNewMonitor(t *testing.T) {
	monitor := NewMonitor()
	
	assert.NotNil(t, monitor)
	assert.NotNil(t, monitor.signals)
	assert.NotNil(t, monitor.results)
	assert.NotNil(t, monitor.metrics)
	assert.NotNil(t, monitor.metrics.SymbolPerformance)
	assert.NotNil(t, monitor.metrics.DailyPerformance)
}

func TestAddSignal(t *testing.T) {
	monitor := NewMonitor()
	
	// Create test signal
	testSignal := createTestSignal("AAPL", signal.BUY, 150.0, 155.0, 148.0)
	
	// Add signal
	monitor.AddSignal(testSignal)
	
	// Verify signal was added
	assert.Len(t, monitor.signals, 1)
	assert.Len(t, monitor.results, 1)
	
	// Verify result
	result := monitor.results[0]
	assert.Equal(t, testSignal.ID, result.SignalID)
	assert.Equal(t, testSignal.Symbol, result.Symbol)
	assert.Equal(t, string(testSignal.Type), result.Type)
	assert.Equal(t, testSignal.Price, result.EntryPrice)
	assert.Equal(t, testSignal.TargetPrice, result.TargetPrice)
	assert.Equal(t, testSignal.StopLoss, result.StopLoss)
	assert.Equal(t, testSignal.ExpectedROI, result.ExpectedROI)
	assert.Equal(t, StatusActive, result.Status)
	
	// Verify metrics
	metrics := monitor.GetMetrics()
	assert.Equal(t, 1, metrics.SignalsCount)
	assert.Equal(t, 0, metrics.SuccessCount)
	assert.Equal(t, 0, metrics.FailureCount)
	assert.Equal(t, 1, metrics.PendingCount)
	
	// Verify symbol metrics
	symbolMetrics, ok := metrics.SymbolPerformance["AAPL"]
	assert.True(t, ok)
	assert.Equal(t, "AAPL", symbolMetrics.Symbol)
	assert.Equal(t, 1, symbolMetrics.SignalsCount)
	assert.Equal(t, 0, symbolMetrics.SuccessCount)
	assert.Equal(t, 0, symbolMetrics.FailureCount)
	assert.Equal(t, 1, symbolMetrics.PendingCount)
	
	// Verify daily metrics
	date := testSignal.GeneratedAt.Format("2006-01-02")
	dailyMetrics, ok := metrics.DailyPerformance[date]
	assert.True(t, ok)
	assert.Equal(t, date, dailyMetrics.Date)
	assert.Equal(t, 1, dailyMetrics.SignalsCount)
	assert.Equal(t, 0, dailyMetrics.SuccessCount)
	assert.Equal(t, 0, dailyMetrics.FailureCount)
	assert.Equal(t, 1, dailyMetrics.PendingCount)
}

func TestUpdateSignalStatus(t *testing.T) {
	monitor := NewMonitor()
	
	// Create and add test signals
	signal1 := createTestSignal("AAPL", signal.BUY, 150.0, 155.0, 148.0)
	signal2 := createTestSignal("MSFT", signal.SELL, 350.0, 345.0, 352.0)
	
	monitor.AddSignal(signal1)
	monitor.AddSignal(signal2)
	
	// Update signal status - success
	monitor.UpdateSignalStatus(signal1.ID, StatusSuccess, 155.0)
	
	// Update signal status - failure
	monitor.UpdateSignalStatus(signal2.ID, StatusFailure, 355.0)
	
	// Verify results
	results := monitor.GetResults()
	assert.Len(t, results, 2)
	
	// Find results by ID
	var result1, result2 *SignalResult
	for _, r := range results {
		if r.SignalID == signal1.ID {
			result1 = r
		} else if r.SignalID == signal2.ID {
			result2 = r
		}
	}
	
	// Verify result 1
	assert.NotNil(t, result1)
	assert.Equal(t, StatusSuccess, result1.Status)
	assert.Equal(t, 155.0, result1.ExitPrice)
	assert.InDelta(t, 3.33, result1.ActualROI, 0.01) // (155-150)/150 * 100 = 3.33%
	
	// Verify result 2
	assert.NotNil(t, result2)
	assert.Equal(t, StatusFailure, result2.Status)
	assert.Equal(t, 355.0, result2.ExitPrice)
	assert.InDelta(t, -1.43, result2.ActualROI, 0.01) // (350-355)/350 * 100 = -1.43%
	
	// Verify metrics
	metrics := monitor.GetMetrics()
	assert.Equal(t, 2, metrics.SignalsCount)
	assert.Equal(t, 1, metrics.SuccessCount)
	assert.Equal(t, 1, metrics.FailureCount)
	assert.Equal(t, 0, metrics.PendingCount)
	assert.InDelta(t, 50.0, metrics.SuccessRate, 0.01) // 1/2 * 100 = 50%
	assert.InDelta(t, 0.95, metrics.AverageROI, 0.01)  // (3.33 - 1.43) / 2 = 0.95
	assert.InDelta(t, 1.9, metrics.TotalProfit, 0.01)  // 3.33 - 1.43 = 1.9
	
	// Verify symbol metrics
	appleMetrics := metrics.SymbolPerformance["AAPL"]
	assert.Equal(t, 1, appleMetrics.SignalsCount)
	assert.Equal(t, 1, appleMetrics.SuccessCount)
	assert.Equal(t, 0, appleMetrics.FailureCount)
	assert.Equal(t, 100.0, appleMetrics.SuccessRate)
	assert.InDelta(t, 3.33, appleMetrics.AverageROI, 0.01)
	
	msftMetrics := metrics.SymbolPerformance["MSFT"]
	assert.Equal(t, 1, msftMetrics.SignalsCount)
	assert.Equal(t, 0, msftMetrics.SuccessCount)
	assert.Equal(t, 1, msftMetrics.FailureCount)
	assert.Equal(t, 0.0, msftMetrics.SuccessRate)
	assert.InDelta(t, -1.43, msftMetrics.AverageROI, 0.01)
}

func TestGetResultsBySymbol(t *testing.T) {
	monitor := NewMonitor()
	
	// Create and add test signals
	signal1 := createTestSignal("AAPL", signal.BUY, 150.0, 155.0, 148.0)
	signal2 := createTestSignal("AAPL", signal.SELL, 160.0, 155.0, 162.0)
	signal3 := createTestSignal("MSFT", signal.BUY, 350.0, 355.0, 348.0)
	
	monitor.AddSignal(signal1)
	monitor.AddSignal(signal2)
	monitor.AddSignal(signal3)
	
	// Get results by symbol
	appleResults := monitor.GetResultsBySymbol("AAPL")
	msftResults := monitor.GetResultsBySymbol("MSFT")
	
	// Verify results
	assert.Len(t, appleResults, 2)
	assert.Len(t, msftResults, 1)
	
	// Verify Apple results
	assert.Equal(t, "AAPL", appleResults[0].Symbol)
	assert.Equal(t, "AAPL", appleResults[1].Symbol)
	
	// Verify MSFT results
	assert.Equal(t, "MSFT", msftResults[0].Symbol)
}

func TestGetResultsByDate(t *testing.T) {
	monitor := NewMonitor()
	
	// Create signals with different dates
	yesterday := time.Now().Add(-24 * time.Hour)
	
	signal1 := createTestSignal("AAPL", signal.BUY, 150.0, 155.0, 148.0)
	signal1.GeneratedAt = yesterday
	
	signal2 := createTestSignal("MSFT", signal.BUY, 350.0, 355.0, 348.0)
	
	monitor.AddSignal(signal1)
	monitor.AddSignal(signal2)
	
	// Get results by date
	yesterdayResults := monitor.GetResultsByDate(yesterday.Format("2006-01-02"))
	todayResults := monitor.GetResultsByDate(time.Now().Format("2006-01-02"))
	
	// Verify results
	assert.Len(t, yesterdayResults, 1)
	assert.Len(t, todayResults, 1)
	
	// Verify yesterday results
	assert.Equal(t, "AAPL", yesterdayResults[0].Symbol)
	
	// Verify today results
	assert.Equal(t, "MSFT", todayResults[0].Symbol)
}

func TestUpdateMetrics(t *testing.T) {
	monitor := NewMonitor()
	
	// Create and add test signals
	signal1 := createTestSignal("AAPL", signal.BUY, 150.0, 155.0, 148.0)
	signal2 := createTestSignal("MSFT", signal.SELL, 350.0, 345.0, 352.0)
	signal3 := createTestSignal("GOOGL", signal.BUY, 180.0, 185.0, 178.0)
	
	monitor.AddSignal(signal1)
	monitor.AddSignal(signal2)
	monitor.AddSignal(signal3)
	
	// Update signal statuses
	monitor.UpdateSignalStatus(signal1.ID, StatusSuccess, 155.0)
	monitor.UpdateSignalStatus(signal2.ID, StatusFailure, 355.0)
	monitor.UpdateSignalStatus(signal3.ID, StatusExpired, 179.0)
	
	// Verify metrics
	metrics := monitor.GetMetrics()
	assert.Equal(t, 3, metrics.SignalsCount)
	assert.Equal(t, 1, metrics.SuccessCount)
	assert.Equal(t, 2, metrics.FailureCount)
	assert.Equal(t, 0, metrics.PendingCount)
	assert.InDelta(t, 33.33, metrics.SuccessRate, 0.01) // 1/3 * 100 = 33.33%
	
	// Verify symbol metrics
	assert.Len(t, metrics.SymbolPerformance, 3)
	
	// Verify daily metrics
	date := time.Now().Format("2006-01-02")
	dailyMetrics := metrics.DailyPerformance[date]
	assert.Equal(t, 3, dailyMetrics.SignalsCount)
	assert.Equal(t, 1, dailyMetrics.SuccessCount)
	assert.Equal(t, 2, dailyMetrics.FailureCount)
	assert.Equal(t, 0, dailyMetrics.PendingCount)
	assert.InDelta(t, 33.33, dailyMetrics.SuccessRate, 0.01)
}

// Helper function to create test signals
func createTestSignal(symbol string, signalType signal.SignalType, price, targetPrice, stopLoss float64) *signal.Signal {
	return &signal.Signal{
		ID:            fmt.Sprintf("SIG-%s-%s-%d", symbol, signalType, time.Now().Unix()),
		Symbol:        symbol,
		Type:          signalType,
		Price:         price,
		TargetPrice:   targetPrice,
		StopLoss:      stopLoss,
		ExpectedROI:   (targetPrice - price) / price * 100,
		Confidence:    0.85,
		GeneratedAt:   time.Now(),
		TimeFrame:     "1-3 hours",
		TechnicalData: map[string]float64{"RSI": 35.0, "Volume": 1500000},
		Status:        "ACTIVE",
	}
}
