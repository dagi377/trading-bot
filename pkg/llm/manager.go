package llm

import (
	"context"
	"fmt"
	"time"

	"github.com/hustler/trading-bot/pkg/config"
	"github.com/hustler/trading-bot/pkg/signal"
)

// Provider represents an LLM provider
type Provider interface {
	GenerateExplanation(ctx context.Context, s *signal.Signal) (string, error)
	Name() string
}

// Manager manages LLM providers
type Manager struct {
	config   *config.LLMConfig
	provider Provider
}

// NewManager creates a new LLM manager
func NewManager(cfg *config.LLMConfig) (*Manager, error) {
	var provider Provider
	var err error

	switch cfg.Provider {
	case "openai":
		provider, err = NewOpenAIProvider(cfg.APIKey, cfg.ModelName, cfg.MaxTokens, cfg.Temperature)
	case "deepseek":
		provider, err = NewDeepSeekProvider(cfg.LocalPath, cfg.ModelName, cfg.MaxTokens, cfg.Temperature)
	case "mock":
		provider = NewMockProvider()
	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s", cfg.Provider)
	}

	if err != nil {
		return nil, err
	}

	return &Manager{
		config:   cfg,
		provider: provider,
	}, nil
}

// GenerateSignalExplanation generates a natural language explanation for a trading signal
func (m *Manager) GenerateSignalExplanation(ctx context.Context, s *signal.Signal) (string, error) {
	return m.provider.GenerateExplanation(ctx, s)
}

// SwitchProvider switches to a different LLM provider
func (m *Manager) SwitchProvider(providerName string, cfg *config.LLMConfig) error {
	var provider Provider
	var err error

	switch providerName {
	case "openai":
		provider, err = NewOpenAIProvider(cfg.APIKey, cfg.ModelName, cfg.MaxTokens, cfg.Temperature)
	case "deepseek":
		provider, err = NewDeepSeekProvider(cfg.LocalPath, cfg.ModelName, cfg.MaxTokens, cfg.Temperature)
	case "mock":
		provider = NewMockProvider()
	default:
		return fmt.Errorf("unsupported LLM provider: %s", providerName)
	}

	if err != nil {
		return err
	}

	m.provider = provider
	m.config = cfg
	return nil
}

// GetCurrentProvider returns the name of the current provider
func (m *Manager) GetCurrentProvider() string {
	return m.provider.Name()
}

// OpenAIProvider implements the Provider interface for OpenAI
type OpenAIProvider struct {
	apiKey      string
	model       string
	maxTokens   int
	temperature float64
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(apiKey, model string, maxTokens int, temperature float64) (*OpenAIProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required")
	}

	if model == "" {
		model = "gpt-4"
	}

	if maxTokens <= 0 {
		maxTokens = 1000
	}

	if temperature < 0 || temperature > 1 {
		temperature = 0.7
	}

	return &OpenAIProvider{
		apiKey:      apiKey,
		model:       model,
		maxTokens:   maxTokens,
		temperature: temperature,
	}, nil
}

// GenerateExplanation generates a natural language explanation using OpenAI
func (p *OpenAIProvider) GenerateExplanation(ctx context.Context, s *signal.Signal) (string, error) {
	// In a real implementation, this would send a request to the OpenAI API
	// For now, we'll return a mock response
	time.Sleep(500 * time.Millisecond) // Simulate processing time

	// Generate a mock explanation based on the signal
	content := generateMockExplanation(s)

	return content, nil
}

// Name returns the provider name
func (p *OpenAIProvider) Name() string {
	return "openai"
}

// DeepSeekProvider implements the Provider interface for DeepSeek
type DeepSeekProvider struct {
	localPath   string
	model       string
	maxTokens   int
	temperature float64
}

// NewDeepSeekProvider creates a new DeepSeek provider
func NewDeepSeekProvider(localPath, model string, maxTokens int, temperature float64) (*DeepSeekProvider, error) {
	if localPath == "" {
		return nil, fmt.Errorf("DeepSeek local path is required")
	}

	if model == "" {
		model = "deepseek-coder"
	}

	if maxTokens <= 0 {
		maxTokens = 1000
	}

	if temperature < 0 || temperature > 1 {
		temperature = 0.7
	}

	return &DeepSeekProvider{
		localPath:   localPath,
		model:       model,
		maxTokens:   maxTokens,
		temperature: temperature,
	}, nil
}

