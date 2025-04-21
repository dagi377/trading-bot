package execution

import (
	"fmt"
	"sync"
	"time"

	"github.com/hustler/trading-bot/pkg/data"
	"github.com/hustler/trading-bot/pkg/strategy"
)

// TradeStatus represents the status of a trade
type TradeStatus string

const (
	Pending   TradeStatus = "PENDING"
	Executed  TradeStatus = "EXECUTED"
	Cancelled TradeStatus = "CANCELLED"
	Completed TradeStatus = "COMPLETED"
)

// Trade represents a trade
type Trade struct {
	ID        string
	Symbol    string
	Quantity  int
	Price     float64
	Type      strategy.TradeSignal
	Status    TradeStatus
	CreatedAt time.Time
	UpdatedAt time.Time
	Reason    string
}

// TradeManager manages trade execution
type TradeManager struct {
	trades         map[string]*Trade
	activeTrades   map[string]*Trade
	capitalPerStock float64
	maxLossPerTrade float64
	mu             sync.RWMutex
}

// NewTradeManager creates a new TradeManager
func NewTradeManager(capitalPerStock, maxLossPerTrade float64) *TradeManager {
	return &TradeManager{
		trades:         make(map[string]*Trade),
		activeTrades:   make(map[string]*Trade),
		capitalPerStock: capitalPerStock,
		maxLossPerTrade: maxLossPerTrade,
	}
}

// ExecuteTrade executes a trade based on a trade decision
func (t *TradeManager) ExecuteTrade(decision *strategy.TradeDecision, stock *data.Stock) (*Trade, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Check if we already have an active trade for this symbol
	if activeTrade, exists := t.getActiveTradeForSymbol(decision.Symbol); exists {
		// If we have an active trade and the decision is to sell, close the position
		if decision.Signal == strategy.Sell {
			return t.closePosition(activeTrade, decision, stock)
		}
		// If we have an active trade and the decision is not to sell, do nothing
		return nil, fmt.Errorf("already have an active trade for %s", decision.Symbol)
	}

	// If we don't have an active trade and the decision is to buy, open a position
	if decision.Signal == strategy.Buy {
		return t.openPosition(decision, stock)
	}

	// If we don't have an active trade and the decision is not to buy, do nothing
	return nil, fmt.Errorf("no action needed for %s", decision.Symbol)
}

// getActiveTradeForSymbol gets an active trade for a symbol
func (t *TradeManager) getActiveTradeForSymbol(symbol string) (*Trade, bool) {
	for _, trade := range t.activeTrades {
		if trade.Symbol == symbol {
			return trade, true
		}
	}
	return nil, false
}

// openPosition opens a new position
func (t *TradeManager) openPosition(decision *strategy.TradeDecision, stock *data.Stock) (*Trade, error) {
	// Calculate quantity based on capital per stock
	quantity := int(t.capitalPerStock / stock.CurrentPrice)
	if quantity <= 0 {
		return nil, fmt.Errorf("insufficient capital to buy %s at $%.2f", stock.Symbol, stock.CurrentPrice)
	}

	// Create a new trade
	trade := &Trade{
		ID:        fmt.Sprintf("%s-%d", stock.Symbol, time.Now().UnixNano()),
		Symbol:    stock.Symbol,
		Quantity:  quantity,
		Price:     stock.CurrentPrice,
		Type:      strategy.Buy,
		Status:    Executed,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Reason:    decision.Rationale,
	}

	// Add to trades and active trades
	t.trades[trade.ID] = trade
	t.activeTrades[trade.ID] = trade

	return trade, nil
}

// closePosition closes an existing position
func (t *TradeManager) closePosition(trade *Trade, decision *strategy.TradeDecision, stock *data.Stock) (*Trade, error) {
	// Create a new trade for the sell
	sellTrade := &Trade{
		ID:        fmt.Sprintf("%s-sell-%d", stock.Symbol, time.Now().UnixNano()),
		Symbol:    stock.Symbol,
		Quantity:  trade.Quantity,
		Price:     stock.CurrentPrice,
		Type:      strategy.Sell,
		Status:    Executed,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Reason:    decision.Rationale,
	}

	// Add to trades
	t.trades[sellTrade.ID] = sellTrade

	// Remove from active trades
	delete(t.activeTrades, trade.ID)

	// Update original trade
	trade.Status = Completed
	trade.UpdatedAt = time.Now()

	return sellTrade, nil
}

