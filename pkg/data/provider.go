package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/hustler/trading-bot/pkg/config"
)

// Provider handles fetching market data from various sources
type Provider struct {
	config *config.Config
}

// MarketData represents market data for a stock
type MarketData struct {
	Symbol     string
	Prices     []float64
	Volumes    []float64
	Timestamps []time.Time
}

// NewProvider creates a new data provider
func NewProvider(cfg *config.Config) *Provider {
	return &Provider{
		config: cfg,
	}
}

// GetMarketData fetches market data for a symbol
func (p *Provider) GetMarketData(symbol string) (*MarketData, error) {
	// Determine which data source to use
	primary := p.config.DataSource.Primary
	
	var data *MarketData
	var err error
	
	// Try primary source
	switch primary {
	case "yahoo":
		data, err = p.fetchYahooFinanceData(symbol)
	case "alphavantage":
		data, err = p.fetchAlphaVantageData(symbol)
	default:
		return nil, fmt.Errorf("unsupported primary data source: %s", primary)
	}
	
	// If primary source fails, try secondary source
	if err != nil {
		secondary := p.config.DataSource.Secondary
		
		switch secondary {
		case "yahoo":
			data, err = p.fetchYahooFinanceData(symbol)
		case "alphavantage":
			data, err = p.fetchAlphaVantageData(symbol)
		default:
			return nil, fmt.Errorf("primary source failed and unsupported secondary data source: %s", secondary)
		}
		
		if err != nil {
			return nil, fmt.Errorf("both primary and secondary data sources failed: %w", err)
		}
	}
	
	return data, nil
}

// fetchYahooFinanceData fetches data from Yahoo Finance API
func (p *Provider) fetchYahooFinanceData(symbol string) (*MarketData, error) {
	// In a real implementation, this would use the Yahoo Finance API
	// For now, we'll use the data API provided in the environment
	
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	// Create request
	url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s", symbol)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Add query parameters
	q := req.URL.Query()
	q.Add("interval", "5m")
	q.Add("range", "1d")
	req.URL.RawQuery = q.Encode()
	
	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}
	
	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	
	// Parse response
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	// Extract data
	chart, ok := response["chart"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format: missing chart")
	}
	
	result, ok := chart["result"].([]interface{})
	if !ok || len(result) == 0 {
		return nil, fmt.Errorf("invalid response format: missing result")
	}
	
	// For now, we'll return mock data since we can't actually call the API
	return createMockMarketData(symbol), nil
}

// fetchAlphaVantageData fetches data from Alpha Vantage API
func (p *Provider) fetchAlphaVantageData(symbol string) (*MarketData, error) {
	// In a real implementation, this would use the Alpha Vantage API
	// For now, we'll return mock data
	
	// Get API key
	apiKey, ok := p.config.DataSource.APIKeys["alphavantage"]
	if !ok || apiKey == "" {
		return nil, fmt.Errorf("Alpha Vantage API key not found")
	}
	
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	// Create request
	url := "https://www.alphavantage.co/query"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Add query parameters
	q := req.URL.Query()
	q.Add("function", "TIME_SERIES_INTRADAY")
	q.Add("symbol", symbol)
	q.Add("interval", "5min")
	q.Add("apikey", apiKey)
	req.URL.RawQuery = q.Encode()
	
	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}
	
	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	
	// Parse response
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	// Check for error message
	if errorMsg, ok := response["Error Message"]; ok {
		return nil, fmt.Errorf("API error: %s", errorMsg)
	}
	
	// For now, we'll return mock data since we can't actually call the API
	return createMockMarketData(symbol), nil
}

// createMockMarketData creates mock market data for testing
func createMockMarketData(symbol string) *MarketData {
	// Create base price based on symbol
	var basePrice float64
	switch symbol {
	case "AAPL":
		basePrice = 175.0
	case "MSFT":
		basePrice = 350.0
	case "GOOGL":
		basePrice = 180.0
	case "AMZN":
		basePrice = 125.0
	case "META":
		basePrice = 300.0
	default:
		basePrice = 100.0
	}
	
	// Create mock data
	now := time.Now()
	dataPoints := 78 // 6.5 hours of 5-minute data
	
	prices := make([]float64, dataPoints)
	volumes := make([]float64, dataPoints)
	timestamps := make([]time.Time, dataPoints)
	
	// Generate data with some randomness and trend
	for i := 0; i < dataPoints; i++ {
		// Calculate time (going backward from now)
		timestamps[dataPoints-1-i] = now.Add(-time.Duration(i*5) * time.Minute)
		
		// Calculate price with some randomness
		randomFactor := 0.002 * (float64(i%10) - 5.0) // -0.01 to 0.01
		trendFactor := 0.0001 * float64(i) // Small upward trend
		
		if i == 0 {
			prices[dataPoints-1-i] = basePrice
		} else {
			prevPrice := prices[dataPoints-i]
			priceChange := prevPrice * (randomFactor + trendFactor)
			prices[dataPoints-1-i] = prevPrice + priceChange
		}
		
		// Calculate volume with some randomness
		baseVolume := 1000000.0
		volumeFactor := 0.5 + float64(i%10)/10.0 // 0.5 to 1.4
		volumes[dataPoints-1-i] = baseVolume * volumeFactor
	}
	
	// Add a volatility spike for testing signal generation
	spikeIndex := dataPoints / 2
	prices[spikeIndex] = prices[spikeIndex-1] * 1.02 // 2% spike
	volumes[spikeIndex] = volumes[spikeIndex-1] * 2.0 // Volume spike
	
	return &MarketData{
		Symbol:     symbol,
		Prices:     prices,
		Volumes:    volumes,
		Timestamps: timestamps,
	}
}

// UpdateConfig updates the provider configuration
func (p *Provider) UpdateConfig(cfg *config.Config) {
	p.config = cfg
}
