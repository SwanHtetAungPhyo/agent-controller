package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Store provides all functions to execute database queries and transactions
type Store interface {
	Querier
	// Add transaction methods here
}

// SQLStore implements Store interface
type SQLStore struct {
	pool *pgxpool.Pool
	*Queries
}

// NewStore creates a new Store instance
func NewStore(pool *pgxpool.Pool) Store {
	return &SQLStore{
		pool:    pool,
		Queries: New(pool),
	}
}

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	// Begin transaction
	tx, err := store.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Create queries with transaction
	q := New(tx)

	// Execute function
	err = fn(q)
	if err != nil {
		// Rollback on error
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	// Commit transaction
	return tx.Commit(ctx)
}
