package indicators

import (
	"math"
	"sync"

	"github.com/hustler/trading-bot/pkg/data"
)

// Indicator represents a technical indicator
type Indicator interface {
	Calculate(stock *data.Stock) float64
	GetName() string
}

// IndicatorProcessor processes technical indicators for stocks
type IndicatorProcessor struct {
	indicators map[string]map[string]float64
	mu         sync.RWMutex
}

// NewIndicatorProcessor creates a new IndicatorProcessor
func NewIndicatorProcessor() *IndicatorProcessor {
	return &IndicatorProcessor{
		indicators: make(map[string]map[string]float64),
	}
}

// UpdateIndicator updates an indicator value for a stock
func (p *IndicatorProcessor) UpdateIndicator(symbol, indicator string, value float64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.indicators[symbol]; !exists {
		p.indicators[symbol] = make(map[string]float64)
	}

	p.indicators[symbol][indicator] = value
}

// GetIndicator gets an indicator value for a stock
func (p *IndicatorProcessor) GetIndicator(symbol, indicator string) (float64, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if indicators, exists := p.indicators[symbol]; exists {
		value, exists := indicators[indicator]
		return value, exists
	}

	return 0, false
}

// GetAllIndicators gets all indicators for a stock
func (p *IndicatorProcessor) GetAllIndicators(symbol string) map[string]float64 {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if indicators, exists := p.indicators[symbol]; exists {
		result := make(map[string]float64)
		for k, v := range indicators {
			result[k] = v
		}
		return result
	}

	return make(map[string]float64)
}

// RSI represents the Relative Strength Index indicator
type RSI struct {
	period      int
	gains       map[string][]float64
	losses      map[string][]float64
	prevPrices  map[string]float64
	mu          sync.RWMutex
	processor   *IndicatorProcessor
}

// NewRSI creates a new RSI indicator
func NewRSI(period int, processor *IndicatorProcessor) *RSI {
	return &RSI{
		period:     period,
		gains:      make(map[string][]float64),
		losses:     make(map[string][]float64),
		prevPrices: make(map[string]float64),
		processor:  processor,
	}
}

// GetName returns the name of the indicator
func (r *RSI) GetName() string {
	return "RSI"
}

// Calculate calculates the RSI value for a stock
func (r *RSI) Calculate(stock *data.Stock) float64 {
	r.mu.Lock()
	defer r.mu.Unlock()

	symbol := stock.Symbol
	currentPrice := stock.CurrentPrice

	// Initialize if this is the first calculation for this symbol
	if _, exists := r.prevPrices[symbol]; !exists {
		r.prevPrices[symbol] = currentPrice
		r.gains[symbol] = make([]float64, 0, r.period)
		r.losses[symbol] = make([]float64, 0, r.period)
		return 50 // Default neutral value
	}

	// Calculate price change
	prevPrice := r.prevPrices[symbol]
	change := currentPrice - prevPrice
	r.prevPrices[symbol] = currentPrice

	// Update gains and losses
	if change > 0 {
		r.gains[symbol] = append(r.gains[symbol], change)
		r.losses[symbol] = append(r.losses[symbol], 0)
	} else {
		r.gains[symbol] = append(r.gains[symbol], 0)
		r.losses[symbol] = append(r.losses[symbol], math.Abs(change))
	}

	// Trim to period length
	if len(r.gains[symbol]) > r.period {
		r.gains[symbol] = r.gains[symbol][len(r.gains[symbol])-r.period:]
		r.losses[symbol] = r.losses[symbol][len(r.losses[symbol])-r.period:]
	}

	// Not enough data yet
	if len(r.gains[symbol]) < r.period {
		return 50 // Default neutral value
	}

	// Calculate average gain and loss
	var avgGain, avgLoss float64
	for _, gain := range r.gains[symbol] {
		avgGain += gain
	}
	for _, loss := range r.losses[symbol] {
		avgLoss += loss
	}
	avgGain /= float64(r.period)
	avgLoss /= float64(r.period)

	// Calculate RSI
	var rsi float64
	if avgLoss == 0 {
		rsi = 100
	} else {
		rs := avgGain / avgLoss
		rsi = 100 - (100 / (1 + rs))
	}

	// Update the indicator processor
	if r.processor != nil {
		r.processor.UpdateIndicator(symbol, r.GetName(), rsi)
	}

	return rsi
}

