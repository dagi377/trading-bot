package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hustler/trading-bot/pkg/config"
	"github.com/hustler/trading-bot/pkg/data"
	"github.com/hustler/trading-bot/pkg/llm"
	"github.com/hustler/trading-bot/pkg/mock"
	"github.com/hustler/trading-bot/pkg/monitor"
	"github.com/hustler/trading-bot/pkg/performance"
	"github.com/hustler/trading-bot/pkg/signal"
	"github.com/hustler/trading-bot/pkg/telegram"
)

func main() {
	log.Println("Starting E2E Test for Hustler Trading Bot...")

	// Create default configuration
	cfg := config.CreateDefaultConfig()
	
	// Modify config for testing
	cfg.CheckInterval = 30 // 30 seconds for faster testing
	cfg.StockSymbols = []string{"AAPL", "MSFT", "GOOGL"}
	cfg.LLM.Provider = "mock" // Use mock LLM provider
	
	// Initialize components
	dataProvider := data.NewProvider(cfg)
	signalGen := signal.NewGenerator(cfg)
	perfMonitor := performance.NewMonitor()
	
	// Use mock Telegram bot
	telegramBot := mock.NewMockTelegramBot()
	
	// Initialize LLM manager
	llmManager, err := llm.NewManager(cfg.LLM)
	if err != nil {
		log.Fatalf("Failed to initialize LLM manager: %v", err)
	}
	
	// Initialize market monitor
	marketMonitor := monitor.NewMarketMonitor(
		cfg,
		dataProvider,
		signalGen,
		llmManager,
		telegramBot,
	)
	
	// Start market monitor
	err = marketMonitor.Start()
	if err != nil {
		log.Fatalf("Failed to start market monitor: %v", err)
	}
	log.Println("Market monitor started")
	
	// Run test for 2 minutes
	log.Println("Running E2E test for 2 minutes...")
	time.Sleep(2 * time.Minute)
	
	// Stop market monitor
	err = marketMonitor.Stop()
	if err != nil {
		log.Printf("Error stopping market monitor: %v", err)
	}
	
	// Get signal history
	signals := marketMonitor.GetSignalHistory()
	log.Printf("Generated %d signals during test", len(signals))
	
	// Get performance metrics
	metrics := perfMonitor.GetMetrics()
	log.Printf("Performance metrics: %d signals, %.2f%% success rate", 
		metrics.SignalsCount, metrics.SuccessRate)
	
	// Get Telegram messages
	messages := telegramBot.GetSentMessages()
	log.Printf("Sent %d messages to Telegram", len(messages))
	
	// Print summary
	fmt.Println("\n=== E2E Test Summary ===")
	fmt.Printf("Generated Signals: %d\n", len(signals))
	fmt.Printf("Telegram Messages: %d\n", len(messages))
	fmt.Println("Test completed successfully!")
	
	// Write results to file
	writeResultsToFile(signals, messages)
	
	log.Println("E2E Test completed")
}

func writeResultsToFile(signals []*signal.Signal, messages []string) {
	// Create results directory if it doesn't exist
	err := os.MkdirAll("./test_results", 0755)
	if err != nil {
		log.Printf("Error creating results directory: %v", err)
		return
	}
	
	// Write signals to file
	signalsFile, err := os.Create("./test_results/signals.txt")
	if err != nil {
		log.Printf("Error creating signals file: %v", err)
		return
	}
	defer signalsFile.Close()
	
	for i, s := range signals {
		signalsFile.WriteString(fmt.Sprintf("Signal %d: %s %s at $%.2f\n", 
			i+1, s.Type, s.Symbol, s.Price))
		signalsFile.WriteString(fmt.Sprintf("  Target: $%.2f, Stop Loss: $%.2f\n", 
			s.TargetPrice, s.StopLoss))
		signalsFile.WriteString(fmt.Sprintf("  Expected ROI: %.2f%%, Confidence: %.0f%%\n", 
			s.ExpectedROI, s.Confidence*100))
		signalsFile.WriteString(fmt.Sprintf("  Rationale: %s\n\n", s.Rationale))
	}
	
	// Write messages to file
	messagesFile, err := os.Create("./test_results/messages.txt")
	if err != nil {
		log.Printf("Error creating messages file: %v", err)
		return
	}
	defer messagesFile.Close()
	
	for i, msg := range messages {
		messagesFile.WriteString(fmt.Sprintf("Message %d:\n%s\n\n", i+1, msg))
	}
	
	log.Println("Test results written to ./test_results/")
}
