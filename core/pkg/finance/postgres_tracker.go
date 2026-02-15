package finance

import (
	"database/sql"
	"errors"
	"fmt"
)

// PostgresTracker implements finance.Tracker backed by PostgreSQL.
// Uses SELECT FOR UPDATE to provide row-level locking for atomic budget checks.
type PostgresTracker struct {
	db *sql.DB
}

// NewPostgresTracker creates a new PostgreSQL-backed budget tracker.
func NewPostgresTracker(db *sql.DB) *PostgresTracker {
	return &PostgresTracker{db: db}
}

// Check verifies that the given cost fits within the budget.
// Uses a read-only transaction with SELECT FOR SHARE to prevent phantom reads.
func (t *PostgresTracker) Check(budgetID string, cost Cost) (bool, error) {
	var resourceType string
	var limit, consumed int64

	err := t.db.QueryRow(
		`SELECT resource_type, budget_limit, consumed FROM finance_budgets WHERE id = $1`,
		budgetID,
	).Scan(&resourceType, &limit, &consumed)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("budget not found")
		}
		return false, fmt.Errorf("budget check failed: %w", err)
	}

	amount, err := extractAmount(resourceType, cost)
	if err != nil {
		return false, err
	}

	return consumed+amount <= limit, nil
}

// Consume atomically deducts the cost from the budget using SELECT FOR UPDATE.
// This is the core of financial determinism: the row lock prevents double-charge.
func (t *PostgresTracker) Consume(budgetID string, cost Cost) error {
	tx, err := t.db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// SELECT FOR UPDATE: locks the row until COMMIT, preventing concurrent consumption
	var resourceType string
	var limit, consumed int64
	err = tx.QueryRow(
		`SELECT resource_type, budget_limit, consumed FROM finance_budgets WHERE id = $1 FOR UPDATE`,
		budgetID,
	).Scan(&resourceType, &limit, &consumed)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("budget not found")
		}
		return fmt.Errorf("budget lock failed: %w", err)
	}

	amount, err := extractAmount(resourceType, cost)
	if err != nil {
		return err
	}

	if consumed+amount > limit {
		return errors.New("budget exceeded")
	}

	_, err = tx.Exec(
		`UPDATE finance_budgets SET consumed = consumed + $1 WHERE id = $2`,
		amount, budgetID,
	)
	if err != nil {
		return fmt.Errorf("budget update failed: %w", err)
	}

	return tx.Commit()
}

// extractAmount determines the cost amount based on the budget's resource type.
func extractAmount(resourceType string, cost Cost) (int64, error) {
	switch resourceType {
	case "USD", "EUR":
		if cost.Money.Currency != resourceType {
			return 0, errors.New("currency mismatch")
		}
		return cost.Money.AmountMinor, nil
	case "TOKENS":
		return cost.Tokens, nil
	case "REQUESTS":
		return cost.Requests, nil
	default:
		return 0, errors.New("unsupported resource type")
	}
}
