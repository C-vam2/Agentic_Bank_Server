package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(dbConn)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// results := make(chan TransferTxResult)
	errs := make(chan error)

	n := 10
	var amount int64 = 10

	for i := 0; i < int(n); i++ {
		fromAccount := account1
		toAccount := account2
		if i%2 == 0 {
			fromAccount = account2
			toAccount = account1
		}

		go func() {

			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        amount,
			})
			errs <- err

		}()
	}

	//check results

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check the final updated balances
	var err error
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NotEmpty(t, updatedAccount1)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NotEmpty(t, updatedAccount2)
	require.NoError(t, err)

	require.Equal(t, updatedAccount1.Balance, account1.Balance)
	require.Equal(t, updatedAccount2.Balance, account2.Balance)

}
