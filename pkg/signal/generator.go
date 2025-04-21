package signal

import (
	"fmt"
	"math"
	"time"

	"github.com/hustler/trading-bot/pkg/config"
)

// SignalType represents the type of trading signal
type SignalType string

const (
	// BUY signal
	BUY SignalType = "BUY"
	// SELL signal
	SELL SignalType = "SELL"
	// HOLD signal (no action)
	HOLD SignalType = "HOLD"
)

// Signal represents a trading signal
type Signal struct {
	ID            string             `json:"id"`
	Symbol        string             `json:"symbol"`
	Type          SignalType         `json:"type"`
	Price         float64            `json:"price"`
	TargetPrice   float64            `json:"target_price"`
	StopLoss      float64            `json:"stop_loss"`
	ExpectedROI   float64            `json:"expected_roi"`
	Confidence    float64            `json:"confidence"`
	Rationale     string             `json:"rationale"`
	GeneratedAt   time.Time          `json:"generated_at"`
	TimeFrame     string             `json:"time_frame"`
	TechnicalData map[string]float64 `json:"technical_data"`
	Status        string             `json:"status"`
}

// Generator is responsible for generating trading signals
type Generator struct {
	config *config.Config
}

// NewGenerator creates a new signal generator
func NewGenerator(cfg *config.Config) *Generator {
	return &Generator{
		config: cfg,
	}
}

// GenerateSignals analyzes market data and generates trading signals
func (g *Generator) GenerateSignals(marketData map[string]MarketData) ([]*Signal, error) {
	signals := []*Signal{}

	for symbol, data := range marketData {
		// Skip if not enough data points
		if len(data.Prices) < 30 || len(data.Volumes) < 30 {
			continue
		}

		// Analyze volatility patterns
		signal, generated := g.analyzeVolatilityPatterns(symbol, data)
		if generated {
			signals = append(signals, signal)
		}
	}

	return signals, nil
}

// analyzeVolatilityPatterns analyzes volatility patterns for a stock
func (g *Generator) analyzeVolatilityPatterns(symbol string, data MarketData) (*Signal, bool) {
	// Get current price
	currentPrice := data.Prices[len(data.Prices)-1]
	
	// Calculate technical indicators
	technicalData := calculateTechnicalIndicators(data, g.config.VolatilityParams, currentPrice)
	
	// Calculate volatility score
	volatilityScore := calculateVolatilityScore(technicalData, g.config.VolatilityParams)
	
	// If volatility score is below threshold, no signal
	if volatilityScore < g.config.VolatilityParams.ConfidenceThreshold {
		return nil, false
	}
	
	// Determine signal type based on indicators
	signalType := determineSignalType(technicalData)
	
	// If HOLD, no signal
	if signalType == HOLD {
		return nil, false
	}
	
	// Calculate target price and stop loss
	targetPrice, stopLoss := calculatePriceLevels(currentPrice, signalType, technicalData, g.config.VolatilityParams)
	
	// Calculate expected ROI
	expectedROI := calculateExpectedROI(currentPrice, targetPrice, signalType)
	
	// If expected ROI is below minimum, no signal
	if expectedROI < g.config.VolatilityParams.MinExpectedROI {
		return nil, false
	}
	
	// Create signal
	signal := &Signal{
		ID:            fmt.Sprintf("SIG-%s-%s-%d", symbol, signalType, time.Now().Unix()),
		Symbol:        symbol,
		Type:          signalType,
		Price:         currentPrice,
		TargetPrice:   targetPrice,
		StopLoss:      stopLoss,
		ExpectedROI:   expectedROI,
		Confidence:    volatilityScore,
		GeneratedAt:   time.Now(),
		TimeFrame:     "1-3 hours",
		TechnicalData: technicalData,
		Status:        "ACTIVE",
	}
	
	return signal, true
}

// MarketData represents market data for a stock
type MarketData struct {
	Symbol     string
	Prices     []float64
	Volumes    []float64
	Timestamps []time.Time
}

