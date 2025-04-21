package llm

import (
	"context"
	"testing"
	"time"

	"github.com/hustler/trading-bot/pkg/config"
	"github.com/hustler/trading-bot/pkg/signal"
	"github.com/stretchr/testify/assert"
)

func TestNewManager(t *testing.T) {
	// Test with mock provider
	cfg := &config.LLMConfig{
		Provider:    "mock",
		APIKey:      "",
		ModelName:   "test-model",
		MaxTokens:   500,
		Temperature: 0.5,
	}

	manager, err := NewManager(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, manager)
	assert.Equal(t, "mock", manager.GetCurrentProvider())

	// Test with unsupported provider
	cfg.Provider = "unsupported"
	manager, err = NewManager(cfg)
	assert.Error(t, err)
	assert.Nil(t, manager)
	assert.Contains(t, err.Error(), "unsupported LLM provider")

	// Test with OpenAI provider but missing API key
	cfg.Provider = "openai"
	cfg.APIKey = ""
	manager, err = NewManager(cfg)
	assert.Error(t, err)
	assert.Nil(t, manager)
	assert.Contains(t, err.Error(), "OpenAI API key is required")

	// Test with OpenAI provider with API key
	cfg.Provider = "openai"
	cfg.APIKey = "test-api-key"
	manager, err = NewManager(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, manager)
	assert.Equal(t, "openai", manager.GetCurrentProvider())

	// Test with DeepSeek provider but missing local path
	cfg.Provider = "deepseek"
	cfg.LocalPath = ""
	manager, err = NewManager(cfg)
	assert.Error(t, err)
	assert.Nil(t, manager)
	assert.Contains(t, err.Error(), "DeepSeek local path is required")

	// Test with DeepSeek provider with local path
	cfg.Provider = "deepseek"
	cfg.LocalPath = "/path/to/model"
	manager, err = NewManager(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, manager)
	assert.Equal(t, "deepseek", manager.GetCurrentProvider())
}

func TestSwitchProvider(t *testing.T) {
	// Start with mock provider
	cfg := &config.LLMConfig{
		Provider:    "mock",
		APIKey:      "",
		ModelName:   "test-model",
		MaxTokens:   500,
		Temperature: 0.5,
	}

	manager, err := NewManager(cfg)
	assert.NoError(t, err)
	assert.Equal(t, "mock", manager.GetCurrentProvider())

	// Switch to OpenAI
	newCfg := &config.LLMConfig{
		Provider:    "openai",
		APIKey:      "test-api-key",
		ModelName:   "gpt-4",
		MaxTokens:   1000,
		Temperature: 0.7,
	}
	err = manager.SwitchProvider("openai", newCfg)
	assert.NoError(t, err)
	assert.Equal(t, "openai", manager.GetCurrentProvider())

	// Switch to DeepSeek
	newCfg = &config.LLMConfig{
		Provider:    "deepseek",
		LocalPath:   "/path/to/model",
		ModelName:   "deepseek-coder",
		MaxTokens:   1000,
		Temperature: 0.7,
	}
	err = manager.SwitchProvider("deepseek", newCfg)
	assert.NoError(t, err)
	assert.Equal(t, "deepseek", manager.GetCurrentProvider())

	// Switch to unsupported provider
	err = manager.SwitchProvider("unsupported", newCfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported LLM provider")
	assert.Equal(t, "deepseek", manager.GetCurrentProvider()) // Should not change
}

func TestGenerateSignalExplanation(t *testing.T) {
	// Create manager with mock provider
	cfg := &config.LLMConfig{
		Provider:    "mock",
		ModelName:   "test-model",
		MaxTokens:   500,
		Temperature: 0.5,
	}

	manager, err := NewManager(cfg)
	assert.NoError(t, err)

	// Create test signal
	testSignal := &signal.Signal{
		Symbol:        "AAPL",
		Type:          signal.BUY,
		Price:         150.25,
		TargetPrice:   155.50,
		StopLoss:      148.00,
		ExpectedROI:   3.5,
		Confidence:    0.85,
		GeneratedAt:   time.Now(),
		TimeFrame:     "1-3 hours",
		TechnicalData: map[string]float64{"RSI": 35.0, "Volume": 1500000},
	}

	// Generate explanation
	ctx := context.Background()
	explanation, err := manager.GenerateSignalExplanation(ctx, testSignal)
	assert.NoError(t, err)
	assert.NotEmpty(t, explanation)

	// Verify explanation contains key information
	assert.Contains(t, explanation, "BUY signal for AAPL")
	assert.Contains(t, explanation, "target price of $155.50")
	assert.Contains(t, explanation, "stop loss at $148.00")

	// Test with SELL signal
	testSignal.Type = signal.SELL
	testSignal.TargetPrice = 145.00
	testSignal.StopLoss = 152.00
	explanation, err = manager.GenerateSignalExplanation(ctx, testSignal)
	assert.NoError(t, err)
	assert.NotEmpty(t, explanation)
	assert.Contains(t, explanation, "SELL signal for AAPL")
}

func TestMockProvider(t *testing.T) {
	provider := NewMockProvider()
	assert.NotNil(t, provider)
	assert.Equal(t, "mock", provider.Name())

	// Create test signal
	testSignal := &signal.Signal{
		Symbol:        "MSFT",
		Type:          signal.BUY,
		Price:         350.75,
		TargetPrice:   360.00,
		StopLoss:      345.00,
		ExpectedROI:   2.6,
		Confidence:    0.75,
		GeneratedAt:   time.Now(),
		TimeFrame:     "2-4 hours",
		TechnicalData: map[string]float64{"RSI": 40.0, "Volume": 2000000},
	}

	// Generate explanation
	ctx := context.Background()
	explanation, err := provider.GenerateExplanation(ctx, testSignal)
	assert.NoError(t, err)
	assert.NotEmpty(t, explanation)
	assert.Contains(t, explanation, "BUY signal for MSFT")
}

func TestCreateSignalPrompt(t *testing.T) {
	// Create test signal
	testSignal := &signal.Signal{
		Symbol:      "GOOGL",
		Type:        signal.SELL,
		Price:       180.50,
		TargetPrice: 175.00,
		StopLoss:    183.00,
		ExpectedROI: 3.0,
		Confidence:  0.80,
		TimeFrame:   "1-2 hours",
		TechnicalData: map[string]float64{
			"RSI":          75.5,
			"Volume":       3000000,
			"price_change": -1.2,
		},
	}

	// Create prompt
	prompt := createSignalPrompt(testSignal)

	// Verify prompt contains key information
	assert.Contains(t, prompt, "Symbol: GOOGL")
	assert.Contains(t, prompt, "Type: SELL")
	assert.Contains(t, prompt, "Current Price: $180.50")
	assert.Contains(t, prompt, "Target Price: $175.00")
	assert.Contains(t, prompt, "Stop Loss: $183.00")
	assert.Contains(t, prompt, "Expected ROI: 3.00%")
	assert.Contains(t, prompt, "Confidence: 80%")
	assert.Contains(t, prompt, "Time Frame: 1-2 hours")
	assert.Contains(t, prompt, "RSI: 75.50")
	assert.Contains(t, prompt, "Volume: 3000000.00")
	assert.Contains(t, prompt, "price_change: -1.20")
	assert.Contains(t, prompt, "Why this SELL signal was generated")
}
