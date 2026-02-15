package finance

import (
	"errors"
	"sync"
	"time"
)

type WindowType string

const (
	WindowDaily   WindowType = "DAILY"
	WindowWeekly  WindowType = "WEEKLY"
	WindowMonthly WindowType = "MONTHLY"
	WindowTotal   WindowType = "TOTAL"
)

// Budget represents a limit on resource consumption.
type Budget struct {
	ID           string     `json:"id"`
	ResourceType string     `json:"resource_type"` // "USD", "TOKENS"
	Limit        int64      `json:"limit"`         // Minor units for money
	Window       WindowType `json:"window"`
	Consumed     int64      `json:"consumed"`
	ResetAt      time.Time  `json:"reset_at"`
}

// Tracker enforces budget changes.
type Tracker interface {
	Check(budgetID string, cost Cost) (bool, error)
	Consume(budgetID string, cost Cost) error
}

// InMemoryTracker is a simple thread-safe budget tracker.
type InMemoryTracker struct {
	mu      sync.RWMutex
	budgets map[string]*Budget
}

func NewInMemoryTracker() *InMemoryTracker {
	return &InMemoryTracker{
		budgets: make(map[string]*Budget),
	}
}

func (t *InMemoryTracker) SetBudget(b Budget) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.budgets[b.ID] = &b
}

func (t *InMemoryTracker) Check(budgetID string, cost Cost) (bool, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	b, ok := t.budgets[budgetID]
	if !ok {
		return false, errors.New("budget not found")
	}

	// Determine cost amount based on resource type
	var amount int64
	switch b.ResourceType {
	case "USD", "EUR":
		if cost.Money.Currency != b.ResourceType {
			return false, errors.New("currency mismatch") // Simplification: strict match
		}
		amount = cost.Money.AmountMinor
	case "TOKENS":
		amount = cost.Tokens
	case "REQUESTS":
		amount = cost.Requests
	default:
		return false, errors.New("unsupported resource type")
	}

	return b.Consumed+amount <= b.Limit, nil
}

func (t *InMemoryTracker) Consume(budgetID string, cost Cost) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Re-check inside lock for atomicity
	b, ok := t.budgets[budgetID]
	if !ok {
		return errors.New("budget not found")
	}

	// Determine cost amount based on resource type
	var amount int64
	switch b.ResourceType {
	case "USD", "EUR":
		if cost.Money.Currency != b.ResourceType {
			return errors.New("currency mismatch")
		}
		amount = cost.Money.AmountMinor
	case "TOKENS":
		amount = cost.Tokens
	case "REQUESTS":
		amount = cost.Requests
	default:
		return errors.New("unsupported resource type")
	}

	if b.Consumed+amount > b.Limit {
		return errors.New("budget exceeded")
	}

	// Update consumption
	switch b.ResourceType {
	case "USD", "EUR":
		b.Consumed += cost.Money.AmountMinor
	case "TOKENS":
		b.Consumed += cost.Tokens
	case "REQUESTS":
		b.Consumed += cost.Requests
	}

	return nil
}