// CancelTrade cancels a trade
func (t *TradeManager) CancelTrade(tradeID string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	trade, exists := t.trades[tradeID]
	if !exists {
		return fmt.Errorf("trade not found: %s", tradeID)
	}

	if trade.Status == Completed {
		return fmt.Errorf("cannot cancel completed trade: %s", tradeID)
	}

	trade.Status = Cancelled
	trade.UpdatedAt = time.Now()

	// Remove from active trades if it's there
	delete(t.activeTrades, tradeID)

	return nil
}

// GetTrade gets a trade by ID
func (t *TradeManager) GetTrade(tradeID string) (*Trade, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	trade, exists := t.trades[tradeID]
	return trade, exists
}

// GetAllTrades gets all trades
func (t *TradeManager) GetAllTrades() []*Trade {
	t.mu.RLock()
	defer t.mu.RUnlock()

	trades := make([]*Trade, 0, len(t.trades))
	for _, trade := range t.trades {
		trades = append(trades, trade)
	}
	return trades
}

// GetActiveTrades gets all active trades
func (t *TradeManager) GetActiveTrades() []*Trade {
	t.mu.RLock()
	defer t.mu.RUnlock()

	trades := make([]*Trade, 0, len(t.activeTrades))
	for _, trade := range t.activeTrades {
		trades = append(trades, trade)
	}
	return trades
}

// CheckStopLoss checks if any active trades have hit their stop loss
func (t *TradeManager) CheckStopLoss(stocks map[string]*data.Stock) []*Trade {
	t.mu.Lock()
	defer t.mu.Unlock()

	closedTrades := make([]*Trade, 0)

	for id, trade := range t.activeTrades {
		stock, exists := stocks[trade.Symbol]
		if !exists {
			continue
		}

		// Calculate current value and loss
		currentValue := float64(trade.Quantity) * stock.CurrentPrice
		entryValue := float64(trade.Quantity) * trade.Price
		loss := entryValue - currentValue

		// If loss exceeds max loss per trade, close the position
		if loss > t.maxLossPerTrade {
			// Create a new trade for the sell
			sellTrade := &Trade{
				ID:        fmt.Sprintf("%s-stoploss-%d", trade.Symbol, time.Now().UnixNano()),
				Symbol:    trade.Symbol,
				Quantity:  trade.Quantity,
				Price:     stock.CurrentPrice,
				Type:      strategy.Sell,
				Status:    Executed,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Reason:    fmt.Sprintf("Stop loss triggered: Loss of $%.2f exceeds max loss of $%.2f", loss, t.maxLossPerTrade),
			}

			// Add to trades
			t.trades[sellTrade.ID] = sellTrade
			closedTrades = append(closedTrades, sellTrade)

			// Remove from active trades
			delete(t.activeTrades, id)

			// Update original trade
			trade.Status = Completed
			trade.UpdatedAt = time.Now()
		}
	}

	return closedTrades
}

// CloseAllPositions closes all active positions
func (t *TradeManager) CloseAllPositions(stocks map[string]*data.Stock) []*Trade {
	t.mu.Lock()
	defer t.mu.Unlock()

	closedTrades := make([]*Trade, 0)

	for id, trade := range t.activeTrades {
		stock, exists := stocks[trade.Symbol]
		if !exists {
			continue
		}

		// Create a new trade for the sell
		sellTrade := &Trade{
			ID:        fmt.Sprintf("%s-close-%d", trade.Symbol, time.Now().UnixNano()),
			Symbol:    trade.Symbol,
			Quantity:  trade.Quantity,
			Price:     stock.CurrentPrice,
			Type:      strategy.Sell,
			Status:    Executed,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Reason:    "End of trading day - closing all positions",
		}

		// Add to trades
		t.trades[sellTrade.ID] = sellTrade
		closedTrades = append(closedTrades, sellTrade)

		// Remove from active trades
		delete(t.activeTrades, id)

		// Update original trade
		trade.Status = Completed
		trade.UpdatedAt = time.Now()
	}

	return closedTrades
}
