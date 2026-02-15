package metering

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// PostgresMeter implements Meter with PostgreSQL storage.
type PostgresMeter struct {
	db *sql.DB
}

// NewPostgresMeter creates a new PostgreSQL-backed meter.
func NewPostgresMeter(db *sql.DB) *PostgresMeter {
	return &PostgresMeter{db: db}
}

const schema = `
CREATE TABLE IF NOT EXISTS usage_events (
	id BIGSERIAL PRIMARY KEY,
	tenant_id TEXT NOT NULL,
	event_type TEXT NOT NULL,
	quantity BIGINT NOT NULL,
	timestamp TIMESTAMP NOT NULL,
	metadata JSONB
);
CREATE INDEX IF NOT EXISTS idx_usage_events_tenant_time ON usage_events(tenant_id, timestamp);
`

// Init creates the necessary database tables.
func (m *PostgresMeter) Init(ctx context.Context) error {
	_, err := m.db.ExecContext(ctx, schema)
	return err
}

// Record stores a single usage event.
func (m *PostgresMeter) Record(ctx context.Context, event Event) error {
	if err := event.Validate(); err != nil {
		return err
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}

	var metadataJSON []byte
	var err error
	if event.Metadata != nil {
		metadataJSON, err = json.Marshal(event.Metadata)
		if err != nil {
			return fmt.Errorf("metering: failed to marshal metadata: %w", err)
		}
	}

	_, err = m.db.ExecContext(ctx, `
		INSERT INTO usage_events (tenant_id, event_type, quantity, timestamp, metadata)
		VALUES ($1, $2, $3, $4, $5)
	`, event.TenantID, event.EventType, event.Quantity, event.Timestamp, metadataJSON)

	if err != nil {
		return fmt.Errorf("metering: failed to record event: %w", err)
	}
	return nil
}

// RecordBatch stores multiple events in a single transaction.
func (m *PostgresMeter) RecordBatch(ctx context.Context, events []Event) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("metering: failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO usage_events (tenant_id, event_type, quantity, timestamp, metadata)
		VALUES ($1, $2, $3, $4, $5)
	`)
	if err != nil {
		return fmt.Errorf("metering: failed to prepare statement: %w", err)
	}
	defer func() { _ = stmt.Close() }()

	now := time.Now().UTC()
	for _, event := range events {
		if err := event.Validate(); err != nil {
			return err
		}
		if event.Timestamp.IsZero() {
			event.Timestamp = now
		}

		var metadataJSON []byte
		if event.Metadata != nil {
			var err error
			metadataJSON, err = json.Marshal(event.Metadata)
			if err != nil {
				return fmt.Errorf("metering: failed to marshal metadata: %w", err)
			}
		}

		_, err := stmt.ExecContext(ctx, event.TenantID, event.EventType, event.Quantity, event.Timestamp, metadataJSON)
		if err != nil {
			return fmt.Errorf("metering: failed to insert event: %w", err)
		}
	}

	return tx.Commit()
}

// GetUsage retrieves aggregated usage for all event types.
func (m *PostgresMeter) GetUsage(ctx context.Context, tenantID string, period Period) (*Usage, error) {
	rows, err := m.db.QueryContext(ctx, `
		SELECT event_type, SUM(quantity) as total
		FROM usage_events
		WHERE tenant_id = $1 AND timestamp >= $2 AND timestamp < $3
		GROUP BY event_type
	`, tenantID, period.Start, period.End)
	if err != nil {
		return nil, fmt.Errorf("metering: failed to query usage: %w", err)
	}
	defer func() { _ = rows.Close() }()

	usage := &Usage{
		TenantID:   tenantID,
		Period:     period,
		Totals:     make(map[EventType]int64),
		LastUpdate: time.Now().UTC(),
	}

	for rows.Next() {
		var eventType EventType
		var total int64
		if err := rows.Scan(&eventType, &total); err != nil {
			return nil, fmt.Errorf("metering: failed to scan row: %w", err)
		}
		usage.Totals[eventType] = total
	}

	return usage, rows.Err()
}

// GetUsageByType retrieves usage for a specific event type.
func (m *PostgresMeter) GetUsageByType(ctx context.Context, tenantID string, eventType EventType, period Period) (int64, error) {
	var total sql.NullInt64
	err := m.db.QueryRowContext(ctx, `
		SELECT SUM(quantity)
		FROM usage_events
		WHERE tenant_id = $1 AND event_type = $2 AND timestamp >= $3 AND timestamp < $4
	`, tenantID, eventType, period.Start, period.End).Scan(&total)

	if err != nil {
		return 0, fmt.Errorf("metering: failed to query usage by type: %w", err)
	}

	return total.Int64, nil
}