// calculateTechnicalIndicators calculates technical indicators from market data
func calculateTechnicalIndicators(data MarketData, params config.VolatilityConfig, currentPrice float64) map[string]float64 {
	indicators := make(map[string]float64)
	
	// Store current price in indicators map
	indicators["price"] = currentPrice
	
	// Calculate Bollinger Bands
	sma := calculateSMA(data.Prices, params.BollingerPeriod)
	stdDev := calculateStdDev(data.Prices, params.BollingerPeriod)
	upperBand := sma + params.BollingerDeviation*stdDev
	lowerBand := sma - params.BollingerDeviation*stdDev
	
	// Calculate RSI
	rsi := calculateRSI(data.Prices, params.RSIPeriod)
	
	// Calculate volume ratio
	avgVolume := calculateSMA(data.Volumes, 10)
	currentVolume := data.Volumes[len(data.Volumes)-1]
	volumeRatio := currentVolume / avgVolume * 100
	
	// Calculate price volatility
	priceChange := calculatePriceChange(data.Prices)
	
	// Store indicators
	indicators["sma"] = sma
	indicators["upper_band"] = upperBand
	indicators["lower_band"] = lowerBand
	indicators["rsi"] = rsi
	indicators["volume_ratio"] = volumeRatio
	indicators["price_change"] = priceChange
	
	return indicators
}

// calculateSMA calculates Simple Moving Average
func calculateSMA(values []float64, period int) float64 {
	if len(values) < period {
		return 0
	}
	
	sum := 0.0
	for i := len(values) - period; i < len(values); i++ {
		sum += values[i]
	}
	
	return sum / float64(period)
}

// calculateStdDev calculates Standard Deviation
func calculateStdDev(values []float64, period int) float64 {
	if len(values) < period {
		return 0
	}
	
	sma := calculateSMA(values, period)
	sumSquaredDiff := 0.0
	
	for i := len(values) - period; i < len(values); i++ {
		diff := values[i] - sma
		sumSquaredDiff += diff * diff
	}
	
	return math.Sqrt(sumSquaredDiff / float64(period))
}

// calculateRSI calculates Relative Strength Index
func calculateRSI(prices []float64, period int) float64 {
	if len(prices) < period+1 {
		return 50 // Default to neutral
	}
	
	gains := 0.0
	losses := 0.0
	
	for i := len(prices) - period; i < len(prices); i++ {
		change := prices[i] - prices[i-1]
		if change >= 0 {
			gains += change
		} else {
			losses -= change
		}
	}
	
	if losses == 0 {
		return 100 // All gains
	}
	
	rs := gains / losses
	rsi := 100 - (100 / (1 + rs))
	
	return rsi
}

// calculatePriceChange calculates recent price change percentage
func calculatePriceChange(prices []float64) float64 {
	if len(prices) < 2 {
		return 0
	}
	
	current := prices[len(prices)-1]
	previous := prices[len(prices)-2]
	
	return (current - previous) / previous * 100
}

// calculateVolatilityScore calculates a volatility score based on technical indicators
func calculateVolatilityScore(indicators map[string]float64, params config.VolatilityConfig) float64 {
	score := 0.0
	
	// Bollinger Band score
	currentPrice := indicators["price"]
	upperBand := indicators["upper_band"]
	lowerBand := indicators["lower_band"]
	
	// Price near or outside Bollinger Bands
	if currentPrice > upperBand*0.98 || currentPrice < lowerBand*1.02 {
		score += 0.3
	}
	
	// RSI score
	rsi := indicators["rsi"]
	if rsi > params.RSIOverbought || rsi < params.RSIOversold {
		score += 0.25
	}
	
	// Volume score
	volumeRatio := indicators["volume_ratio"]
	if volumeRatio > params.VolumeThreshold {
		score += 0.25
	}
	
	// Price change score
	priceChange := math.Abs(indicators["price_change"])
	if priceChange > params.MinVolatilityPercent {
		score += 0.2
	}
	
	return score
}

