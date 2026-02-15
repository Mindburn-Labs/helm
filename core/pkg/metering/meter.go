// Package metering provides per-tenant usage tracking for HELM.
// It tracks requests, tool invocations, LLM tokens, and storage usage.
package metering

import (
	"context"
	"errors"
	"time"
)

var (
	// ErrEmptyTenantID is returned when a metering event has no tenant ID.
	ErrEmptyTenantID = errors.New("metering: tenant_id must not be empty")
	// ErrNegativeQuantity is returned when a metering event has a negative quantity.
	ErrNegativeQuantity = errors.New("metering: quantity must not be negative")
	// ErrInvalidEventType is returned when the event type is empty.
	ErrInvalidEventType = errors.New("metering: event_type must not be empty")
)

// EventType defines the type of metered event.
type EventType string

const (
	EventRequest      EventType = "request"
	EventToolCall     EventType = "tool_call"
	EventLLMToken     EventType = "llm_token"
	EventStorageByte  EventType = "storage_byte"
	EventExecution    EventType = "execution"
	EventReceiptStore EventType = "receipt_store"
	EventIngestion    EventType = "ingestion"
)

// Event represents a single metered usage event.
type Event struct {
	TenantID  string         `json:"tenant_id"`
	EventType EventType      `json:"event_type"`
	Quantity  int64          `json:"quantity"`
	Timestamp time.Time      `json:"timestamp"`
	Metadata  map[string]any `json:"metadata,omitempty"`
}

// Validate checks that the event has valid fields.
func (e Event) Validate() error {
	if e.TenantID == "" {
		return ErrEmptyTenantID
	}
	if e.Quantity < 0 {
		return ErrNegativeQuantity
	}
	if e.EventType == "" {
		return ErrInvalidEventType
	}
	return nil
}

// Period defines a time range for usage aggregation.
type Period struct {
	Start time.Time
	End   time.Time
}

// DailyPeriod returns a Period for the current day.
func DailyPeriod() Period {
	now := time.Now().UTC()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return Period{Start: start, End: start.Add(24 * time.Hour)}
}

// MonthlyPeriod returns a Period for the current month.
func MonthlyPeriod() Period {
	now := time.Now().UTC()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)
	return Period{Start: start, End: end}
}

// Usage contains aggregated usage for a tenant.
type Usage struct {
	TenantID   string
	Period     Period
	Totals     map[EventType]int64
	LastUpdate time.Time
}

// Meter is the interface for recording and querying usage.
type Meter interface {
	// Record stores a usage event.
	Record(ctx context.Context, event Event) error

	// RecordBatch stores multiple events atomically.
	RecordBatch(ctx context.Context, events []Event) error

	// GetUsage retrieves aggregated usage for a tenant in a period.
	GetUsage(ctx context.Context, tenantID string, period Period) (*Usage, error)

	// GetUsageByType retrieves usage for a specific event type.
	GetUsageByType(ctx context.Context, tenantID string, eventType EventType, period Period) (int64, error)
}
