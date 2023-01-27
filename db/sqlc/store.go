package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTrans(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal("Failed to initialize store structure", err)
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		log.Fatal("Failed to execute target query. Starting to roll back...", err)
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("=== Failed to execute database transaction ===\nError: %v\nRollback Error: %v", err, rbErr)
		}
	}
	// commit tx and return its error
	return tx.Commit()
}

// structure for transfer execution and transfer result
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

func (store *Store) TransferTx(ctx context.Context, data TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	// txName := ctx.Value(txKey) -- for debugging

	err := store.execTrans(ctx, func(q *Queries) error {
		// declare error var
		var err error

		// create transfer
		if result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: data.FromAccountID,
			ToAccountID:   data.ToAccountID,
			Amount:        data.Amount,
		}); err != nil {
			return err
		}

		// create entry record for the account transferring money
		if result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: data.FromAccountID,
			Amount:    -data.Amount,
		}); err != nil {
			return err
		}

		// create entry record for the account receiving money
		if result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: data.ToAccountID,
			Amount:    data.Amount,
		}); err != nil {
			return err
		}

		// INFO: Deal with the account balance
		if result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     data.FromAccountID,
			Amount: -data.Amount,
		}); err != nil {
			return err
		}

		if result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     data.ToAccountID,
			Amount: data.Amount,
		}); err != nil {
			return err
		}

		return nil
	})

	return result, err
}
