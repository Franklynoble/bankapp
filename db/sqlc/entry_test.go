package db

import (
	"context"
	"testing"
	"time"

	"github.com/Franklynoble/bankapp/db/util"
	"github.com/stretchr/testify/require"
)

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	CreateRandomEntry(t, account)

	//return account2

}

func CreateRandomEntry(t *testing.T, account Account) Entry {
	arg := createEntriesParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}
	entry, err := testQueries.createEntries(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry

}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry1 := CreateRandomEntry(t, account)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	account := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		CreateRandomEntry(t, account)
	}
	args := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}
	entries, err := testQueries.ListEntries(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, args.AccountID, entry.AccountID)

	}
}
