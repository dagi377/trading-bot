package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

// Config represents the application configuration
type Config struct {
	Admin          AdminConfig     `json:"admin"`
	Telegram       TelegramConfig  `json:"telegram"`
	DataSource     DataSourceConfig `json:"data_source"`
	LLM            LLMConfig       `json:"llm"`
	StockSymbols   []string        `json:"stock_symbols"`
	TradingHours   TradingHoursConfig `json:"trading_hours"`
	VolatilityParams VolatilityConfig `json:"volatility_params"`
	CheckInterval  int             `json:"check_interval"` // in seconds
	LogLevel       string          `json:"log_level"`
}

// AdminConfig represents admin-specific configuration
type AdminConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Port     int    `json:"port"`
}

// TelegramConfig represents Telegram-specific configuration
type TelegramConfig struct {
	BotToken     string  `json:"bot_token"`
	ChannelID    string  `json:"channel_id"`
	AdminUserIDs []int64 `json:"admin_user_ids"`
}

// DataSourceConfig represents data source configuration
type DataSourceConfig struct {
	Primary   string            `json:"primary"`
	Secondary string            `json:"secondary"`
	APIKeys   map[string]string `json:"api_keys"`
}

// LLMConfig represents LLM provider configuration
type LLMConfig struct {
	Provider   string `json:"provider"`
	APIKey     string `json:"api_key"`
	ModelName  string `json:"model_name"`
	LocalPath  string `json:"local_path"`
	MaxTokens  int    `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
}

// TradingHoursConfig represents trading hours configuration
type TradingHoursConfig struct {
	StartTime string `json:"start_time"` // Format: "HH:MM" in 24-hour format
	EndTime   string `json:"end_time"`   // Format: "HH:MM" in 24-hour format
	Start     string `json:"start"`      // Alias for StartTime for backward compatibility
	End       string `json:"end"`        // Alias for EndTime for backward compatibility
	TimeZone  string `json:"time_zone"`  // e.g., "America/New_York"
	Weekend   bool   `json:"weekend"`    // Whether to trade on weekends
}

// VolatilityConfig represents volatility detection parameters
type VolatilityConfig struct {
	MinVolatilityPercent float64 `json:"min_volatility_percent"`
	MinExpectedROI       float64 `json:"min_expected_roi"`
	StopLossPercent      float64 `json:"stop_loss_percent"`
	BollingerPeriod      int     `json:"bollinger_period"`
	BollingerDeviation   float64 `json:"bollinger_deviation"`
	RSIPeriod            int     `json:"rsi_period"`
	RSIOverbought        float64 `json:"rsi_overbought"`
	RSIOversold          float64 `json:"rsi_oversold"`
	VolumeThreshold      float64 `json:"volume_threshold"` // % above average
	ConfidenceThreshold  float64 `json:"confidence_threshold"`
}

// LoadConfigFromFile loads configuration from a file
func LoadConfigFromFile(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// SaveConfig saves configuration to a file
func SaveConfig(config *Config, path string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// SaveConfigToFile saves configuration to a file (alias for SaveConfig for backward compatibility)
func SaveConfigToFile(config *Config, path string) error {
	return SaveConfig(config, path)
}

// CreateDefaultConfig creates a default configuration
func CreateDefaultConfig() *Config {
	return &Config{
		Admin: AdminConfig{
			Username: "admin",
			Password: "hustler123",
			Port:     8080,
		},
		Telegram: TelegramConfig{
			BotToken:     "",
			ChannelID:    "",
			AdminUserIDs: []int64{},
		},
		DataSource: DataSourceConfig{
			Primary:   "yahoo",
			Secondary: "alphavantage",
			APIKeys: map[string]string{
				"alphavantage": "",
				"finnhub":      "",
			},
		},
		LLM: LLMConfig{
			Provider:    "openai",
			APIKey:      "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",
			ModelName:   "gpt-4",
			LocalPath:   "",
			MaxTokens:   1000,
			Temperature: 0.7,
		},
		StockSymbols: []string{"AAPL", "MSFT", "GOOGL", "AMZN", "META"},
		TradingHours: TradingHoursConfig{
			StartTime: "09:30",
			EndTime:   "15:30",
			Start:     "09:30", // For backward compatibility
			End:       "15:30", // For backward compatibility
			TimeZone:  "America/New_York",
			Weekend:   false,
		},
		VolatilityParams: VolatilityConfig{
			MinVolatilityPercent: 1.0,
			MinExpectedROI:       1.5,
			StopLossPercent:      0.5,
			BollingerPeriod:      20,
			BollingerDeviation:   2.0,
			RSIPeriod:            14,
			RSIOverbought:        70.0,
			RSIOversold:          30.0,
			VolumeThreshold:      150.0,
			ConfidenceThreshold:  0.7,
		},
		CheckInterval: 300, // 5 minutes
		LogLevel:      "info",
	}
}

// Variable for time.Now to allow mocking in tests
var timeNow = time.Now

// IsWithinTradingHours checks if the current time is within trading hours
func (c *Config) IsWithinTradingHours() (bool, error) {
	// Parse time zone
	loc, err := time.LoadLocation(c.TradingHours.TimeZone)
	if err != nil {
		return false, fmt.Errorf("invalid time zone: %w", err)
	}

	// Get current time in the specified time zone
	now := timeNow().In(loc)

	// Check if it's a weekend
	if !c.TradingHours.Weekend && (now.Weekday() == time.Saturday || now.Weekday() == time.Sunday) {
		return false, nil
	}

	// Parse start and end times
	var startHour, startMin, endHour, endMin int
	
	// Use Start/StartTime field (whichever is set)
	startTimeStr := c.TradingHours.StartTime
	if startTimeStr == "" {
		startTimeStr = c.TradingHours.Start
	}
	
	// Use End/EndTime field (whichever is set)
	endTimeStr := c.TradingHours.EndTime
	if endTimeStr == "" {
		endTimeStr = c.TradingHours.End
	}
	
	// Validate time format - only validate if we have values
	if startTimeStr != "" {
		_, err = time.Parse("15:04", startTimeStr)
		if err != nil {
			return false, fmt.Errorf("invalid start time format: %s", startTimeStr)
		}
	} else {
		return false, fmt.Errorf("missing start time")
	}
	
	if endTimeStr != "" {
		_, err = time.Parse("15:04", endTimeStr)
		if err != nil {
			return false, fmt.Errorf("invalid end time format: %s", endTimeStr)
		}
	} else {
		return false, fmt.Errorf("missing end time")
	}
	
	// Parse start time
	startParts := strings.Split(startTimeStr, ":")
	if len(startParts) != 2 {
		return false, fmt.Errorf("invalid start time format: %s", startTimeStr)
	}
	
	_, err = fmt.Sscanf(startParts[0], "%d", &startHour)
	if err != nil {
		return false, fmt.Errorf("invalid start hour: %w", err)
	}
	
	_, err = fmt.Sscanf(startParts[1], "%d", &startMin)
	if err != nil {
		return false, fmt.Errorf("invalid start minute: %w", err)
	}
	
	// Parse end time
	endParts := strings.Split(endTimeStr, ":")
	if len(endParts) != 2 {
		return false, fmt.Errorf("invalid end time format: %s", endTimeStr)
	}
	
	_, err = fmt.Sscanf(endParts[0], "%d", &endHour)
	if err != nil {
		return false, fmt.Errorf("invalid end hour: %w", err)
	}
	
	_, err = fmt.Sscanf(endParts[1], "%d", &endMin)
	if err != nil {
		return false, fmt.Errorf("invalid end minute: %w", err)
	}

	// Create time objects for start and end times
	startTime := time.Date(now.Year(), now.Month(), now.Day(), startHour, startMin, 0, 0, loc)
	endTime := time.Date(now.Year(), now.Month(), now.Day(), endHour, endMin, 0, 0, loc)

	// Check if current time is within trading hours
	return (now.Equal(startTime) || now.After(startTime)) && now.Before(endTime), nil
}

// ValidateConfig validates the configuration
func ValidateConfig(config *Config) error {
	// Validate trading hours
	startTimeStr := config.TradingHours.StartTime
	if startTimeStr == "" {
		startTimeStr = config.TradingHours.Start
	}
	
	endTimeStr := config.TradingHours.EndTime
	if endTimeStr == "" {
		endTimeStr = config.TradingHours.End
	}
	
	if _, err := time.Parse("15:04", startTimeStr); err != nil {
		return fmt.Errorf("invalid start time format: %w", err)
	}
	if _, err := time.Parse("15:04", endTimeStr); err != nil {
		return fmt.Errorf("invalid end time format: %w", err)
	}
	if _, err := time.LoadLocation(config.TradingHours.TimeZone); err != nil {
		return fmt.Errorf("invalid time zone: %w", err)
	}

	// Validate volatility parameters
	if config.VolatilityParams.MinVolatilityPercent <= 0 {
		return fmt.Errorf("min_volatility_percent must be positive")
	}
	if config.VolatilityParams.MinExpectedROI <= 0 {
		return fmt.Errorf("min_expected_roi must be positive")
	}
	if config.VolatilityParams.StopLossPercent <= 0 {
		return fmt.Errorf("stop_loss_percent must be positive")
	}
	if config.VolatilityParams.BollingerPeriod <= 0 {
		return fmt.Errorf("bollinger_period must be positive")
	}
	if config.VolatilityParams.RSIPeriod <= 0 {
		return fmt.Errorf("rsi_period must be positive")
	}

	// Validate check interval
	if config.CheckInterval <= 0 {
		return fmt.Errorf("check_interval must be positive")
	}

	return nil
}
