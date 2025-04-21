package strategy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/hustler/trading-bot/pkg/data"
	"github.com/hustler/trading-bot/pkg/indicators"
)

// TradeSignal represents a trading signal
type TradeSignal string

const (
	Buy  TradeSignal = "BUY"
	Sell TradeSignal = "SELL"
	Hold TradeSignal = "HOLD"
)

// TradeDecision represents a trading decision
type TradeDecision struct {
	Symbol    string
	Signal    TradeSignal
	Price     float64
	Timestamp time.Time
	Rationale string
	Score     float64
}

// LLMConfig represents the configuration for the LLM
type LLMConfig struct {
	Provider   string // "openai", "anthropic", "deepseek", or "mock"
	APIKey     string
	ModelName  string
	LocalPath  string // Path to local model (for deepseek)
	MaxTokens  int
	Temperature float64
}

// LLMAdvisor uses an LLM to provide trading advice
type LLMAdvisor struct {
	config       LLMConfig
	indicatorProc *indicators.IndicatorProcessor
	mu           sync.Mutex
}

// NewLLMAdvisor creates a new LLMAdvisor
func NewLLMAdvisor(config LLMConfig, indicatorProc *indicators.IndicatorProcessor) *LLMAdvisor {
	// Set default values if not provided
	if config.ModelName == "" {
		switch config.Provider {
		case "openai":
			config.ModelName = "gpt-4"
		case "anthropic":
			config.ModelName = "claude-3-opus-20240229"
		case "deepseek":
			config.ModelName = "deepseek-coder"
		case "mock":
			config.ModelName = "mock-model"
		}
	}

	if config.MaxTokens == 0 {
		config.MaxTokens = 1000
	}

	if config.Temperature == 0 {
		config.Temperature = 0.7
	}

	return &LLMAdvisor{
		config:       config,
		indicatorProc: indicatorProc,
	}
}

// OpenAIRequest represents a request to the OpenAI API
type OpenAIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
}

// Message represents a message in the OpenAI API request
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIResponse represents a response from the OpenAI API
type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// AnthropicRequest represents a request to the Anthropic API
type AnthropicRequest struct {
	Model       string  `json:"model"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature,omitempty"`
	Messages    []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}

// AnthropicResponse represents a response from the Anthropic API
type AnthropicResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
}

// GetTradeAdvice gets trading advice from the LLM
func (l *LLMAdvisor) GetTradeAdvice(stock *data.Stock) (*TradeDecision, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Get all indicators for the stock
	indicators := l.indicatorProc.GetAllIndicators(stock.Symbol)

	// Prepare context for LLM
	context := fmt.Sprintf(`
Stock: %s
Current Price: $%.2f
Previous Close: $%.2f
Change: $%.2f (%.2f%%)
Daily High: $%.2f
Daily Low: $%.2f
Volume: %d
Last Updated: %s

Technical Indicators:
`,
		stock.Symbol,
		stock.CurrentPrice,
		stock.PreviousClose,
		stock.Change,
		stock.ChangePercent,
		stock.DailyHigh,
		stock.DailyLow,
		stock.Volume,
		stock.LastUpdated.Format(time.RFC3339),
	)

	// Add indicators to context
	for name, value := range indicators {
		context += fmt.Sprintf("%s: %.2f\n", name, value)
	}

	// Add prompt for LLM
	prompt := context + `
Based on the above market data and technical indicators, provide a trading recommendation (BUY, SELL, or HOLD) for this stock.
Include a brief rationale for your recommendation. Format your response as JSON with the following structure:
{
  "signal": "BUY|SELL|HOLD",
  "rationale": "Your reasoning here",
  "confidence": 0.XX
}
`

	var response string
	var err error

	switch l.config.Provider {
	case "openai":
		response, err = l.callOpenAI(prompt)
	case "anthropic":
		response, err = l.callAnthropic(prompt)
	case "deepseek":
		response, err = l.callDeepSeek(prompt)
	case "mock":
		response, err = l.mockLLMResponse(stock)
	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s", l.config.Provider)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get LLM response: %w", err)
	}

	// Parse LLM response
	var result struct {
		Signal     string  `json:"signal"`
		Rationale  string  `json:"rationale"`
		Confidence float64 `json:"confidence"`
	}

	if err := json.Unmarshal([]byte(response), &result); err != nil {
		// If JSON parsing fails, try to extract JSON from the response
		jsonStart := strings.Index(response, "{")
		jsonEnd := strings.LastIndex(response, "}")
		if jsonStart >= 0 && jsonEnd > jsonStart {
			jsonStr := response[jsonStart : jsonEnd+1]
			if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
				return nil, fmt.Errorf("failed to parse LLM response: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to parse LLM response: %w", err)
		}
	}

	// Convert signal string to TradeSignal
	var signal TradeSignal
	switch result.Signal {
	case "BUY":
		signal = Buy
	case "SELL":
		signal = Sell
	default:
		signal = Hold
	}

	return &TradeDecision{
		Symbol:    stock.Symbol,
		Signal:    signal,
		Price:     stock.CurrentPrice,
		Timestamp: time.Now(),
		Rationale: result.Rationale,
		Score:     result.Confidence,
	}, nil
}

