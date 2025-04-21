package monitor

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/hustler/trading-bot/pkg/config"
	"github.com/hustler/trading-bot/pkg/data"
	"github.com/hustler/trading-bot/pkg/llm"
	"github.com/hustler/trading-bot/pkg/signal"
	"github.com/hustler/trading-bot/pkg/telegram"
)

// MarketMonitor monitors the market and generates trading signals
type MarketMonitor struct {
	config        *config.Config
	dataProvider  *data.Provider
	signalGen     *signal.Generator
	llmManager    *llm.Manager
	telegramBot   *telegram.Bot
	isRunning     bool
	stopChan      chan struct{}
	signalHistory []*signal.Signal
	mu            sync.RWMutex
}

// NewMarketMonitor creates a new market monitor
func NewMarketMonitor(
	cfg *config.Config,
	dataProvider *data.Provider,
	signalGen *signal.Generator,
	llmManager *llm.Manager,
	telegramBot *telegram.Bot,
) *MarketMonitor {
	return &MarketMonitor{
		config:        cfg,
		dataProvider:  dataProvider,
		signalGen:     signalGen,
		llmManager:    llmManager,
		telegramBot:   telegramBot,
		isRunning:     false,
		stopChan:      make(chan struct{}),
		signalHistory: []*signal.Signal{},
		mu:            sync.RWMutex{},
	}
}

// Start starts the market monitor
func (m *MarketMonitor) Start() error {
	m.mu.Lock()
	if m.isRunning {
		m.mu.Unlock()
		return fmt.Errorf("market monitor is already running")
	}
	m.isRunning = true
	m.stopChan = make(chan struct{})
	m.mu.Unlock()

	log.Println("Starting market monitor")

	// Start monitoring in a goroutine
	go m.monitorMarket()

	return nil
}

// Stop stops the market monitor
func (m *MarketMonitor) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.isRunning {
		return fmt.Errorf("market monitor is not running")
	}

	log.Println("Stopping market monitor")
	close(m.stopChan)
	m.isRunning = false

	return nil
}

// IsRunning returns whether the market monitor is running
func (m *MarketMonitor) IsRunning() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.isRunning
}

// GetSignalHistory returns the signal history
func (m *MarketMonitor) GetSignalHistory() []*signal.Signal {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return a copy to avoid race conditions
	history := make([]*signal.Signal, len(m.signalHistory))
	copy(history, m.signalHistory)

	return history
}

// monitorMarket monitors the market and generates signals
func (m *MarketMonitor) monitorMarket() {
	// Calculate initial check time
	nextCheckTime := time.Now()

	for {
		select {
		case <-m.stopChan:
			log.Println("Market monitor stopped")
			return
		case <-time.After(time.Until(nextCheckTime)):
			// Check if within trading hours
			withinHours, err := m.config.IsWithinTradingHours()
			if err != nil {
				log.Printf("Error checking trading hours: %v", err)
				nextCheckTime = time.Now().Add(time.Minute) // Retry in 1 minute
				continue
			}

			if !withinHours {
				log.Println("Outside trading hours, skipping check")
				// Calculate next check time (next minute)
				nextCheckTime = time.Now().Add(time.Minute)
				continue
			}

			// Perform market check
			log.Println("Performing market check")
			err = m.performMarketCheck()
			if err != nil {
				log.Printf("Error performing market check: %v", err)
			}

			// Calculate next check time
			nextCheckTime = time.Now().Add(time.Duration(m.config.CheckInterval) * time.Second)
		}
	}
}

// performMarketCheck performs a market check and generates signals
func (m *MarketMonitor) performMarketCheck() error {
	// Get stock symbols
	m.mu.RLock()
	symbols := m.config.StockSymbols
	m.mu.RUnlock()

	// Fetch market data for all symbols
	marketData := make(map[string]signal.MarketData)
	for _, symbol := range symbols {
		data, err := m.dataProvider.GetMarketData(symbol)
		if err != nil {
			log.Printf("Error fetching market data for %s: %v", symbol, err)
			continue
		}
		marketData[symbol] = signal.MarketData{
			Symbol:     symbol,
			Prices:     data.Prices,
			Volumes:    data.Volumes,
			Timestamps: data.Timestamps,
		}
	}

	// Generate signals
	signals, err := m.signalGen.GenerateSignals(marketData)
	if err != nil {
		return fmt.Errorf("error generating signals: %w", err)
	}

	// Process signals
	for _, s := range signals {
		// Generate explanation using LLM
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		explanation, err := m.llmManager.GenerateSignalExplanation(ctx, s)
		cancel()
		if err != nil {
			log.Printf("Error generating explanation for signal %s: %v", s.ID, err)
		} else {
			s.Rationale = explanation
		}

		// Send signal to Telegram
		err = m.telegramBot.SendSignal(s)
		if err != nil {
			log.Printf("Error sending signal to Telegram: %v", err)
		}

		// Add signal to history
		m.mu.Lock()
		m.signalHistory = append(m.signalHistory, s)
		// Limit history size to 100 signals
		if len(m.signalHistory) > 100 {
			m.signalHistory = m.signalHistory[len(m.signalHistory)-100:]
		}
		m.mu.Unlock()

		log.Printf("Generated and sent %s signal for %s", s.Type, s.Symbol)
	}

	log.Printf("Market check completed, generated %d signals", len(signals))
	return nil
}

// UpdateConfig updates the monitor configuration
func (m *MarketMonitor) UpdateConfig(cfg *config.Config) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.config = cfg
}
