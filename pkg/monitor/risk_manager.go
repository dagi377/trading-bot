package monitor

import (
	"fmt"
	"sync"
	"time"

	"github.com/hustler/trading-bot/pkg/data"
	"github.com/hustler/trading-bot/pkg/execution"
)

// RiskManager monitors and enforces risk limits
type RiskManager struct {
	maxDailyLoss    float64
	maxLossPerTrade float64
	dailyPnL        float64
	tradeManager    *execution.TradeManager
	mu              sync.RWMutex
	tradingDay      time.Time
}

// NewRiskManager creates a new RiskManager
func NewRiskManager(maxDailyLoss, maxLossPerTrade float64, tradeManager *execution.TradeManager) *RiskManager {
	return &RiskManager{
		maxDailyLoss:    maxDailyLoss,
		maxLossPerTrade: maxLossPerTrade,
		tradeManager:    tradeManager,
		tradingDay:      time.Now().Truncate(24 * time.Hour),
	}
}

// CheckDailyLoss checks if the daily loss limit has been reached
func (r *RiskManager) CheckDailyLoss(stocks map[string]*data.Stock) (bool, float64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Reset PnL if it's a new trading day
	today := time.Now().Truncate(24 * time.Hour)
	if !today.Equal(r.tradingDay) {
		r.dailyPnL = 0
		r.tradingDay = today
	}

	// Calculate current PnL for all active trades
	currentPnL := r.dailyPnL
	activeTrades := r.tradeManager.GetActiveTrades()

	for _, trade := range activeTrades {
		stock, exists := stocks[trade.Symbol]
		if !exists {
			continue
		}

		// Calculate trade PnL
		entryValue := float64(trade.Quantity) * trade.Price
		currentValue := float64(trade.Quantity) * stock.CurrentPrice
		tradePnL := currentValue - entryValue

		// Add to current PnL
		currentPnL += tradePnL
	}

	// Check if daily loss limit has been reached
	if currentPnL < -r.maxDailyLoss {
		return true, currentPnL
	}

	return false, currentPnL
}

// UpdateDailyPnL updates the daily PnL with a completed trade
func (r *RiskManager) UpdateDailyPnL(buyTrade, sellTrade *execution.Trade) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Reset PnL if it's a new trading day
	today := time.Now().Truncate(24 * time.Hour)
	if !today.Equal(r.tradingDay) {
		r.dailyPnL = 0
		r.tradingDay = today
	}

	// Calculate trade PnL
	buyValue := float64(buyTrade.Quantity) * buyTrade.Price
	sellValue := float64(sellTrade.Quantity) * sellTrade.Price
	tradePnL := sellValue - buyValue

	// Update daily PnL
	r.dailyPnL += tradePnL
}

// GetDailyPnL gets the current daily PnL
func (r *RiskManager) GetDailyPnL() float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.dailyPnL
}

// IsTradingHours checks if it's currently trading hours (9:30 AM - 4:00 PM EST)
func (r *RiskManager) IsTradingHours() bool {
	now := time.Now()
	
	// Convert to EST
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		// If we can't load the location, use a simple approximation
		// Assuming UTC-5 for EST (not accounting for daylight saving)
		offset := -5 * 60 * 60
		now = now.UTC().Add(time.Duration(offset) * time.Second)
	} else {
		now = now.In(loc)
	}
	
	// Check if it's a weekday
	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		return false
	}
	
	// Check if it's between 9:30 AM and 4:00 PM
	marketOpen := time.Date(now.Year(), now.Month(), now.Day(), 9, 30, 0, 0, now.Location())
	marketClose := time.Date(now.Year(), now.Month(), now.Day(), 16, 0, 0, 0, now.Location())
	
	return now.After(marketOpen) && now.Before(marketClose)
}

// ShouldCloseAllPositions checks if all positions should be closed (end of trading day)
func (r *RiskManager) ShouldCloseAllPositions() bool {
	now := time.Now()
	
	// Convert to EST
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		// If we can't load the location, use a simple approximation
		offset := -5 * 60 * 60
		now = now.UTC().Add(time.Duration(offset) * time.Second)
	} else {
		now = now.In(loc)
	}
	
	// Check if it's a weekday
	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		return false
	}
	
	// Check if it's close to 4:00 PM (within 5 minutes)
	marketClose := time.Date(now.Year(), now.Month(), now.Day(), 16, 0, 0, 0, now.Location())
	fiveMinBeforeClose := marketClose.Add(-5 * time.Minute)
	
	return now.After(fiveMinBeforeClose) && now.Before(marketClose)
}

// GenerateRiskReport generates a risk report
func (r *RiskManager) GenerateRiskReport(stocks map[string]*data.Stock) string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	report := "Risk Management Report\n"
	report += "=====================\n\n"
	
	report += fmt.Sprintf("Trading Day: %s\n", r.tradingDay.Format("2006-01-02"))
	report += fmt.Sprintf("Current Daily P&L: $%.2f\n", r.dailyPnL)
	report += fmt.Sprintf("Max Daily Loss Limit: $%.2f\n", r.maxDailyLoss)
	report += fmt.Sprintf("Max Loss Per Trade: $%.2f\n\n", r.maxLossPerTrade)
	
	report += "Active Positions:\n"
	report += "-----------------\n"
	
	activeTrades := r.tradeManager.GetActiveTrades()
	if len(activeTrades) == 0 {
		report += "No active positions\n\n"
	} else {
		for _, trade := range activeTrades {
			stock, exists := stocks[trade.Symbol]
			if !exists {
				continue
			}
			
			entryValue := float64(trade.Quantity) * trade.Price
			currentValue := float64(trade.Quantity) * stock.CurrentPrice
			tradePnL := currentValue - entryValue
			pnlPercent := (tradePnL / entryValue) * 100
			
			report += fmt.Sprintf("Symbol: %s\n", trade.Symbol)
			report += fmt.Sprintf("Quantity: %d\n", trade.Quantity)
			report += fmt.Sprintf("Entry Price: $%.2f\n", trade.Price)
			report += fmt.Sprintf("Current Price: $%.2f\n", stock.CurrentPrice)
			report += fmt.Sprintf("P&L: $%.2f (%.2f%%)\n", tradePnL, pnlPercent)
			report += fmt.Sprintf("Entry Value: $%.2f\n", entryValue)
			report += fmt.Sprintf("Current Value: $%.2f\n", currentValue)
			report += fmt.Sprintf("Time in Trade: %s\n\n", time.Since(trade.CreatedAt).Round(time.Second))
		}
	}
	
	report += "Risk Status:\n"
	report += "-----------\n"
	
	dailyLossLimitReached, _ := r.CheckDailyLoss(stocks)
	if dailyLossLimitReached {
		report += "WARNING: Daily loss limit reached!\n"
	} else {
		report += "Daily loss limit not reached\n"
	}
	
	if r.IsTradingHours() {
		report += "Currently within trading hours\n"
	} else {
		report += "Outside of trading hours\n"
	}
	
	if r.ShouldCloseAllPositions() {
		report += "WARNING: Close to market close, should close all positions\n"
	}
	
	return report
}
