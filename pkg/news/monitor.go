package news

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/hustler/trading-bot/pkg/auth"
	"github.com/hustler/trading-bot/pkg/config"
)

// Article represents a financial news article
type Article struct {
	Title       string
	Description string
	URL         string
	Source      string
	PublishedAt time.Time
	Sentiment   float64 // -1.0 to 1.0 (negative to positive)
	Symbols     []string
	Keywords    []string
}

// Monitor watches for financial news from various sources
type Monitor struct {
	config      config.NewsConfig
	authManager *auth.AuthManager
	articles    []Article
	mu          sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
	callbacks   []func([]Article)
}

// NewMonitor creates a new news monitor
func NewMonitor(cfg config.NewsConfig, authManager *auth.AuthManager) *Monitor {
	ctx, cancel := context.WithCancel(context.Background())
	return &Monitor{
		config:      cfg,
		authManager: authManager,
		articles:    make([]Article, 0),
		ctx:         ctx,
		cancel:      cancel,
		callbacks:   make([]func([]Article), 0),
	}
}

// Start begins monitoring for news
func (m *Monitor) Start() {
	go func() {
		// Initial fetch
		m.fetchAllNews()

		// Set up ticker for periodic fetching
		ticker := time.NewTicker(time.Duration(m.config.PollInterval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				m.fetchAllNews()
			case <-m.ctx.Done():
				return
			}
		}
	}()
}

// Stop stops monitoring for news
func (m *Monitor) Stop() {
	m.cancel()
}

// GetLatestArticles returns the latest news articles
func (m *Monitor) GetLatestArticles(limit int) []Article {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if limit <= 0 || limit > len(m.articles) {
		limit = len(m.articles)
	}

	result := make([]Article, limit)
	copy(result, m.articles[:limit])
	return result
}

// GetArticlesForSymbol returns news articles for a specific stock symbol
func (m *Monitor) GetArticlesForSymbol(symbol string, limit int) []Article {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]Article, 0)
	for _, article := range m.articles {
		for _, s := range article.Symbols {
			if strings.EqualFold(s, symbol) {
				result = append(result, article)
				break
			}
		}

		if limit > 0 && len(result) >= limit {
			break
		}
	}

	return result
}

// RegisterCallback registers a callback function to be called when new articles are fetched
func (m *Monitor) RegisterCallback(callback func([]Article)) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.callbacks = append(m.callbacks, callback)
}

// fetchAllNews fetches news from all configured sources
func (m *Monitor) fetchAllNews() {
	var newArticles []Article

	for _, source := range m.config.Sources {
		var articles []Article
		var err error

		switch source {
		case "marketaux":
			articles, err = m.fetchMarketauxNews()
		case "twitter":
			articles, err = m.fetchTwitterNews()
		default:
			log.Printf("Unsupported news source: %s", source)
			continue
		}

		if err != nil {
			log.Printf("Error fetching news from %s: %v", source, err)
			continue
		}

		newArticles = append(newArticles, articles...)
	}

	if len(newArticles) > 0 {
		m.updateArticles(newArticles)
	}
}

// updateArticles updates the articles list with new articles
func (m *Monitor) updateArticles(newArticles []Article) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Add new articles to the beginning of the list
	combined := append(newArticles, m.articles...)

	// Remove duplicates (based on URL)
	seen := make(map[string]bool)
	unique := make([]Article, 0, len(combined))

	for _, article := range combined {
		if !seen[article.URL] {
			seen[article.URL] = true
			unique = append(unique, article)
		}
	}

	// Limit to a reasonable number to prevent memory issues
	const maxArticles = 1000
	if len(unique) > maxArticles {
		unique = unique[:maxArticles]
	}

	m.articles = unique

	// Notify callbacks
	callbacks := make([]func([]Article), len(m.callbacks))
	copy(callbacks, m.callbacks)

	// Call callbacks outside of the lock
	go func(articles []Article, callbacks []func([]Article)) {
		for _, callback := range callbacks {
			callback(articles)
		}
	}(newArticles, callbacks)
}