// GenerateExplanation generates a natural language explanation using DeepSeek
func (p *DeepSeekProvider) GenerateExplanation(ctx context.Context, s *signal.Signal) (string, error) {
	// In a real implementation, this would use the DeepSeek API or local model
	// For now, we'll implement a simplified version that returns a mock response

	// In a real implementation, this would send a request to the local DeepSeek server
	// For now, we'll return a mock response
	time.Sleep(500 * time.Millisecond) // Simulate processing time

	// Generate a mock explanation based on the signal
	explanation := generateMockExplanation(s)

	return explanation, nil
}

// Name returns the provider name
func (p *DeepSeekProvider) Name() string {
	return "deepseek"
}

// MockProvider implements the Provider interface for testing
type MockProvider struct{}

// NewMockProvider creates a new mock provider
func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

// GenerateExplanation generates a mock explanation
func (p *MockProvider) GenerateExplanation(ctx context.Context, s *signal.Signal) (string, error) {
	return generateMockExplanation(s), nil
}

// Name returns the provider name
func (p *MockProvider) Name() string {
	return "mock"
}

// Helper functions

// createSignalPrompt creates a prompt for the LLM based on the signal
func createSignalPrompt(s *signal.Signal) string {
	// Format technical data
	technicalData := ""
	for key, value := range s.TechnicalData {
		technicalData += fmt.Sprintf("- %s: %.2f\n", key, value)
	}

	// Create prompt
	prompt := fmt.Sprintf(`
Analyze the following trading signal and provide a clear, concise explanation for why this signal was generated and what it means for traders.

Signal Details:
- Symbol: %s
- Type: %s
- Current Price: $%.2f
- Target Price: $%.2f
- Stop Loss: $%.2f
- Expected ROI: %.2f%%
- Confidence: %.0f%%
- Time Frame: %s

Technical Indicators:
%s

Based on these details, explain:
1. Why this %s signal was generated
2. What technical factors support this signal
3. What risks to be aware of
4. How traders should approach this opportunity

Keep your explanation concise, informative, and suitable for both novice and experienced traders.
`, s.Symbol, s.Type, s.Price, s.TargetPrice, s.StopLoss, s.ExpectedROI, s.Confidence*100, s.TimeFrame, technicalData, s.Type)

	return prompt
}

// generateMockExplanation generates a mock explanation based on the signal
func generateMockExplanation(s *signal.Signal) string {
	var explanation string

	if s.Type == signal.BUY {
		explanation = fmt.Sprintf(`
This BUY signal for %s is based on a strong volatility pattern indicating potential upward movement. 

Key factors supporting this signal:
1. The price has shown increased volatility with a bullish bias
2. Technical indicators suggest the stock is currently undervalued
3. Volume has increased significantly, confirming buying interest

With a target price of $%.2f and stop loss at $%.2f, this trade offers a favorable risk-reward ratio of approximately %.1f:1. The %.0f%% confidence score indicates a relatively strong signal based on our algorithm's analysis.

Traders should consider entering this position soon, as the expected timeframe for this movement is %s. However, always adhere to your personal risk management rules and consider the broader market context before entering any trade.
`, s.Symbol, s.TargetPrice, s.StopLoss, s.ExpectedROI/(s.Price-s.StopLoss)*100, s.Confidence*100, s.TimeFrame)
	} else {
		explanation = fmt.Sprintf(`
This SELL signal for %s is based on a volatility pattern indicating potential downward movement.

Key factors supporting this signal:
1. The price has shown increased volatility with a bearish bias
2. Technical indicators suggest the stock may be currently overvalued
3. Recent price action shows weakening momentum

With a target price of $%.2f and stop loss at $%.2f, this trade offers a favorable risk-reward ratio of approximately %.1f:1. The %.0f%% confidence score indicates a relatively strong signal based on our algorithm's analysis.

Traders should consider entering this short position soon, as the expected timeframe for this movement is %s. However, always remember that short positions carry additional risks, and you should adhere to strict risk management practices.
`, s.Symbol, s.TargetPrice, s.StopLoss, s.ExpectedROI/(s.StopLoss-s.Price)*100, s.Confidence*100, s.TimeFrame)
	}

	return explanation
}
