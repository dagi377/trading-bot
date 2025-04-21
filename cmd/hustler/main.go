package main

import (
	"log"
	"os"
	ossignal "os/signal"
	"syscall"
	"time"

	"github.com/hustler/trading-bot/pkg/config"
	"github.com/hustler/trading-bot/pkg/data"
	"github.com/hustler/trading-bot/pkg/llm"
	"github.com/hustler/trading-bot/pkg/monitor"
	"github.com/hustler/trading-bot/pkg/signal"
	"github.com/hustler/trading-bot/pkg/telegram"
)

func main() {
	log.Println("Starting Hustler Trading Bot...")

	// Load configuration
	cfg := config.CreateDefaultConfig()
	if len(os.Args) > 1 {
		configFile := os.Args[1]
		loadedCfg, err := config.LoadConfigFromFile(configFile)
		if err != nil {
			log.Printf("Warning: Failed to load config from %s: %v", configFile, err)
			log.Println("Using default configuration")
		} else {
			cfg = loadedCfg
			log.Printf("Loaded configuration from %s", configFile)
		}
	} else {
		log.Println("No config file specified, using default configuration")
	}

	// Initialize components
	dataProvider := data.NewProvider(cfg)
	signalGen := signal.NewGenerator(cfg)
	telegramBot := telegram.NewBot(cfg.Telegram)

	// Initialize LLM manager
	llmManager, err := llm.NewManager(&cfg.LLM)
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

	// Start processing Telegram updates in a separate goroutine
	go func() {
		for {
			err := telegramBot.ProcessUpdates()
			if err != nil {
				log.Printf("Error processing Telegram updates: %v", err)
			}
			time.Sleep(5 * time.Second)
		}
	}()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	ossignal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for termination signal
	sig := <-sigChan
	log.Printf("Received signal %v, shutting down...", sig)

	// Stop market monitor
	err = marketMonitor.Stop()
	if err != nil {
		log.Printf("Error stopping market monitor: %v", err)
	}

	log.Println("Hustler Trading Bot shutdown complete")
}
