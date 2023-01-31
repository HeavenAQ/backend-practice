package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type EntryType int64

const (
	from EntryType = 0
	to
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	// run a concurrent database transaction
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

	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccountID, account1.ID)
		require.Equal(t, transfer.ToAccountID, account2.ID)
		require.Equal(t, transfer.Amount, amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.Created)

		// check entry
		chkEntry := func(ty EntryType) {
			fromEntry := result.FromEntry
			require.NotEmpty(t, fromEntry)
			require.Equal(t, fromEntry.AccountID, account1.ID)
			require.NotZero(t, fromEntry.ID)
			require.NotZero(t, fromEntry.Created)
			if ty == from {
				require.Equal(t, fromEntry.Amount, -amount)
			} else {
				require.Equal(t, fromEntry.Amount, amount)
			}
		}
		chkEntry(from)
		chkEntry(to)

		// check account
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, account1.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, account2.ID)

		fmt.Printf("Tx ==> Account1: %v Account2: %v\n", fromAccount.Balance, toAccount.Balance)

		// check balance
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)

		// check difference is the multiples of the given amount
		require.True(t, diff1%amount == 0)
		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount1)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount2)

	fmt.Printf("After ==> Account1: %v Account2: %v\n", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance-int64(amount)*int64(n), updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(amount)*int64(n), updatedAccount2.Balance)

}

func TestTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	// run a concurrent database transaction
	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		transfer := func(reversed bool) {
			fromAccount := account1
			toAccount := account2

			if reversed {
				fromAccount = account2
				toAccount = account1
			}
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        amount,
			})
			errs <- err
		}
		go transfer(i%2 == 0)
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount1)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount2)

	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
