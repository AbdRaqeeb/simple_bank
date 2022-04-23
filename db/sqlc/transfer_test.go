package db

import (
	"context"
	"github.com/AbdRaqeeb/simple_bank/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomTransfer(t *testing.T, fromAccountID, toAccountID int64) Transfer {
	arg := CreateTransferParams{
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, transfer.FromAccountID, fromAccountID)
	require.Equal(t, transfer.ToAccountID, toAccountID)
	require.Equal(t, transfer.Amount, arg.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	fromAccount, toAccount := createRandomAccount(t), createRandomAccount(t)

	createRandomTransfer(t, fromAccount.ID, toAccount.ID)
}

func TestGetTransfer(t *testing.T) {
	fromAccount, toAccount := createRandomAccount(t), createRandomAccount(t)

	transfer := createRandomTransfer(t, fromAccount.ID, toAccount.ID)

	foundTransfer, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, foundTransfer)

	require.Equal(t, foundTransfer.ID, transfer.ID)
	require.Equal(t, foundTransfer.FromAccountID, transfer.FromAccountID)
	require.Equal(t, foundTransfer.ToAccountID, transfer.ToAccountID)
	require.Equal(t, foundTransfer.Amount, transfer.Amount)
	require.WithinDuration(t, transfer.CreatedAt, foundTransfer.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T) {
	fromAccount, toAccount := createRandomAccount(t), createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomTransfer(t, fromAccount.ID, toAccount.ID)
	}

	args := ListTransfersParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}
