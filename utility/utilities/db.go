package utilities

import (
	"context"
	"database/sql"
	"fmt"

	"gitlab.com/tuneverse/toolkit/core/logger"
)

// ExecuteTx executes a transaction on the database.
// It takes a context and a function that takes a transaction and returns an error.
// It returns an error if the transaction fails to begin, rollback, or commit.
func ExecuteTx(ctx context.Context, db *sql.DB, fn func(tx *sql.Tx) error) error {

	log := logger.Log().WithContext(ctx)

	// Begin a db transaction with the provided context.
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Errorf("ExecuteTx: Failed to begin transaction: %s", err.Error())
		return err
	}

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			// Log the error if a rollback operation fails.
			log.Errorf("ExecuteTx: Rollback failed: %s", rbErr.Error())
			return fmt.Errorf("tx error %v, rbErr %v", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		// Log the error if the commit operation fails.
		log.Errorf("ExecuteTx: Commit failed: %s", err.Error())
		return err
	}

	return nil
}
