package signal

import (
	"testing"
	"time"

	"github.com/hustler/trading-bot/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestNewGenerator(t *testing.T) {
	cfg := config.CreateDefaultConfig()
	generator := NewGenerator(cfg)
	
	assert.NotNil(t, generator)
	assert.Equal(t, cfg, generator.config)
}

func TestGenerateSignals(t *testing.T) {
	// Create test configuration
	cfg := config.CreateDefaultConfig()
	cfg.VolatilityParams.MinVolatilityPercent = 1.0
	cfg.VolatilityParams.MinExpectedROI = 1.5
	cfg.VolatilityParams.ConfidenceThreshold = 0.6
	
	// Create generator
	generator := NewGenerator(cfg)
	
	// Create test market data
	marketData := map[string]MarketData{
		"AAPL": createTestMarketData("AAPL", true),  // Bullish pattern
		"MSFT": createTestMarketData("MSFT", false), // Bearish pattern
		"GOOGL": {
			Symbol:     "GOOGL",
			Prices:     []float64{150.0, 151.0}, // Not enough data points
			Volumes:    []float64{1000000, 1100000},
			Timestamps: []time.Time{time.Now().Add(-1 * time.Hour), time.Now()},
		},
	}
	
	// Generate signals
	signals, err := generator.GenerateSignals(marketData)
	
	// Verify results
	assert.NoError(t, err)
	assert.Len(t, signals, 2) // Should generate signals for AAPL and MSFT, but not GOOGL
	
	// Find signals by symbol
	var appleSignal, msftSignal *Signal
	for _, s := range signals {
		if s.Symbol == "AAPL" {
			appleSignal = s
		} else if s.Symbol == "MSFT" {
			msftSignal = s
		}
	}
	
	// Verify AAPL signal
	assert.NotNil(t, appleSignal)
	assert.Equal(t, BUY, appleSignal.Type)
	assert.Greater(t, appleSignal.TargetPrice, appleSignal.Price)
	assert.Less(t, appleSignal.StopLoss, appleSignal.Price)
	assert.GreaterOrEqual(t, appleSignal.ExpectedROI, cfg.VolatilityParams.MinExpectedROI)
	assert.GreaterOrEqual(t, appleSignal.Confidence, cfg.VolatilityParams.ConfidenceThreshold)
	
	// Verify MSFT signal
	assert.NotNil(t, msftSignal)
	assert.Equal(t, SELL, msftSignal.Type)
	assert.Less(t, msftSignal.TargetPrice, msftSignal.Price)
	assert.Greater(t, msftSignal.StopLoss, msftSignal.Price)
	assert.GreaterOrEqual(t, msftSignal.ExpectedROI, cfg.VolatilityParams.MinExpectedROI)
	assert.GreaterOrEqual(t, msftSignal.Confidence, cfg.VolatilityParams.ConfidenceThreshold)
}

func TestCalculateTechnicalIndicators(t *testing.T) {
	// Create test market data
	data := createTestMarketData("TEST", true)
	
	// Create test parameters
	params := config.VolatilityConfig{
		BollingerPeriod:    20,
		BollingerDeviation: 2.0,
		RSIPeriod:          14,
	}
	
	// Calculate indicators
	indicators := calculateTechnicalIndicators(data, params)
	
	// Verify indicators were calculated
	assert.Contains(t, indicators, "sma")
	assert.Contains(t, indicators, "upper_band")
	assert.Contains(t, indicators, "lower_band")
	assert.Contains(t, indicators, "rsi")
	assert.Contains(t, indicators, "volume_ratio")
	assert.Contains(t, indicators, "price_change")
	
	// Verify SMA is reasonable
	assert.InDelta(t, 150.0, indicators["sma"], 10.0)
	
	// Verify Bollinger Bands
	assert.Greater(t, indicators["upper_band"], indicators["sma"])
	assert.Less(t, indicators["lower_band"], indicators["sma"])
	
	// Verify RSI is between 0 and 100
	assert.GreaterOrEqual(t, indicators["rsi"], 0.0)
	assert.LessOrEqual(t, indicators["rsi"], 100.0)
}

func TestCalculateSMA(t *testing.T) {
	// Test with valid data
	values := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	sma := calculateSMA(values, 3)
	assert.Equal(t, 4.0, sma) // (3+4+5)/3 = 4
	
	// Test with period larger than data
	sma = calculateSMA(values, 10)
	assert.Equal(t, 0.0, sma)
}

func TestCalculateRSI(t *testing.T) {
	// Test with all gains
	prices := []float64{100.0, 101.0, 102.0, 103.0, 104.0}
	rsi := calculateRSI(prices, 3)
	assert.Equal(t, 100.0, rsi)
	
	// Test with all losses
	prices = []float64{100.0, 99.0, 98.0, 97.0, 96.0}
	rsi = calculateRSI(prices, 3)
	assert.InDelta(t, 0.0, rsi, 0.1)
	
	// Test with mixed gains and losses
	prices = []float64{100.0, 101.0, 99.0, 102.0, 98.0}
	rsi = calculateRSI(prices, 3)
	assert.Greater(t, rsi, 0.0)
	assert.Less(t, rsi, 100.0)
	
	// Test with insufficient data
	prices = []float64{100.0}
	rsi = calculateRSI(prices, 3)
	assert.Equal(t, 50.0, rsi) // Default to neutral
}