// fetchMarketauxNews fetches news from Marketaux API
func (m *Monitor) fetchMarketauxNews() ([]Article, error) {
	apiKey, err := m.authManager.GetAPIKey("marketaux")
	if err != nil {
		return nil, fmt.Errorf("failed to get Marketaux API key: %w", err)
	}

	baseURL := "https://api.marketaux.com/v1/news/all"
	params := url.Values{}
	params.Add("api_token", apiKey)
	params.Add("language", "en")
	params.Add("limit", "10")

	// Add keywords if configured
	if len(m.config.Keywords) > 0 {
		params.Add("keywords", strings.Join(m.config.Keywords, ","))
	}

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get data, status: %d, body: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response struct {
		Data []struct {
			Title       string    `json:"title"`
			Description string    `json:"description"`
			URL         string    `json:"url"`
			Source      string    `json:"source"`
			PublishedAt string    `json:"published_at"`
			Sentiment   float64   `json:"sentiment"`
			Entities    []struct {
				Symbol string `json:"symbol"`
			} `json:"entities"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	articles := make([]Article, 0, len(response.Data))
	for _, item := range response.Data {
		publishedAt, _ := time.Parse(time.RFC3339, item.PublishedAt)

		symbols := make([]string, 0, len(item.Entities))
		for _, entity := range item.Entities {
			if entity.Symbol != "" {
				symbols = append(symbols, entity.Symbol)
			}
		}

		article := Article{
			Title:       item.Title,
			Description: item.Description,
			URL:         item.URL,
			Source:      item.Source,
			PublishedAt: publishedAt,
			Sentiment:   item.Sentiment,
			Symbols:     symbols,
			Keywords:    extractKeywords(item.Title + " " + item.Description),
		}

		articles = append(articles, article)
	}

	return articles, nil
}

// fetchTwitterNews fetches financial news from Twitter
func (m *Monitor) fetchTwitterNews() ([]Article, error) {
	// Use the Twitter API from the datasource module
	// This is a simplified implementation that would need to be expanded
	// based on the actual Twitter API implementation

	// Create a list of search queries based on keywords and stock symbols
	queries := m.config.Keywords

	articles := make([]Article, 0)

	for _, query := range queries {
		// Mock implementation - in a real scenario, this would call the Twitter API
		// through the datasource module
		mockArticles := createMockTwitterArticles(query, 5)
		articles = append(articles, mockArticles...)
	}

	return articles, nil
}

// Helper function to extract keywords from text
func extractKeywords(text string) []string {
	// This is a simplified implementation
	// In a real scenario, this would use NLP techniques to extract relevant keywords
	
	// Convert to lowercase
	text = strings.ToLower(text)
	
	// Remove punctuation
	text = strings.Map(func(r rune) rune {
		if r == '.' || r == ',' || r == '!' || r == '?' || r == ';' || r == ':' || r == '"' || r == '\'' {
			return ' '
		}
		return r
	}, text)
	
	// Split into words
	words := strings.Fields(text)
	
	// Filter out common stop words
	stopWords := map[string]bool{
		"a": true, "an": true, "the": true, "and": true, "or": true, "but": true,
		"is": true, "are": true, "was": true, "were": true, "be": true, "been": true,
		"to": true, "of": true, "in": true, "for": true, "with": true, "on": true,
		"at": true, "by": true, "from": true, "up": true, "down": true, "this": true,
		"that": true, "these": true, "those": true, "it": true, "its": true,
	}
	
	filtered := make([]string, 0)
	seen := make(map[string]bool)
	
	for _, word := range words {
		if !stopWords[word] && len(word) > 3 && !seen[word] {
			filtered = append(filtered, word)
			seen[word] = true
			
			// Limit to a reasonable number of keywords
			if len(filtered) >= 10 {
				break
			}
		}
	}
	
	return filtered
}

// Helper function to create mock Twitter articles for testing
func createMockTwitterArticles(query string, count int) []Article {
	articles := make([]Article, count)
	
	for i := 0; i < count; i++ {
		// Extract potential stock symbols from the query
		symbols := make([]string, 0)
		words := strings.Fields(query)
		for _, word := range words {
			if strings.HasPrefix(word, "$") {
				symbol := strings.TrimPrefix(word, "$")
				symbols = append(symbols, symbol)
			}
		}
		
		// If no symbols were found, add some common ones
		if len(symbols) == 0 {
			symbols = []string{"AAPL", "MSFT", "GOOGL", "AMZN"}
		}
		
		// Generate a random sentiment between -1 and 1
		sentiment := float64(i%3-1) * 0.5 // -0.5, 0, or 0.5
		
		articles[i] = Article{
			Title:       fmt.Sprintf("Latest update on %s - Tweet %d", strings.Join(symbols, ", "), i+1),
			Description: fmt.Sprintf("This is a mock Twitter post about %s with some financial insights. #stocks #investing", query),
			URL:         fmt.Sprintf("https://twitter.com/user/status/%d", 1000000000+i),
			Source:      "Twitter",
			PublishedAt: time.Now().Add(-time.Duration(i) * time.Hour),
			Sentiment:   sentiment,
			Symbols:     symbols,
			Keywords:    []string{"stocks", "investing", "finance", query},
		}
	}
	
	return articles
}
