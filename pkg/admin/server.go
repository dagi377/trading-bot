package admin

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/hustler/trading-bot/pkg/config"
)

// Server represents the admin web interface server
type Server struct {
	config     *config.Config
	configPath string
	templates  *template.Template
	mu         sync.RWMutex
}

// NewServer creates a new admin server
func NewServer(cfg *config.Config, configPath string, templatesDir string) (*Server, error) {
	// Load templates
	templates, err := template.ParseGlob(filepath.Join(templatesDir, "*.html"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	return &Server{
		config:     cfg,
		configPath: configPath,
		templates:  templates,
		mu:         sync.RWMutex{},
	}, nil
}

// Start starts the admin server
func (s *Server) Start() error {
	// Set up routes
	http.HandleFunc("/", s.authMiddleware(s.handleDashboard))
	http.HandleFunc("/login", s.handleLogin)
	http.HandleFunc("/logout", s.handleLogout)
	http.HandleFunc("/stocks", s.authMiddleware(s.handleStocks))
	http.HandleFunc("/settings", s.authMiddleware(s.handleSettings))
	http.HandleFunc("/api/config", s.authMiddleware(s.handleAPIConfig))
	http.HandleFunc("/api/stocks", s.authMiddleware(s.handleAPIStocks))
	http.HandleFunc("/api/signals", s.authMiddleware(s.handleAPISignals))
	http.HandleFunc("/api/performance", s.authMiddleware(s.handleAPIPerformance))

	// Serve static files
	fs := http.FileServer(http.Dir(filepath.Join(templatesDir, "static")))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Start server
	addr := fmt.Sprintf(":%d", s.config.Admin.Port)
	log.Printf("Starting admin server on %s", addr)
	return http.ListenAndServe(addr, nil)
}

// authMiddleware checks if the user is authenticated
func (s *Server) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if user is authenticated
		cookie, err := r.Cookie("auth")
		if err != nil || cookie.Value != "authenticated" {
			// Redirect to login page
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// User is authenticated, proceed to next handler
		next(w, r)
	}
}

// handleLogin handles the login page
func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse form
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		// Check credentials
		username := r.FormValue("username")
		password := r.FormValue("password")

		s.mu.RLock()
		validUsername := s.config.Admin.Username
		validPassword := s.config.Admin.Password
		s.mu.RUnlock()

		if username == validUsername && password == validPassword {
			// Set authentication cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "auth",
				Value:    "authenticated",
				Path:     "/",
				HttpOnly: true,
			})

			// Redirect to dashboard
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Invalid credentials
		s.templates.ExecuteTemplate(w, "login.html", map[string]interface{}{
			"Error": "Invalid username or password",
		})
		return
	}

	// Show login page
	s.templates.ExecuteTemplate(w, "login.html", nil)
}

// handleLogout handles the logout request
func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	// Clear authentication cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	// Redirect to login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// handleDashboard handles the dashboard page
func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	cfg := s.config
	s.mu.RUnlock()

	// Render dashboard template
	s.templates.ExecuteTemplate(w, "dashboard.html", map[string]interface{}{
		"Config": cfg,
		"Active": "dashboard",
	})
}

// handleStocks handles the stocks management page
func (s *Server) handleStocks(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	cfg := s.config
	s.mu.RUnlock()

	// Render stocks template
	s.templates.ExecuteTemplate(w, "stocks.html", map[string]interface{}{
		"Config": cfg,
		"Active": "stocks",
	})
}

// handleSettings handles the settings page
func (s *Server) handleSettings(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	cfg := s.config
	s.mu.RUnlock()

	// Render settings template
	s.templates.ExecuteTemplate(w, "settings.html", map[string]interface{}{
		"Config": cfg,
		"Active": "settings",
	})
}

