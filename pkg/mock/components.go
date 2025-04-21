package mock

import (
	"context"
	"fmt"
	"sync"

	"github.com/hustler/trading-bot/pkg/signal"
)

// MockTelegramBot is a mock implementation of the Telegram bot
type MockTelegramBot struct {
	messages []string
	mu       sync.RWMutex
}

// NewMockTelegramBot creates a new mock Telegram bot
func NewMockTelegramBot() *MockTelegramBot {
	return &MockTelegramBot{
		messages: []string{},
		mu:       sync.RWMutex{},
	}
}

// SendSignal sends a signal to the mock Telegram bot
func (m *MockTelegramBot) SendSignal(s *signal.Signal) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Format signal message
	message := formatSignalMessage(s)
	
	// Store message
	m.messages = append(m.messages, message)
	
	return nil
}

// ProcessUpdates processes mock updates
func (m *MockTelegramBot) ProcessUpdates() error {
	// No-op for mock
	return nil
}

// GetSentMessages returns all sent messages
func (m *MockTelegramBot) GetSentMessages() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	// Return a copy to avoid race conditions
	messagesCopy := make([]string, len(m.messages))
	copy(messagesCopy, m.messages)
	
	return messagesCopy
}

// formatSignalMessage formats a signal for Telegram message
func formatSignalMessage(s *signal.Signal) string {
	// Use the same format as the real Telegram bot
	if s.Type == signal.BUY {
		return formatBuySignal(s)
	} else {
		return formatSellSignal(s)
	}
}

// formatBuySignal formats a BUY signal
func formatBuySignal(s *signal.Signal) string {
	message := "🚨 BUY SIGNAL: " + s.Symbol + " 🚨\n\n"
	message += "💰 Entry Price: $" + formatFloat(s.Price) + "\n"
	message += "🎯 Target Price: $" + formatFloat(s.TargetPrice) + "\n"
	message += "🛑 Stop Loss: $" + formatFloat(s.StopLoss) + "\n"
	message += "📈 Expected ROI: +" + formatFloat(s.ExpectedROI) + "%\n"
	message += "🔍 Confidence: " + formatConfidence(s.Confidence) + "%\n"
	message += "⏱ Time Frame: " + s.TimeFrame + "\n\n"
	
	if s.Rationale != "" {
		message += "📝 Rationale:\n" + s.Rationale + "\n\n"
	}
	
	message += "⏰ Generated at: " + s.GeneratedAt.Format("2006-01-02 15:04:05")
	
	return message
}

// formatSellSignal formats a SELL signal
func formatSellSignal(s *signal.Signal) string {
	message := "🚨 SELL SIGNAL: " + s.Symbol + " 🚨\n\n"
	message += "💰 Entry Price: $" + formatFloat(s.Price) + "\n"
	message += "🎯 Target Price: $" + formatFloat(s.TargetPrice) + "\n"
	message += "🛑 Stop Loss: $" + formatFloat(s.StopLoss) + "\n"
	message += "📈 Expected ROI: -" + formatFloat(s.ExpectedROI) + "%\n"
	message += "🔍 Confidence: " + formatConfidence(s.Confidence) + "%\n"
	message += "⏱ Time Frame: " + s.TimeFrame + "\n\n"
	
	if s.Rationale != "" {
		message += "📝 Rationale:\n" + s.Rationale + "\n\n"
	}
	
	message += "⏰ Generated at: " + s.GeneratedAt.Format("2006-01-02 15:04:05")
	
	return message
}

// Helper functions
func formatFloat(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

func formatConfidence(c float64) string {
	return fmt.Sprintf("%.0f", c*100)
}

// MockLLMProvider is a mock implementation of the LLM provider
type MockLLMProvider struct{}

// NewMockLLMProvider creates a new mock LLM provider
func NewMockLLMProvider() *MockLLMProvider {
	return &MockLLMProvider{}
}

// GenerateExplanation generates a mock explanation
func (p *MockLLMProvider) GenerateExplanation(ctx context.Context, s *signal.Signal) (string, error) {
	if s.Type == signal.BUY {
		return generateMockBuyExplanation(s), nil
	} else {
		return generateMockSellExplanation(s), nil
	}
}

// Name returns the provider name
func (p *MockLLMProvider) Name() string {
	return "mock"
}

// generateMockBuyExplanation generates a mock explanation for a BUY signal
func generateMockBuyExplanation(s *signal.Signal) string {
	return fmt.Sprintf(`
This BUY signal for %s is based on a clear volatility pattern indicating potential upward movement in the short term.

Key factors supporting this signal:
1. The price has shown increased volatility with a bullish bias
2. Technical indicators suggest the stock is currently undervalued
3. Volume has increased significantly, confirming buying interest

With a target price of $%.2f and stop loss at $%.2f, this trade offers a favorable risk-reward ratio of approximately %.1f:1. The %.0f%% confidence score indicates a relatively strong signal based on our algorithm's analysis.

Traders should consider entering this position soon, as the expected timeframe for this movement is %s. However, always adhere to your personal risk management rules and consider the broader market context before entering any trade.
`, s.Symbol, s.TargetPrice, s.StopLoss, s.ExpectedROI/(s.Price-s.StopLoss)*100, s.Confidence*100, s.TimeFrame)
}

// generateMockSellExplanation generates a mock explanation for a SELL signal
func generateMockSellExplanation(s *signal.Signal) string {
	return fmt.Sprintf(`
This SELL signal for %s is based on a volatility pattern indicating potential downward movement in the short term.

Key factors supporting this signal:
1. The price has shown increased volatility with a bearish bias
2. Technical indicators suggest the stock may be currently overvalued
3. Recent price action shows weakening momentum

With a target price of $%.2f and stop loss at $%.2f, this trade offers a favorable risk-reward ratio of approximately %.1f:1. The %.0f%% confidence score indicates a relatively strong signal based on our algorithm's analysis.

Traders should consider entering this short position soon, as the expected timeframe for this movement is %s. However, always remember that short positions carry additional risks, and you should adhere to strict risk management practices.
`, s.Symbol, s.TargetPrice, s.StopLoss, s.ExpectedROI/(s.StopLoss-s.Price)*100, s.Confidence*100, s.TimeFrame)
}