func TestDetermineSignalType(t *testing.T) {
	// Test BUY signal - oversold
	indicators := map[string]float64{
		"price":      95.0,
		"upper_band": 110.0,
		"lower_band": 100.0,
		"rsi":        25.0,
		"price_change": 0.5,
	}
	signalType := determineSignalType(indicators)
	assert.Equal(t, BUY, signalType)
	
	// Test SELL signal - overbought
	indicators = map[string]float64{
		"price":      115.0,
		"upper_band": 110.0,
		"lower_band": 100.0,
		"rsi":        75.0,
		"price_change": -0.5,
	}
	signalType = determineSignalType(indicators)
	assert.Equal(t, SELL, signalType)
	
	// Test HOLD signal - neutral
	indicators = map[string]float64{
		"price":      105.0,
		"upper_band": 110.0,
		"lower_band": 100.0,
		"rsi":        50.0,
		"price_change": 0.1,
	}
	signalType = determineSignalType(indicators)
	assert.Equal(t, HOLD, signalType)
}

func TestCalculatePriceLevels(t *testing.T) {
	// Test BUY signal
	indicators := map[string]float64{
		"sma":        100.0,
		"upper_band": 110.0,
		"lower_band": 90.0,
	}
	params := config.VolatilityConfig{
		MinExpectedROI:  2.0,
		StopLossPercent: 1.0,
	}
	currentPrice := 100.0
	targetPrice, stopLoss := calculatePriceLevels(currentPrice, BUY, indicators, params)
	
	assert.InDelta(t, 102.0, targetPrice, 0.1) // Min of upper band and 2% gain
	assert.InDelta(t, 99.0, stopLoss, 0.1)     // Max of lower band and 1% loss
	
	// Test SELL signal
	targetPrice, stopLoss = calculatePriceLevels(currentPrice, SELL, indicators, params)
	
	assert.InDelta(t, 98.0, targetPrice, 0.1)  // Max of lower band and 2% drop
	assert.InDelta(t, 101.0, stopLoss, 0.1)    // Min of upper band and 1% gain
}

func TestCalculateExpectedROI(t *testing.T) {
	// Test BUY signal
	currentPrice := 100.0
	targetPrice := 105.0
	roi := calculateExpectedROI(currentPrice, targetPrice, BUY)
	assert.Equal(t, 5.0, roi) // (105-100)/100 * 100 = 5%
	
	// Test SELL signal
	targetPrice = 95.0
	roi = calculateExpectedROI(currentPrice, targetPrice, SELL)
	assert.Equal(t, 5.0, roi) // (100-95)/100 * 100 = 5%
}

func TestFormatSignalMessage(t *testing.T) {
	// Create test signal
	signal := &Signal{
		Symbol:      "AAPL",
		Type:        BUY,
		Price:       150.25,
		TargetPrice: 155.50,
		StopLoss:    148.00,
		ExpectedROI: 3.5,
		Confidence:  0.85,
		Rationale:   "Strong momentum with increasing volume",
		GeneratedAt: time.Date(2025, 4, 20, 10, 15, 0, 0, time.UTC),
		TimeFrame:   "1-3 hours",
	}
	
	// Format message
	message := FormatSignalMessage(signal)
	
	// Verify message contains key information
	assert.Contains(t, message, "BUY SIGNAL: AAPL")
	assert.Contains(t, message, "Entry Price: $150.25")
	assert.Contains(t, message, "Target Price: $155.50")
	assert.Contains(t, message, "Stop Loss: $148.00")
	assert.Contains(t, message, "Expected ROI: +3.50%")
	assert.Contains(t, message, "Confidence: 85%")
	assert.Contains(t, message, "Strong momentum with increasing volume")
	assert.Contains(t, message, "2025-04-20 10:15:00")
	
	// Test SELL signal formatting
	signal.Type = SELL
	message = FormatSignalMessage(signal)
	assert.Contains(t, message, "SELL SIGNAL: AAPL")
	assert.Contains(t, message, "Expected ROI: -3.50%")
}

// Helper function to create test market data
func createTestMarketData(symbol string, bullish bool) MarketData {
	// Create base prices
	basePrices := make([]float64, 50)
	for i := 0; i < 50; i++ {
		basePrices[i] = 150.0
	}
	
	// Create volumes
	volumes := make([]float64, 50)
	for i := 0; i < 50; i++ {
		volumes[i] = 1000000.0
	}
	
	// Create timestamps
	timestamps := make([]time.Time, 50)
	now := time.Now()
	for i := 0; i < 50; i++ {
		timestamps[49-i] = now.Add(time.Duration(-i) * time.Hour)
	}
	
	// Modify recent prices to create pattern
	if bullish {
		// Create bullish pattern (price increasing with higher volume)
		for i := 40; i < 50; i++ {
			basePrices[i] = basePrices[i-1] * (1 + 0.005*float64(i-39))
			volumes[i] = volumes[i-1] * 1.1
		}
	} else {
		// Create bearish pattern (price decreasing with higher volume)
		for i := 40; i < 50; i++ {
			basePrices[i] = basePrices[i-1] * (1 - 0.005*float64(i-39))
			volumes[i] = volumes[i-1] * 1.1
		}
	}
	
	return MarketData{
		Symbol:     symbol,
		Prices:     basePrices,
		Volumes:    volumes,
		Timestamps: timestamps,
	}
}
