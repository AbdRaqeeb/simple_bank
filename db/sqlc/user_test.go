package db

import (
	"context"
	"github.com/AbdRaqeeb/simple_bank/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	_ = createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)
	foundUser, err := testQueries.GetUser(context.Background(), user.Username)

	require.NoError(t, err)
	require.NotEmpty(t, foundUser)
	require.Equal(t, user.Username, foundUser.Username)
	require.Equal(t, user.FullName, foundUser.FullName)
	require.Equal(t, user.Email, foundUser.Email)
	require.Equal(t, user.HashedPassword, foundUser.HashedPassword)
	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)
	require.WithinDuration(t, user.CreatedAt, foundUser.CreatedAt, time.Second)
}
