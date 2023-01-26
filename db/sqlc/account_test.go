package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/HeavenAQ/simple-bank/utils"
	"github.com/stretchr/testify/require"
)

func CreateRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	//ensure account could be added successfully
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	// ensure ID and Created are auto generated
	require.NotZero(t, account.ID)
	require.NotZero(t, account.Created)
	return account
}

func TestCreateAccount(t *testing.T) {
	CreateRandomAccount(t)
}

func TestUpdateAccount(t *testing.T) {
	account := CreateRandomAccount(t)
	arg := UpdateAccountParams{
		ID:      account.ID,
		Balance: utils.RandomMoney(),
	}

	account, err := testQueries.UpdateAccount(context.Background(), arg)

	// Ensure the returned value is not empty
	require.NoError(t, err)
	require.NotEmpty(t, account)

	// Ensure the value is updated
	require.Equal(t, arg.ID, account.ID)
	require.Equal(t, arg.Balance, account.Balance)
}

// ensure account data can be fetched
func TestGetAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.Created, account2.Created)
}

// ensure all accounts can be fetched
func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

// Ensure account is deleted
func TestDeleteAccount(t *testing.T) {
	account := CreateRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}