// callOpenAI calls the OpenAI API
func (l *LLMAdvisor) callOpenAI(prompt string) (string, error) {
	request := OpenAIRequest{
		Model: l.config.ModelName,
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a financial advisor specialized in stock trading. Analyze the provided market data and technical indicators to make a trading recommendation.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens:   l.config.MaxTokens,
		Temperature: l.config.Temperature,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", l.config.APIKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get response, status: %d, body: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var openAIResp OpenAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return openAIResp.Choices[0].Message.Content, nil
}

// callAnthropic calls the Anthropic API
func (l *LLMAdvisor) callAnthropic(prompt string) (string, error) {
	request := AnthropicRequest{
		Model:       l.config.ModelName,
		MaxTokens:   l.config.MaxTokens,
		Temperature: l.config.Temperature,
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", l.config.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get response, status: %d, body: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var anthropicResp AnthropicResponse
	if err := json.Unmarshal(body, &anthropicResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(anthropicResp.Content) == 0 {
		return "", fmt.Errorf("no content in response")
	}

	return anthropicResp.Content[0].Text, nil
}

// callDeepSeek calls the local DeepSeek model
func (l *LLMAdvisor) callDeepSeek(prompt string) (string, error) {
	if l.config.LocalPath == "" {
		return "", fmt.Errorf("local path for DeepSeek model not specified")
	}

	// Create a command to run the local DeepSeek model
	// This is a simplified example - actual implementation would depend on how DeepSeek is deployed
	cmd := exec.Command(l.config.LocalPath, "--prompt", prompt)
	
	var out bytes.Buffer
	cmd.Stdout = &out
	
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to run DeepSeek model: %w", err)
	}
	
	return out.String(), nil
}

// mockLLMResponse generates a mock response for testing
func (l *LLMAdvisor) mockLLMResponse(stock *data.Stock) (string, error) {
	// Simple logic to generate a mock trading signal based on price movement
	var signal string
	var rationale string
	var confidence float64

	// If price is higher than previous close, generate a BUY signal
	if stock.CurrentPrice > stock.PreviousClose {
		signal = "BUY"
		rationale = fmt.Sprintf("The stock price for %s has increased by %.2f%% from the previous close. Technical indicators suggest continued upward momentum.", 
			stock.Symbol, stock.ChangePercent)
		confidence = 0.7 + (stock.ChangePercent / 100)
		if confidence > 0.95 {
			confidence = 0.95
		}
	} else if stock.CurrentPrice < stock.PreviousClose {
		// If price is lower than previous close, generate a SELL signal
		signal = "SELL"
		rationale = fmt.Sprintf("The stock price for %s has decreased by %.2f%% from the previous close. Technical indicators suggest continued downward pressure.", 
			stock.Symbol, -stock.ChangePercent)
		confidence = 0.7 + (-stock.ChangePercent / 100)
		if confidence > 0.95 {
			confidence = 0.95
		}
	} else {
		// If price is the same as previous close, generate a HOLD signal
		signal = "HOLD"
		rationale = fmt.Sprintf("The stock price for %s has remained stable. No clear directional movement is detected in the technical indicators.", 
			stock.Symbol)
		confidence = 0.6
	}

	// Create a JSON response
	response := fmt.Sprintf(`{
  "signal": "%s",
  "rationale": "%s",
  "confidence": %.2f
}`, signal, rationale, confidence)

	return response, nil
}
