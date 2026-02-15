package memory

import (
	"database/sql"
)

// PostgresMemoryStore provides a stub for memory storage.
// Full ingestion pipeline is not part of the kernel TCB.
type PostgresMemoryStore struct {
	db *sql.DB
}

// NewPostgresMemoryStore creates a new store instance.
func NewPostgresMemoryStore(db *sql.DB) *PostgresMemoryStore {
	return &PostgresMemoryStore{db: db}
}
