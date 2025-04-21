package ui

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/hustler/trading-bot/pkg/config"
	"github.com/hustler/trading-bot/pkg/data"
	"github.com/hustler/trading-bot/pkg/llm"
	"github.com/hustler/trading-bot/pkg/news"
	"github.com/hustler/trading-bot/pkg/signal"
	"github.com/hustler/trading-bot/pkg/telegram"
)

// Controller handles the web UI and API endpoints
type Controller struct {
	config        *config.AppConfig
	marketWatcher *data.MarketWatcher
	newsMonitor   *news.Monitor
	llmManager    *llm.Manager
	signalGen     *signal.Generator
	telegramBot   *telegram.Bot
}

// NewController creates a new UI controller
func NewController(
	config *config.AppConfig,
	marketWatcher *data.MarketWatcher,
	newsMonitor *news.Monitor,
	llmManager *llm.Manager,
	signalGen *signal.Generator,
	telegramBot *telegram.Bot,
) *Controller {
	return &Controller{
		config:        config,
		marketWatcher: marketWatcher,
		newsMonitor:   newsMonitor,
		llmManager:    llmManager,
		signalGen:     signalGen,
		telegramBot:   telegramBot,
	}
}

// Start starts the web server
func (c *Controller) Start(port int) error {
	// Set up API routes
	http.HandleFunc("/api/stocks", c.handleStocks)
	http.HandleFunc("/api/stock", c.handleStock)
	http.HandleFunc("/api/signals", c.handleSignals)
	http.HandleFunc("/api/signal", c.handleSignal)
	http.HandleFunc("/api/news", c.handleNews)
	http.HandleFunc("/api/config", c.handleConfig)
	http.HandleFunc("/api/telegram/test", c.handleTelegramTest)
	http.HandleFunc("/api/llm/switch", c.handleLLMSwitch)
	http.HandleFunc("/api/generate-signals", c.handleGenerateSignals)

	// Serve static files
	http.Handle("/", http.FileServer(http.Dir("./web/admin")))

	// Start the server
	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting admin UI server on %s", addr)
	return http.ListenAndServe(addr, nil)
}

// handleStocks handles requests for all stocks
func (c *Controller) handleStocks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stocks := c.marketWatcher.GetAllStocks()
	writeJSON(w, stocks)
}

// handleStock handles requests for a specific stock
func (c *Controller) handleStock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		http.Error(w, "Symbol parameter is required", http.StatusBadRequest)
		return
	}

	stock, exists := c.marketWatcher.GetStock(symbol)
	if !exists {
		http.Error(w, "Stock not found", http.StatusNotFound)
		return
	}

	writeJSON(w, stock)
}

// handleSignals handles requests for all signals
func (c *Controller) handleSignals(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	signals := c.signalGen.GetAllSignals()
	writeJSON(w, signals)
}

// handleSignal handles requests for a specific signal
func (c *Controller) handleSignal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		http.Error(w, "Symbol parameter is required", http.StatusBadRequest)
		return
	}

	signal, exists := c.signalGen.GetSignal(symbol)
	if !exists {
		http.Error(w, "Signal not found", http.StatusNotFound)
		return
	}

	writeJSON(w, signal)
}

// handleNews handles requests for news articles
func (c *Controller) handleNews(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	symbol := r.URL.Query().Get("symbol")
	limitStr := r.URL.Query().Get("limit")

	limit := 10 // Default limit
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
	}

	var articles []news.Article
	if symbol != "" {
		articles = c.newsMonitor.GetArticlesForSymbol(symbol, limit)
	} else {
		articles = c.newsMonitor.GetLatestArticles(limit)
	}

	writeJSON(w, articles)
}

// handleConfig handles requests for configuration
func (c *Controller) handleConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Return current configuration
		writeJSON(w, c.config)
		return
	}

	if r.Method == http.MethodPost {
		// Update configuration
		var newConfig config.AppConfig
		if err := json.NewDecoder(r.Body).Decode(&newConfig); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Save the new configuration
		if err := config.SaveConfig(&newConfig, "config.json"); err != nil {
			http.Error(w, fmt.Sprintf("Failed to save configuration: %v", err), http.StatusInternalServerError)
			return
		}

		// Update the current configuration
		*c.config = newConfig

		// Apply configuration changes
		// This would need to be implemented based on the specific components
		// that need to be updated when configuration changes

		w.WriteHeader(http.StatusOK)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// handleTelegramTest handles requests to test the Telegram bot
func (c *Controller) handleTelegramTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	message := r.URL.Query().Get("message")
	if message == "" {
		message = "This is a test message from the Hustler Trading Bot."
	}

	if err := c.telegramBot.SendMessage(message); err != nil {
		http.Error(w, fmt.Sprintf("Failed to send message: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Test message sent successfully"))
}

// handleLLMSwitch handles requests to switch the LLM provider
func (c *Controller) handleLLMSwitch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	provider := r.URL.Query().Get("provider")
	if provider == "" {
		http.Error(w, "Provider parameter is required", http.StatusBadRequest)
		return
	}

	// Get the indicator processor from the LLM advisor
	// This is a simplified approach - in a real implementation, you would need
	// to properly handle the dependency injection
	indicatorProcessor := c.llmManager.GetAdvisor().GetIndicatorProcessor()

	if err := c.llmManager.SwitchProvider(provider, indicatorProcessor); err != nil {
		http.Error(w, fmt.Sprintf("Failed to switch LLM provider: %v", err), http.StatusInternalServerError)
		return
	}

	// Update the configuration
	c.config.LLM.Provider = provider
	if err := config.SaveConfig(c.config, "config.json"); err != nil {
		log.Printf("Warning: Failed to save configuration after LLM switch: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("LLM provider switched to %s", provider)))
}

// handleGenerateSignals handles requests to generate signals
func (c *Controller) handleGenerateSignals(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	symbol := r.URL.Query().Get("symbol")

	var signals []*signal.Signal
	var err error

	if symbol != "" {
		// Generate signal for a specific symbol
		var s *signal.Signal
		s, err = c.signalGen.GenerateSignal(symbol)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to generate signal: %v", err), http.StatusInternalServerError)
			return
		}
		signals = []*signal.Signal{s}
	} else {
		// Generate signals for all stocks
		signals = c.signalGen.GenerateAllSignals()
	}

	// Send signals via Telegram if requested
	sendTelegram := r.URL.Query().Get("telegram") == "true"
	if sendTelegram && len(signals) > 0 {
		for _, s := range signals {
			message := signal.FormatSignalMessage(s)
			if err := c.telegramBot.SendMessage(message); err != nil {
				log.Printf("Warning: Failed to send signal via Telegram: %v", err)
			}
			// Add a small delay between messages to avoid rate limiting
			time.Sleep(100 * time.Millisecond)
		}
	}

	writeJSON(w, signals)
}

// Helper function to write JSON response
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
	}
}
