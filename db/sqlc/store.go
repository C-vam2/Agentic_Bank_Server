package db

import (
	"context"
	"database/sql"
	"fmt"
)

var txKey = struct{}{}

// Store provides all functions to execure db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	queries := New(tx)
	err = fn(queries)

	if err != nil {
		rbErr := tx.Rollback()

		if rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err:%v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to the other.
// It creates a transfer record, add account entries, and update account's balance within a single database transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		// txName := ctx.Value(txKey)

		// fmt.Println(txName, "create transfer")
		if result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		}); err != nil {
			return err
		}

		// fmt.Println(txName, "Get account for update 1")
		account1, err := q.GetAccountForUpdate(context.Background(), arg.FromAccountID)
		if err != nil {
			return err
		}

		// fmt.Println(txName, "update account 1")
		result.FromAccount, err = q.UpdateAccount(context.Background(), UpdateAccountParams{
			ID:      arg.FromAccountID,
			Balance: account1.Balance - arg.Amount,
		})
		if err != nil {
			return err
		}

		// fmt.Println(txName, "Get account for update 2")
		account2, err := q.GetAccountForUpdate(context.Background(), arg.ToAccountID)
		if err != nil {
			return err
		}

		// fmt.Println(txName, "update account 2")
		result.ToAccount, err = q.UpdateAccount(context.Background(), UpdateAccountParams{ID: arg.ToAccountID, Balance: account2.Balance + arg.Amount})
		if err != nil {
			return err
		}

		// fmt.Println(txName, "create entry 1")
		if result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{AccountID: arg.FromAccountID, Amount: -arg.Amount}); err != nil {
			return err
		}

		// fmt.Println(txName, "create entry 2")
		if result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{AccountID: arg.ToAccountID, Amount: arg.Amount}); err != nil {
			return err
		}

		return nil
	})
	return result, err
}
