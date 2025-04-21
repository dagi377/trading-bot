package telegram

import (
	"testing"
	"time"

	"github.com/hustler/trading-bot/pkg/config"
	"github.com/hustler/trading-bot/pkg/signal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTelegramAPI is a mock implementation of the Telegram API client
type MockTelegramAPI struct {
	mock.Mock
}

func (m *MockTelegramAPI) SendMessage(chatID int64, text string, parseMode string) error {
	args := m.Called(chatID, text, parseMode)
	return args.Error(0)
}

func (m *MockTelegramAPI) GetUpdates(offset int, limit int) ([]Update, error) {
	args := m.Called(offset, limit)
	return args.Get(0).([]Update), args.Error(1)
}

func TestNewBot(t *testing.T) {
	// Create config
	cfg := &config.TelegramConfig{
		Token:   "test-token",
		Channel: "@test_channel",
	}

	// Create bot
	bot := NewBot(cfg)

	// Verify bot
	assert.NotNil(t, bot)
	assert.Equal(t, cfg, bot.config)
	assert.NotNil(t, bot.api)
	assert.NotNil(t, bot.subscribers)
	assert.Equal(t, 0, len(bot.subscribers))
}

func TestSendSignal(t *testing.T) {
	// Create mock API
	mockAPI := new(MockTelegramAPI)

	// Create config
	cfg := &config.TelegramConfig{
		Token:   "test-token",
		Channel: "@test_channel",
	}

	// Create bot with mock API
	bot := NewBot(cfg)
	bot.api = mockAPI

	// Create test signal
	testSignal := &signal.Signal{
		ID:          "SIG-AAPL-BUY-1",
		Symbol:      "AAPL",
		Type:        signal.BUY,
		Price:       150.0,
		TargetPrice: 155.0,
		StopLoss:    148.0,
		ExpectedROI: 3.33,
		Confidence:  0.85,
		Rationale:   "This is a test rationale",
		GeneratedAt: time.Now(),
		TimeFrame:   "1-3 hours",
	}

	// Set up mock expectations
	// The message text is complex, so we use mock.Anything
	mockAPI.On("SendMessage", int64(0), mock.Anything, "HTML").Return(nil)

	// Send signal
	err := bot.SendSignal(testSignal)

	// Verify no error
	assert.NoError(t, err)

	// Verify mock was called
	mockAPI.AssertExpectations(t)
}

func TestSendSignalToSubscribers(t *testing.T) {
	// Create mock API
	mockAPI := new(MockTelegramAPI)

	// Create config
	cfg := &config.TelegramConfig{
		Token:   "test-token",
		Channel: "@test_channel",
	}

	// Create bot with mock API
	bot := NewBot(cfg)
	bot.api = mockAPI

	// Add subscribers
	bot.subscribers = []int64{123456789, 987654321}

	// Create test signal
	testSignal := &signal.Signal{
		ID:          "SIG-AAPL-BUY-1",
		Symbol:      "AAPL",
		Type:        signal.BUY,
		Price:       150.0,
		TargetPrice: 155.0,
		StopLoss:    148.0,
		ExpectedROI: 3.33,
		Confidence:  0.85,
		Rationale:   "This is a test rationale",
		GeneratedAt: time.Now(),
		TimeFrame:   "1-3 hours",
	}

	// Set up mock expectations
	// The message text is complex, so we use mock.Anything
	mockAPI.On("SendMessage", int64(0), mock.Anything, "HTML").Return(nil)
	mockAPI.On("SendMessage", int64(123456789), mock.Anything, "HTML").Return(nil)
	mockAPI.On("SendMessage", int64(987654321), mock.Anything, "HTML").Return(nil)

	// Send signal
	err := bot.SendSignal(testSignal)

	// Verify no error
	assert.NoError(t, err)

	// Verify mock was called
	mockAPI.AssertExpectations(t)
}

func TestProcessUpdates(t *testing.T) {
	// Create mock API
	mockAPI := new(MockTelegramAPI)

	// Create config
	cfg := &config.TelegramConfig{
		Token:   "test-token",
		Channel: "@test_channel",
	}

	// Create bot with mock API
	bot := NewBot(cfg)
	bot.api = mockAPI

	// Create test updates
	updates := []Update{
		{
			UpdateID: 1,
			Message: Message{
				MessageID: 1,
				From: User{
					ID:        123456789,
					FirstName: "Test",
					LastName:  "User",
					Username:  "testuser",
				},
				Chat: Chat{
					ID:   123456789,
					Type: "private",
				},
				Text: "/start",
				Date: int(time.Now().Unix()),
			},
		},
		{
			UpdateID: 2,
			Message: Message{
				MessageID: 2,
				From: User{
					ID:        987654321,
					FirstName: "Another",
					LastName:  "User",
					Username:  "anotheruser",
				},
				Chat: Chat{
					ID:   987654321,
					Type: "private",
				},
				Text: "/help",
				Date: int(time.Now().Unix()),
			},
		},
	}

	// Set up mock expectations
	mockAPI.On("GetUpdates", 0, 100).Return(updates, nil)
	mockAPI.On("SendMessage", int64(123456789), mock.Anything, "HTML").Return(nil)
	mockAPI.On("SendMessage", int64(987654321), mock.Anything, "HTML").Return(nil)

	// Process updates
	err := bot.ProcessUpdates()

	// Verify no error
	assert.NoError(t, err)

	// Verify subscribers were added
	assert.Contains(t, bot.subscribers, int64(123456789))
	assert.NotContains(t, bot.subscribers, int64(987654321)) // /help doesn't add subscriber

	// Verify mock was called
	mockAPI.AssertExpectations(t)
}

func TestHandleCommand(t *testing.T) {
	// Create mock API
	mockAPI := new(MockTelegramAPI)

	// Create config
	cfg := &config.TelegramConfig{
		Token:   "test-token",
		Channel: "@test_channel",
	}

	// Create bot with mock API
	bot := NewBot(cfg)
	bot.api = mockAPI

	// Test cases
	testCases := []struct {
		command    string
		chatID     int64
		shouldSend bool
		shouldAdd  bool
	}{
		{"/start", 123456789, true, true},
		{"/help", 987654321, true, false},
		{"/stop", 123456789, true, false},
		{"/unknown", 555555555, false, false},
	}

	for _, tc := range testCases {
		t.Run(tc.command, func(t *testing.T) {
			// Reset subscribers
			bot.subscribers = []int64{}

			// Set up mock expectations if needed
			if tc.shouldSend {
				mockAPI.On("SendMessage", tc.chatID, mock.Anything, "HTML").Return(nil).Once()
			}

			// Handle command
			bot.handleCommand(tc.command, tc.chatID)

			// Verify subscriber was added if expected
			if tc.shouldAdd {
				assert.Contains(t, bot.subscribers, tc.chatID)
			} else {
				assert.NotContains(t, bot.subscribers, tc.chatID)
			}

			// Verify mock was called
			mockAPI.AssertExpectations(t)
		})
	}
}

func TestFormatSignalMessage(t *testing.T) {
	// Create test signal
	testSignal := &signal.Signal{
		ID:          "SIG-AAPL-BUY-1",
		Symbol:      "AAPL",
		Type:        signal.BUY,
		Price:       150.0,
		TargetPrice: 155.0,
		StopLoss:    148.0,
		ExpectedROI: 3.33,
		Confidence:  0.85,
		Rationale:   "This is a test rationale",
		GeneratedAt: time.Date(2025, 4, 20, 10, 15, 0, 0, time.UTC),
		TimeFrame:   "1-3 hours",
	}

	// Format message
	message := formatSignalMessage(testSignal)

	// Verify message contains key information
	assert.Contains(t, message, "BUY SIGNAL: AAPL")
	assert.Contains(t, message, "Entry Price: $150.00")
	assert.Contains(t, message, "Target Price: $155.00")
	assert.Contains(t, message, "Stop Loss: $148.00")
	assert.Contains(t, message, "Expected ROI: +3.33%")
	assert.Contains(t, message, "Confidence: 85%")
	assert.Contains(t, message, "This is a test rationale")
	assert.Contains(t, message, "2025-04-20 10:15:00")

	// Test SELL signal
	testSignal.Type = signal.SELL
	message = formatSignalMessage(testSignal)
	assert.Contains(t, message, "SELL SIGNAL: AAPL")
	assert.Contains(t, message, "Expected ROI: -3.33%")
}
