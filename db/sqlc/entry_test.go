package db

import (
	"context"
	"testing"
	"time"

	"github.com/Agentic_Bank_Server/utils"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) (Entry, Account) {
	account := createRandomAccount(t)
	args := CreateEntryParams{
		AccountID: account.ID,
		Amount:    utils.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, account.ID, entry.AccountID)
	require.Equal(t, args.Amount, entry.Amount)
	return entry, account
}
func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	entry1, _ := createRandomEntry(t)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, 2*time.Second)
}