// handleAPIConfig handles the API endpoint for configuration
func (s *Server) handleAPIConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		// Return current configuration
		s.mu.RLock()
		cfg := s.config
		s.mu.RUnlock()

		json.NewEncoder(w).Encode(cfg)
		return
	}

	if r.Method == http.MethodPOST {
		// Update configuration
		var newConfig config.Config
		err := json.NewDecoder(r.Body).Decode(&newConfig)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse request body: %v", err), http.StatusBadRequest)
			return
		}

		// Validate configuration
		err = config.ValidateConfig(&newConfig)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid configuration: %v", err), http.StatusBadRequest)
			return
		}

		// Update configuration
		s.mu.Lock()
		s.config = &newConfig
		s.mu.Unlock()

		// Save configuration to file
		err = config.SaveConfig(&newConfig, s.configPath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to save configuration: %v", err), http.StatusInternalServerError)
			return
		}

		// Return success
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
		return
	}

	// Method not allowed
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// handleAPIStocks handles the API endpoint for stocks
func (s *Server) handleAPIStocks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		// Return current stocks
		s.mu.RLock()
		stocks := s.config.StockSymbols
		s.mu.RUnlock()

		json.NewEncoder(w).Encode(stocks)
		return
	}

	if r.Method == http.MethodPOST {
		// Update stocks
		var stocks []string
		err := json.NewDecoder(r.Body).Decode(&stocks)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse request body: %v", err), http.StatusBadRequest)
			return
		}

		// Validate stocks
		for i, symbol := range stocks {
			stocks[i] = strings.ToUpper(strings.TrimSpace(symbol))
			if stocks[i] == "" {
				http.Error(w, "Stock symbols cannot be empty", http.StatusBadRequest)
				return
			}
		}

		// Update configuration
		s.mu.Lock()
		s.config.StockSymbols = stocks
		s.mu.Unlock()

		// Save configuration to file
		err = config.SaveConfig(s.config, s.configPath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to save configuration: %v", err), http.StatusInternalServerError)
			return
		}

		// Return success
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
		return
	}

	// Method not allowed
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// handleAPISignals handles the API endpoint for signals
func (s *Server) handleAPISignals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Mock signals data for now
	signals := []map[string]interface{}{
		{
			"id":           "sig-001",
			"symbol":       "AAPL",
			"type":         "BUY",
			"price":        175.50,
			"target_price": 180.25,
			"stop_loss":    173.00,
			"roi":          2.7,
			"confidence":   0.85,
			"timestamp":    "2025-04-20T10:15:30Z",
			"status":       "ACTIVE",
		},
		{
			"id":           "sig-002",
			"symbol":       "MSFT",
			"type":         "SELL",
			"price":        350.75,
			"target_price": 345.00,
			"stop_loss":    353.50,
			"roi":          1.6,
			"confidence":   0.75,
			"timestamp":    "2025-04-20T09:45:12Z",
			"status":       "ACTIVE",
		},
	}

	json.NewEncoder(w).Encode(signals)
}

// handleAPIPerformance handles the API endpoint for performance metrics
func (s *Server) handleAPIPerformance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Mock performance data for now
	performance := map[string]interface{}{
		"signals_count": 32,
		"success_rate":  68.5,
		"average_roi":   1.8,
		"total_profit":  12.5,
		"by_symbol": []map[string]interface{}{
			{
				"symbol":       "AAPL",
				"signals":      8,
				"success_rate": 75.0,
				"average_roi":  2.1,
			},
			{
				"symbol":       "MSFT",
				"signals":      6,
				"success_rate": 66.7,
				"average_roi":  1.5,
			},
			{
				"symbol":       "GOOGL",
				"signals":      5,
				"success_rate": 60.0,
				"average_roi":  1.9,
			},
		},
		"by_day": []map[string]interface{}{
			{
				"date":         "2025-04-19",
				"signals":      12,
				"success_rate": 75.0,
				"profit":       5.2,
			},
			{
				"date":         "2025-04-18",
				"signals":      10,
				"success_rate": 60.0,
				"profit":       3.8,
			},
			{
				"date":         "2025-04-17",
				"signals":      10,
				"success_rate": 70.0,
				"profit":       3.5,
			},
		},
	}

	json.NewEncoder(w).Encode(performance)
}

// UpdateConfig updates the server configuration
func (s *Server) UpdateConfig(cfg *config.Config) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.config = cfg
}

// GetConfig returns the current server configuration
func (s *Server) GetConfig() *config.Config {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.config
}
