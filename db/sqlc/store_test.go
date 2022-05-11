package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	//run n cuncurrent transfer transactions
	n := 5
	amount := int64(10)
	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}
	//check the results
	for i := 0; i < n; i++ {
		err := <-errs // pass the errors recieved from channel  and  check the result
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		//check the transfer
		transfer := result.Transfer

		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		// get the  account entry from the database and  make sure it was created
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//check  entries
		fromEntry := result.FromEntry

		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		//require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		// get the  account entry from the database and  make sure it was created

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		//check  toentries
		toEntry := result.ToEntry

		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		//require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		// get the  account entry from the database and  make sure it was created

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//TODO account balance
	}

}
