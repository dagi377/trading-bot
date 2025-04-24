package mock

import (
	"context"
	"fmt"
	"sync"

	"github.com/hustler/trading-bot/pkg/config"
	"github.com/hustler/trading-bot/pkg/signal"
	"github.com/hustler/trading-bot/pkg/telegram"
)

// MockTelegramBot is a mock implementation of the Telegram bot
type MockTelegramBot struct {
	messages []string
	mu       sync.RWMutex
}

// NewMockTelegramBot creates a new mock Telegram bot
func NewMockTelegramBot() *telegram.Bot {
	return telegram.NewBotWithMode(config.TelegramConfig{}, true)
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
func (m *MockLLMProvider) GenerateExplanation(ctx context.Context, signal *signal.Signal) (string, error) {
	return "This is a mock explanation for the signal.", nil
}