// determineSignalType determines the signal type based on technical indicators
func determineSignalType(indicators map[string]float64) SignalType {
	// Get indicators
	currentPrice := indicators["price"]
	upperBand := indicators["upper_band"]
	lowerBand := indicators["lower_band"]
	rsi := indicators["rsi"]
	priceChange := indicators["price_change"]
	
	// Bullish conditions
	if (currentPrice < lowerBand*1.02 && rsi < 30) || // Oversold
	   (priceChange > 0 && rsi > 50 && rsi < 70) {    // Uptrend with momentum
		return BUY
	}
	
	// Bearish conditions
	if (currentPrice > upperBand*0.98 && rsi > 70) || // Overbought
	   (priceChange < 0 && rsi < 50 && rsi > 30) {    // Downtrend with momentum
		return SELL
	}
	
	// No clear signal
	return HOLD
}

// calculatePriceLevels calculates target price and stop loss levels
func calculatePriceLevels(currentPrice float64, signalType SignalType, indicators map[string]float64, params config.VolatilityConfig) (float64, float64) {
	// Get indicators
	upperBand := indicators["upper_band"]
	lowerBand := indicators["lower_band"]
	
	var targetPrice, stopLoss float64
	
	if signalType == BUY {
		// Target price: either upper band or a percentage gain
		targetPrice = math.Min(upperBand, currentPrice*(1+params.MinExpectedROI/100))
		
		// Stop loss: either lower band or a percentage loss
		stopLoss = math.Max(lowerBand, currentPrice*(1-params.StopLossPercent/100))
	} else { // SELL
		// Target price: either lower band or a percentage drop
		targetPrice = math.Max(lowerBand, currentPrice*(1-params.MinExpectedROI/100))
		
		// Stop loss: either upper band or a percentage gain
		stopLoss = math.Min(upperBand, currentPrice*(1+params.StopLossPercent/100))
	}
	
	return targetPrice, stopLoss
}

// calculateExpectedROI calculates the expected ROI for a signal
func calculateExpectedROI(currentPrice, targetPrice float64, signalType SignalType) float64 {
	if signalType == BUY {
		return (targetPrice - currentPrice) / currentPrice * 100
	} else { // SELL
		return (currentPrice - targetPrice) / currentPrice * 100
	}
}

// FormatSignalMessage formats a signal for Telegram message
func FormatSignalMessage(s *Signal) string {
	// Format ROI with sign
	roiSign := "+"
	if s.Type == SELL {
		roiSign = "-"
	}
	
	// Format confidence as percentage
	confidencePercent := math.Round(s.Confidence * 100)
	
	// Create message
	message := fmt.Sprintf("🚨 <b>%s SIGNAL: %s</b> 🚨\n\n", s.Type, s.Symbol)
	message += fmt.Sprintf("💰 <b>Entry Price:</b> $%.2f\n", s.Price)
	message += fmt.Sprintf("🎯 <b>Target Price:</b> $%.2f\n", s.TargetPrice)
	message += fmt.Sprintf("🛑 <b>Stop Loss:</b> $%.2f\n", s.StopLoss)
	message += fmt.Sprintf("📈 <b>Expected ROI:</b> %s%.2f%%\n", roiSign, s.ExpectedROI)
	message += fmt.Sprintf("🔍 <b>Confidence:</b> %.0f%%\n", confidencePercent)
	message += fmt.Sprintf("⏱ <b>Time Frame:</b> %s\n\n", s.TimeFrame)
	
	if s.Rationale != "" {
		message += fmt.Sprintf("📝 <b>Rationale:</b>\n%s\n\n", s.Rationale)
	}
	
	message += fmt.Sprintf("⏰ Generated at: %s", s.GeneratedAt.Format("2006-01-02 15:04:05"))
	
	return message
}