// MovingAverage represents a moving average indicator
type MovingAverage struct {
	period    int
	prices    map[string][]float64
	mu        sync.RWMutex
	processor *IndicatorProcessor
	maType    string // "SMA" or "EMA"
}

// NewSMA creates a new Simple Moving Average indicator
func NewSMA(period int, processor *IndicatorProcessor) *MovingAverage {
	return &MovingAverage{
		period:    period,
		prices:    make(map[string][]float64),
		processor: processor,
		maType:    "SMA",
	}
}

// NewEMA creates a new Exponential Moving Average indicator
func NewEMA(period int, processor *IndicatorProcessor) *MovingAverage {
	return &MovingAverage{
		period:    period,
		prices:    make(map[string][]float64),
		processor: processor,
		maType:    "EMA",
	}
}

// GetName returns the name of the indicator
func (m *MovingAverage) GetName() string {
	return m.maType + "-" + string(rune(m.period+'0'))
}

// Calculate calculates the moving average value for a stock
func (m *MovingAverage) Calculate(stock *data.Stock) float64 {
	m.mu.Lock()
	defer m.mu.Unlock()

	symbol := stock.Symbol
	currentPrice := stock.CurrentPrice

	// Initialize if this is the first calculation for this symbol
	if _, exists := m.prices[symbol]; !exists {
		m.prices[symbol] = make([]float64, 0, m.period)
	}

	// Add current price
	m.prices[symbol] = append(m.prices[symbol], currentPrice)

	// Trim to period length
	if len(m.prices[symbol]) > m.period {
		m.prices[symbol] = m.prices[symbol][len(m.prices[symbol])-m.period:]
	}

	// Not enough data yet
	if len(m.prices[symbol]) < m.period {
		return currentPrice // Default to current price
	}

	var ma float64
	if m.maType == "SMA" {
		// Calculate Simple Moving Average
		var sum float64
		for _, price := range m.prices[symbol] {
			sum += price
		}
		ma = sum / float64(m.period)
	} else {
		// Calculate Exponential Moving Average
		k := 2.0 / float64(m.period+1)
		ma = m.prices[symbol][0]
		for i := 1; i < len(m.prices[symbol]); i++ {
			ma = m.prices[symbol][i]*k + ma*(1-k)
		}
	}

	// Update the indicator processor
	if m.processor != nil {
		m.processor.UpdateIndicator(symbol, m.GetName(), ma)
	}

	return ma
}

// VolumeAnalyzer analyzes volume changes
type VolumeAnalyzer struct {
	prevVolumes map[string]int64
	mu          sync.RWMutex
	processor   *IndicatorProcessor
}

// NewVolumeAnalyzer creates a new VolumeAnalyzer
func NewVolumeAnalyzer(processor *IndicatorProcessor) *VolumeAnalyzer {
	return &VolumeAnalyzer{
		prevVolumes: make(map[string]int64),
		processor:   processor,
	}
}

// GetName returns the name of the indicator
func (v *VolumeAnalyzer) GetName() string {
	return "VolumeSurge"
}

// Calculate calculates the volume surge indicator for a stock
func (v *VolumeAnalyzer) Calculate(stock *data.Stock) float64 {
	v.mu.Lock()
	defer v.mu.Unlock()

	symbol := stock.Symbol
	currentVolume := stock.Volume

	// Initialize if this is the first calculation for this symbol
	if _, exists := v.prevVolumes[symbol]; !exists {
		v.prevVolumes[symbol] = currentVolume
		return 0 // No surge yet
	}

	prevVolume := v.prevVolumes[symbol]
	v.prevVolumes[symbol] = currentVolume

	// Calculate volume change percentage
	var volumeChange float64
	if prevVolume > 0 {
		volumeChange = float64(currentVolume-prevVolume) / float64(prevVolume) * 100
	}

	// Update the indicator processor
	if v.processor != nil {
		v.processor.UpdateIndicator(symbol, v.GetName(), volumeChange)
	}

	return volumeChange
}
