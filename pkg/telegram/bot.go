package telegram

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/hustler/trading-bot/pkg/config"
	"github.com/hustler/trading-bot/pkg/signal"
)

// Bot represents a Telegram bot for sending trading signals
type Bot struct {
	config      config.TelegramConfig
	mockMode    bool
	mockMessages []string
	subscribers  map[int64]bool
	adminUsers   map[int64]bool
	mu           sync.RWMutex
}

// NewBot creates a new Telegram bot
func NewBot(config config.TelegramConfig) *Bot {
	return NewBotWithMode(config, false)
}

// NewBotWithMode creates a new Telegram bot with specified mock mode
func NewBotWithMode(config config.TelegramConfig, mockMode bool) *Bot {
	adminUsers := make(map[int64]bool)
	for _, id := range config.AdminUserIDs {
		adminUsers[id] = true
	}

	return &Bot{
		config:      config,
		mockMode:    mockMode,
		mockMessages: []string{},
		subscribers:  make(map[int64]bool),
		adminUsers:   adminUsers,
		mu:           sync.RWMutex{},
	}
}

// SendMessage sends a message to the configured Telegram channel
func (b *Bot) SendMessage(message string) error {
	if b.mockMode {
		b.mu.Lock()
		b.mockMessages = append(b.mockMessages, message)
		b.mu.Unlock()
		log.Printf("[MOCK] Telegram message sent: %s", message)
		return nil
	}

	// In a real implementation, this would use the Telegram Bot API
	// to send the message to the configured channel
	log.Printf("Would send to Telegram: %s", message)
	
	// TODO: Implement actual Telegram API call
	// Example:
	// url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", b.config.BotToken)
	// payload := map[string]interface{}{
	//     "chat_id": b.config.ChannelID,
	//     "text": message,
	//     "parse_mode": "HTML",
	// }
	// ... HTTP POST request with payload

	return nil
}

// SendSignal formats and sends a trading signal via Telegram
func (b *Bot) SendSignal(s *signal.Signal) error {
	message := signal.FormatSignalMessage(s)
	return b.SendMessage(message)
}

// HandleCommand processes a command from a user
func (b *Bot) HandleCommand(userID int64, command string, args []string) (string, error) {
	command = strings.ToLower(command)
	
	switch command {
	case "/start":
		return b.handleStartCommand(userID)
	case "/settings":
		return b.handleSettingsCommand(userID, args)
	case "/performance":
		return b.handlePerformanceCommand(userID)
	case "/help":
		return b.handleHelpCommand(userID)
	default:
		return "Unknown command. Type /help for available commands.", nil
	}
}

// handleStartCommand handles the /start command
func (b *Bot) handleStartCommand(userID int64) (string, error) {
	b.mu.Lock()
	b.subscribers[userID] = true
	b.mu.Unlock()
	
	return "Welcome to Hustler Trading Bot! You are now subscribed to trading signals.\n\n" +
		"You will receive intraday trading signals based on volatility patterns.\n\n" +
		"Type /help to see available commands.", nil
}

// handleSettingsCommand handles the /settings command
func (b *Bot) handleSettingsCommand(userID int64, args []string) (string, error) {
	// In a real implementation, this would allow users to configure their preferences
	return "Settings functionality will be available soon.", nil
}

// handlePerformanceCommand handles the /performance command
func (b *Bot) handlePerformanceCommand(userID int64) (string, error) {
	// In a real implementation, this would return performance statistics
	return fmt.Sprintf("Performance Statistics (Last 7 Days):\n\n" +
		"Signals Sent: 32\n" +
		"Success Rate: 68%%\n" +
		"Average ROI: 1.2%%\n" +
		"Best Signal: AAPL +3.5%%\n" +
		"Last Updated: %s", time.Now().Format("2006-01-02 15:04:05")), nil
}

// handleHelpCommand handles the /help command
func (b *Bot) handleHelpCommand(userID int64) (string, error) {
	return "Available Commands:\n\n" +
		"/start - Subscribe to trading signals\n" +
		"/settings - Configure your preferences\n" +
		"/performance - View bot performance statistics\n" +
		"/help - Show this help message", nil
}

// IsAdmin checks if a user is an admin
func (b *Bot) IsAdmin(userID int64) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.adminUsers[userID]
}

// GetSubscribers returns the list of subscriber IDs
func (b *Bot) GetSubscribers() []int64 {
	b.mu.RLock()
	defer b.mu.RUnlock()
	
	subscribers := make([]int64, 0, len(b.subscribers))
	for id := range b.subscribers {
		subscribers = append(subscribers, id)
	}
	
	return subscribers
}

// GetMockMessages returns the list of mock messages (for testing)
func (b *Bot) GetMockMessages() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	
	messages := make([]string, len(b.mockMessages))
	copy(messages, b.mockMessages)
	
	return messages
}

// UpdateConfig updates the bot configuration
func (b *Bot) UpdateConfig(config config.TelegramConfig) {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	b.config = config
	
	// Update admin users
	b.adminUsers = make(map[int64]bool)
	for _, id := range config.AdminUserIDs {
		b.adminUsers[id] = true
	}
}

// ProcessUpdates processes incoming updates from Telegram
func (b *Bot) ProcessUpdates() error {
	// In a real implementation, this would poll the Telegram API for updates
	// and process incoming messages
	
	// For now, we'll just return nil
	return nil
}
