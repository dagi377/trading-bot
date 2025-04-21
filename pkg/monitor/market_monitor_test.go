package monitor

import (
	"testing"
	"time"

	"github.com/hustler/trading-bot/pkg/config"
	"github.com/hustler/trading-bot/pkg/data"
	"github.com/hustler/trading-bot/pkg/llm"
	"github.com/hustler/trading-bot/pkg/signal"
	"github.com/hustler/trading-bot/pkg/telegram"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementations
type MockDataProvider struct {
	mock.Mock
}

func (m *MockDataProvider) GetMarketData(symbol string) (*data.MarketData, error) {
	args := m.Called(symbol)
	return args.Get(0).(*data.MarketData), args.Error(1)
}

type MockSignalGenerator struct {
	mock.Mock
}

func (m *MockSignalGenerator) GenerateSignals(marketData map[string]signal.MarketData) ([]*signal.Signal, error) {
	args := m.Called(marketData)
	return args.Get(0).([]*signal.Signal), args.Error(1)
}

type MockLLMManager struct {
	mock.Mock
}

func (m *MockLLMManager) GenerateSignalExplanation(ctx context.Context, s *signal.Signal) (string, error) {
	args := m.Called(ctx, s)
	return args.String(0), args.Error(1)
}

type MockTelegramBot struct {
	mock.Mock
}

func (m *MockTelegramBot) SendSignal(s *signal.Signal) error {
	args := m.Called(s)
	return args.Error(0)
}

func TestNewMarketMonitor(t *testing.T) {
	// Create mocks
	cfg := config.CreateDefaultConfig()
	dataProvider := &MockDataProvider{}
	signalGen := &MockSignalGenerator{}
	llmManager := &MockLLMManager{}
	telegramBot := &MockTelegramBot{}

	// Create monitor
	monitor := NewMarketMonitor(cfg, dataProvider, signalGen, llmManager, telegramBot)

	// Verify monitor
	assert.NotNil(t, monitor)
	assert.Equal(t, cfg, monitor.config)
	assert.Equal(t, dataProvider, monitor.dataProvider)
	assert.Equal(t, signalGen, monitor.signalGen)
	assert.Equal(t, llmManager, monitor.llmManager)
	assert.Equal(t, telegramBot, monitor.telegramBot)
	assert.False(t, monitor.isRunning)
	assert.NotNil(t, monitor.stopChan)
	assert.Empty(t, monitor.signalHistory)
}

func TestStartStop(t *testing.T) {
	// Create mocks
	cfg := config.CreateDefaultConfig()
	dataProvider := &MockDataProvider{}
	signalGen := &MockSignalGenerator{}
	llmManager := &MockLLMManager{}
	telegramBot := &MockTelegramBot{}

	// Create monitor
	monitor := NewMarketMonitor(cfg, dataProvider, signalGen, llmManager, telegramBot)

	// Test start
	err := monitor.Start()
	assert.NoError(t, err)
	assert.True(t, monitor.IsRunning())

	// Test start when already running
	err = monitor.Start()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already running")

	// Test stop
	err = monitor.Stop()
	assert.NoError(t, err)
	assert.False(t, monitor.IsRunning())

	// Test stop when not running
	err = monitor.Stop()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not running")
}

func TestGetSignalHistory(t *testing.T) {
	// Create mocks
	cfg := config.CreateDefaultConfig()
	dataProvider := &MockDataProvider{}
	signalGen := &MockSignalGenerator{}
	llmManager := &MockLLMManager{}
	telegramBot := &MockTelegramBot{}

	// Create monitor
	monitor := NewMarketMonitor(cfg, dataProvider, signalGen, llmManager, telegramBot)

	// Add signals to history
	signal1 := &signal.Signal{ID: "1", Symbol: "AAPL"}
	signal2 := &signal.Signal{ID: "2", Symbol: "MSFT"}

	monitor.mu.Lock()
	monitor.signalHistory = append(monitor.signalHistory, signal1, signal2)
	monitor.mu.Unlock()

	// Get history
	history := monitor.GetSignalHistory()

	// Verify history
	assert.Len(t, history, 2)
	assert.Equal(t, "1", history[0].ID)
	assert.Equal(t, "2", history[1].ID)

	// Verify it's a copy (modify original)
	monitor.mu.Lock()
	monitor.signalHistory[0].Symbol = "CHANGED"
	monitor.mu.Unlock()

	// Get history again
	history = monitor.GetSignalHistory()

	// Verify history reflects changes
	assert.Equal(t, "CHANGED", history[0].Symbol)
}

func TestUpdateConfig(t *testing.T) {
	// Create mocks
	cfg := config.CreateDefaultConfig()
	dataProvider := &MockDataProvider{}
	signalGen := &MockSignalGenerator{}
	llmManager := &MockLLMManager{}
	telegramBot := &MockTelegramBot{}

	// Create monitor
	monitor := NewMarketMonitor(cfg, dataProvider, signalGen, llmManager, telegramBot)

	// Create new config
	newCfg := config.CreateDefaultConfig()
	newCfg.CheckInterval = 60
	newCfg.StockSymbols = []string{"AAPL", "MSFT", "GOOGL"}

	// Update config
	monitor.UpdateConfig(newCfg)

	// Verify config was updated
	assert.Equal(t, newCfg, monitor.config)
	assert.Equal(t, 60, monitor.config.CheckInterval)
	assert.Equal(t, []string{"AAPL", "MSFT", "GOOGL"}, monitor.config.StockSymbols)
}

func TestPerformMarketCheck(t *testing.T) {
	// Create mocks
	cfg := config.CreateDefaultConfig()
	cfg.StockSymbols = []string{"AAPL", "MSFT"}

	dataProvider := &MockDataProvider{}
	signalGen := &MockSignalGenerator{}
	llmManager := &MockLLMManager{}
	telegramBot := &MockTelegramBot{}

	// Create monitor
	monitor := NewMarketMonitor(cfg, dataProvider, signalGen, llmManager, telegramBot)

	// Create mock data
	appleData := &data.MarketData{
		Symbol:     "AAPL",
		Prices:     []float64{150.0, 151.0, 152.0},
		Volumes:    []float64{1000000, 1100000, 1200000},
		Timestamps: []time.Time{time.Now().Add(-2 * time.Hour), time.Now().Add(-1 * time.Hour), time.Now()},
	}

	msftData := &data.MarketData{
		Symbol:     "MSFT",
		Prices:     []float64{350.0, 351.0, 352.0},
		Volumes:    []float64{2000000, 2100000, 2200000},
		Timestamps: []time.Time{time.Now().Add(-2 * time.Hour), time.Now().Add(-1 * time.Hour), time.Now()},
	}

	// Create mock signals
	mockSignals := []*signal.Signal{
		{
			ID:          "SIG-AAPL-BUY-1",
			Symbol:      "AAPL",
			Type:        signal.BUY,
			Price:       152.0,
			TargetPrice: 155.0,
			StopLoss:    150.0,
			ExpectedROI: 1.97,
			Confidence:  0.85,
			GeneratedAt: time.Now(),
		},
	}

	// Set up mock expectations
	dataProvider.On("GetMarketData", "AAPL").Return(appleData, nil)
	dataProvider.On("GetMarketData", "MSFT").Return(msftData, nil)

	// The marketData parameter is complex, so we use mock.Anything
	signalGen.On("GenerateSignals", mock.Anything).Return(mockSignals, nil)

	// LLM should be called for the signal
	llmManager.On("GenerateSignalExplanation", mock.Anything, mockSignals[0]).Return("This is a test explanation", nil)

	// Telegram bot should be called to send the signal
	telegramBot.On("SendSignal", mockSignals[0]).Return(nil)

	// Perform market check
	err := monitor.performMarketCheck()

	// Verify no error
	assert.NoError(t, err)

	// Verify mocks were called
	dataProvider.AssertExpectations(t)
	signalGen.AssertExpectations(t)
	llmManager.AssertExpectations(t)
	telegramBot.AssertExpectations(t)

	// Verify signal was added to history
	history := monitor.GetSignalHistory()
	assert.Len(t, history, 1)
	assert.Equal(t, "SIG-AAPL-BUY-1", history[0].ID)
	assert.Equal(t, "This is a test explanation", history[0].Rationale)
}
