package db

import (
	"context"
	"testing"

	"github.com/Franklynoble/bankapp/db/util"
	"github.com/stretchr/testify/require"
	//_ "github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	//this would automatically check if the error is nil and will automattically check the error if it is not nil
	require.NoError(t, err)

	require.NotEmpty(t, account) //check that the account is not an empty object

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

}
