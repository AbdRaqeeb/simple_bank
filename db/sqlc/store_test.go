package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDb)

	accountOne := createRandomAccount(t)
	accountTwo := createRandomAccount(t)
	fmt.Println(">> before:", accountOne.Balance, accountTwo.Balance)

	/**
	Run concurrent transfers to ensure the transactions works efficiently
	About five transfers will be run concurrently to test
	*/
	n := 5
	amount := int64(10)

	results := make(chan TransferTxResult)
	errs := make(chan error)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: accountOne.ID,
				ToAccountID:   accountTwo.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	existed := map[int]bool{}
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccountID, accountOne.ID)
		require.Equal(t, transfer.ToAccountID, accountTwo.ID)
		require.Equal(t, transfer.Amount, amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		// check transfer record in db
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check FromEntry
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromEntry.AccountID, accountOne.ID)
		require.Equal(t, fromEntry.Amount, -amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		// check FromEntry in db
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// check ToEntry
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toEntry.AccountID, accountTwo.ID)
		require.Equal(t, toEntry.Amount, amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		// check ToEntry in db
		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check FromAccount
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, accountOne.ID)

		// check ToAccount
		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, accountTwo.ID)

		// check Accounts Balances
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)

		diffOne := accountOne.Balance - fromAccount.Balance // this is the amount that was deducted from account one balance
		diffTwo := toAccount.Balance - accountTwo.Balance   // this is the amount that was added to account two balance
		require.Equal(t, diffOne, diffTwo)
		require.True(t, diffOne > 0)
		// 1 * amount, 2 * amount, .... n * amount.
		//Where n is the number of times the deduction was done because of the concurrent goroutines.
		require.True(t, diffOne%amount == 0)

		k := int(diffOne / amount)
		require.True(t, k >= 1 && k <= n)
		// k should not be in existed yet
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check updated account balances
	updatedAccountOne, err := testQueries.GetAccount(context.Background(), accountOne.ID)
	require.NoError(t, err)

	updatedAccountTwo, err := testQueries.GetAccount(context.Background(), accountTwo.ID)
	require.NoError(t, err)

	accountOneCurrentBalance := accountOne.Balance - int64(n)*amount // the amount times the number of times it was deducted
	accountTwoCurrentBalance := accountTwo.Balance + int64(n)*amount // the amount times the number of times it was added

	fmt.Println(">> after:", updatedAccountOne.Balance, updatedAccountTwo.Balance)
	require.Equal(t, accountOneCurrentBalance, updatedAccountOne.Balance)
	require.Equal(t, accountTwoCurrentBalance, updatedAccountTwo.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDb)

	accountOne := createRandomAccount(t)
	accountTwo := createRandomAccount(t)
	fmt.Println(">> before:", accountOne.Balance, accountTwo.Balance)

	/**
	Run concurrent transfers to ensure the transactions works efficiently
	About five transfers will be run concurrently to test
	*/
	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := accountOne.ID
		toAccountID := accountTwo.ID

		if i%2 == 1 {
			fromAccountID = accountTwo.ID
			toAccountID = accountOne.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check updated account balances
	updatedAccountOne, err := testQueries.GetAccount(context.Background(), accountOne.ID)
	require.NoError(t, err)

	updatedAccountTwo, err := testQueries.GetAccount(context.Background(), accountTwo.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccountOne.Balance, updatedAccountTwo.Balance)
	require.Equal(t, accountOne.Balance, updatedAccountOne.Balance)
	require.Equal(t, accountTwo.Balance, updatedAccountTwo.Balance)
}
