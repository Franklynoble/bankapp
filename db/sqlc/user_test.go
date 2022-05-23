package db

import (
	"context"
	"testing"
	"time"

	"github.com/Franklynoble/bankapp/db/util"
	"github.com/stretchr/testify/require"
	//_ "github.com/stretchr/testify/require"
)

// this func does not have a Test prefix, therefore it would not be run as a unit test
func createRandomUser(t *testing.T) User {

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: "secret",
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	//this would automatically check if the error is nil and will automattically check the error if it is not nil
	require.NoError(t, err)
	require.NotEmpty(t, user) //check that the account is not an empty object
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user

}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)

}

func TestGetUser(t *testing.T) {

	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.Equal(t, user1.CreatedAt, user2.CreatedAt, time.Second)

}
