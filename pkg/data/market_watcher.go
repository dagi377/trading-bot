package data

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/hustler/trading-bot/pkg/auth"
)

// Stock represents a stock with its current market data
type Stock struct {
	Symbol        string
	CurrentPrice  float64
	PreviousClose float64
	Volume        int64
	LastUpdated   time.Time
	DailyHigh     float64
	DailyLow      float64
	Bid           float64
	Ask           float64
	Change        float64
	ChangePercent float64
}

// MarketWatcher watches real-time market data for a list of stocks
type MarketWatcher struct {
	stocks      map[string]*Stock
	authManager *auth.AuthManager
	dataSource  string
	pollInterval time.Duration
	mu          sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
}

// YahooFinanceResponse represents the response from Yahoo Finance API
type YahooFinanceResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				Currency             string  `json:"currency"`
				Symbol               string  `json:"symbol"`
				RegularMarketPrice   float64 `json:"regularMarketPrice"`
				PreviousClose        float64 `json:"chartPreviousClose"`
				RegularMarketVolume  int64   `json:"regularMarketVolume"`
				RegularMarketDayHigh float64 `json:"regularMarketDayHigh"`
				RegularMarketDayLow  float64 `json:"regularMarketDayLow"`
			} `json:"meta"`
			Timestamp  []int64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Open   []float64 `json:"open"`
					High   []float64 `json:"high"`
					Low    []float64 `json:"low"`
					Close  []float64 `json:"close"`
					Volume []int64   `json:"volume"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
	} `json:"chart"`
}

// AlphaVantageResponse represents the response from Alpha Vantage API
type AlphaVantageResponse struct {
	GlobalQuote struct {
		Symbol           string `json:"01. symbol"`
		Open             string `json:"02. open"`
		High             string `json:"03. high"`
		Low              string `json:"04. low"`
		Price            string `json:"05. price"`
		Volume           string `json:"06. volume"`
		LatestTradingDay string `json:"07. latest trading day"`
		PreviousClose    string `json:"08. previous close"`
		Change           string `json:"09. change"`
		ChangePercent    string `json:"10. change percent"`
	} `json:"Global Quote"`
}

// FinnhubResponse represents the response from Finnhub API
type FinnhubResponse struct {
	CurrentPrice  float64 `json:"c"`
	Change        float64 `json:"d"`
	PercentChange float64 `json:"dp"`
	High          float64 `json:"h"`
	Low           float64 `json:"l"`
	Open          float64 `json:"o"`
	PreviousClose float64 `json:"pc"`
	Timestamp     int64   `json:"t"`
}

// NewMarketWatcher creates a new MarketWatcher instance
func NewMarketWatcher(authManager *auth.AuthManager, dataSource string, pollInterval int) *MarketWatcher {
	ctx, cancel := context.WithCancel(context.Background())
	return &MarketWatcher{
		stocks:       make(map[string]*Stock),
		authManager:  authManager,
		dataSource:   dataSource,
		pollInterval: time.Duration(pollInterval) * time.Second,
		ctx:          ctx,
		cancel:       cancel,
	}
}

// AddStock adds a stock to the watch list
func (m *MarketWatcher) AddStock(symbol string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if _, exists := m.stocks[symbol]; !exists {
		m.stocks[symbol] = &Stock{Symbol: symbol}
	}
}

// RemoveStock removes a stock from the watch list
func (m *MarketWatcher) RemoveStock(symbol string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	delete(m.stocks, symbol)
}

// GetStock returns the current stock data
func (m *MarketWatcher) GetStock(symbol string) (*Stock, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	stock, exists := m.stocks[symbol]
	return stock, exists
}

// GetAllStocks returns all stocks being watched
func (m *MarketWatcher) GetAllStocks() []*Stock {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	stocks := make([]*Stock, 0, len(m.stocks))
	for _, stock := range m.stocks {
		stocks = append(stocks, stock)
	}
	return stocks
}

// StartWatching starts watching all stocks in the watch list
func (m *MarketWatcher) StartWatching() {
	go func() {
		ticker := time.NewTicker(m.pollInterval)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				m.updateAllStocks()
			case <-m.ctx.Done():
				return
			}
		}
	}()
}

// StopWatching stops watching all stocks
func (m *MarketWatcher) StopWatching() {
	m.cancel()
}

// updateAllStocks updates market data for all stocks
func (m *MarketWatcher) updateAllStocks() {
	m.mu.RLock()
	symbols := make([]string, 0, len(m.stocks))
	for symbol := range m.stocks {
		symbols = append(symbols, symbol)
	}
	m.mu.RUnlock()
	
	if len(symbols) == 0 {
		return
	}
	
	for _, symbol := range symbols {
		if err := m.updateStock(symbol); err != nil {
			fmt.Printf("Error updating stock %s: %v\n", symbol, err)
		}
	}
}

// updateStock updates market data for a single stock
func (m *MarketWatcher) updateStock(symbol string) error {
	switch m.dataSource {
	case "yahoo":
		return m.updateStockYahooFinance(symbol)
	case "alphavantage":
		return m.updateStockAlphaVantage(symbol)
	case "finnhub":
		return m.updateStockFinnhub(symbol)
	default:
		return fmt.Errorf("unsupported data source: %s", m.dataSource)
	}
}

// updateStockYahooFinance updates stock data using Yahoo Finance API
func (m *MarketWatcher) updateStockYahooFinance(symbol string) error {
	// Using the YahooFinance/get_stock_chart API from the datasource module
	client := &http.Client{}
	
	// Create the API URL with parameters
	baseURL := "https://query1.finance.yahoo.com/v8/finance/chart/" + symbol
	params := url.Values{}
	params.Add("interval", "1d")
	params.Add("range", "1d")
	
	req, err := http.NewRequest("GET", baseURL+"?"+params.Encode(), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Add("User-Agent", "Mozilla/5.0")
	
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to get data, status: %d, body: %s", resp.StatusCode, string(body))
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	
	var yahooResp YahooFinanceResponse
	if err := json.Unmarshal(body, &yahooResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}
	
	if len(yahooResp.Chart.Result) == 0 {
		return fmt.Errorf("no data found for symbol: %s", symbol)
	}
	
	result := yahooResp.Chart.Result[0]
	meta := result.Meta
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	stock, exists := m.stocks[symbol]
	if !exists {
		return fmt.Errorf("stock not found in watch list: %s", symbol)
	}
	
	stock.CurrentPrice = meta.RegularMarketPrice
	stock.PreviousClose = meta.PreviousClose
	stock.Volume = meta.RegularMarketVolume
	stock.LastUpdated = time.Now()
	stock.DailyHigh = meta.RegularMarketDayHigh
	stock.DailyLow = meta.RegularMarketDayLow
	stock.Change = meta.RegularMarketPrice - meta.PreviousClose
	stock.ChangePercent = (stock.Change / meta.PreviousClose) * 100
	
	// Bid and Ask are not directly available in this API
	// Using current price as an approximation
	stock.Bid = meta.RegularMarketPrice
	stock.Ask = meta.RegularMarketPrice
	
	return nil
}

// updateStockAlphaVantage updates stock data using Alpha Vantage API
func (m *MarketWatcher) updateStockAlphaVantage(symbol string) error {
	apiKey, err := m.authManager.GetAPIKey("alphavantage")
	if err != nil {
		return fmt.Errorf("failed to get Alpha Vantage API key: %w", err)
	}
	
	baseURL := "https://www.alphavantage.co/query"
	params := url.Values{}
	params.Add("function", "GLOBAL_QUOTE")
	params.Add("symbol", symbol)
	params.Add("apikey", apiKey)
	
	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to get data, status: %d, body: %s", resp.StatusCode, string(body))
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	
	var avResp AlphaVantageResponse
	if err := json.Unmarshal(body, &avResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}
	
	quote := avResp.GlobalQuote
	if quote.Symbol == "" {
		return fmt.Errorf("no data found for symbol: %s", symbol)
	}
	
	// Parse string values to float64
	price, _ := strconv.ParseFloat(quote.Price, 64)
	prevClose, _ := strconv.ParseFloat(quote.PreviousClose, 64)
	high, _ := strconv.ParseFloat(quote.High, 64)
	low, _ := strconv.ParseFloat(quote.Low, 64)
	volume, _ := strconv.ParseInt(quote.Volume, 10, 64)
	change, _ := strconv.ParseFloat(quote.Change, 64)
	
	// Remove % from change percent and parse
	changePercentStr := quote.ChangePercent
	if len(changePercentStr) > 0 && changePercentStr[len(changePercentStr)-1] == '%' {
		changePercentStr = changePercentStr[:len(changePercentStr)-1]
	}
	changePercent, _ := strconv.ParseFloat(changePercentStr, 64)
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	stock, exists := m.stocks[symbol]
	if !exists {
		return fmt.Errorf("stock not found in watch list: %s", symbol)
	}
	
	stock.CurrentPrice = price
	stock.PreviousClose = prevClose
	stock.Volume = volume
	stock.LastUpdated = time.Now()
	stock.DailyHigh = high
	stock.DailyLow = low
	stock.Change = change
	stock.ChangePercent = changePercent
	
	// Bid and Ask are not available in this API
	// Using current price as an approximation
	stock.Bid = price
	stock.Ask = price
	
	return nil
}

// updateStockFinnhub updates stock data using Finnhub API
func (m *MarketWatcher) updateStockFinnhub(symbol string) error {
	apiKey, err := m.authManager.GetAPIKey("finnhub")
	if err != nil {
		return fmt.Errorf("failed to get Finnhub API key: %w", err)
	}
	
	baseURL := "https://finnhub.io/api/v1/quote"
	params := url.Values{}
	params.Add("symbol", symbol)
	
	req, err := http.NewRequest("GET", baseURL+"?"+params.Encode(), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Add("X-Finnhub-Token", apiKey)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to get data, status: %d, body: %s", resp.StatusCode, string(body))
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	
	var finnhubResp FinnhubResponse
	if err := json.Unmarshal(body, &finnhubResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}
	
	if finnhubResp.CurrentPrice == 0 {
		return fmt.Errorf("no data found for symbol: %s", symbol)
	}
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	stock, exists := m.stocks[symbol]
	if !exists {
		return fmt.Errorf("stock not found in watch list: %s", symbol)
	}
	
	stock.CurrentPrice = finnhubResp.CurrentPrice
	stock.PreviousClose = finnhubResp.PreviousClose
	stock.Volume = 0 // Not provided in this API response
	stock.LastUpdated = time.Now()
	stock.DailyHigh = finnhubResp.High
	stock.DailyLow = finnhubResp.Low
	stock.Change = finnhubResp.Change
	stock.ChangePercent = finnhubResp.PercentChange
	
	// Bid and Ask are not available in this API
	// Using current price as an approximation
	stock.Bid = finnhubResp.CurrentPrice
	stock.Ask = finnhubResp.CurrentPrice
	
	return nil
}
