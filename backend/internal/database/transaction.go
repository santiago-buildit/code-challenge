package database

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
)

// WithTransaction wraps a function with a DB transaction (for functions with single error return)
func WithTransaction(ctx context.Context, db *sqlx.DB, fn func(tx *sqlx.Tx) error) error {
	_, err := withTransactionInternal(ctx, db, func(tx *sqlx.Tx) (struct{}, error) {
		return struct{}{}, fn(tx)
	})
	return err
}

// WithTransactionResult wraps a function with a DB transaction (for functions with result and error return)
func WithTransactionResult[T any](ctx context.Context, db *sqlx.DB, fn func(tx *sqlx.Tx) (T, error)) (T, error) {
	return withTransactionInternal(ctx, db, fn)
}

// withTransactionInternal implements common wrapping logic
func withTransactionInternal[T any](ctx context.Context, db *sqlx.DB, fn func(tx *sqlx.Tx) (T, error)) (T, error) {

	var zero T // Zero value of T

	// Begin transaction
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return zero, err
	}

	// Defer rollback (no effect if tx.Commit() is called)
	defer func() {
		if rErr := tx.Rollback(); rErr != nil && !errors.Is(rErr, sql.ErrTxDone) {
			log.Printf("rollback failed: %v", rErr)
		}
	}()

	// Run the wrapped function within the transaction
	result, err := fn(tx)
	if err != nil {
		return zero, err // Rollback executed on defer
	}

	// No errors, commit the transaction
	if err := tx.Commit(); err != nil {
		return zero, err
	}

	return result, nil
}
