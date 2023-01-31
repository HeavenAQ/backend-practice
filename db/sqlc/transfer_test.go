package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/HeavenAQ/simple-bank/utils"
	"github.com/stretchr/testify/require"
)

func CreateRandomTransfer(t *testing.T) Transfer {
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)
	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        utils.RandomMoney(),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, transfer.Amount, arg.Amount)
	require.Equal(t, transfer.FromAccountID, arg.FromAccountID)
	require.Equal(t, transfer.ToAccountID, arg.ToAccountID)

	require.NotEmpty(t, transfer.ID)
	require.NotEmpty(t, transfer.Created)
	return transfer
}

func TestCreateTransfer(t *testing.T) {
	CreateRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	transfer1 := CreateRandomTransfer(t)
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)
	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
}

func TestUpdateTransfer(t *testing.T) {
	transfer1 := CreateRandomTransfer(t)
	args := UpdateTransferParams{
		ID:            transfer1.ID,
		FromAccountID: transfer1.FromAccountID,
		ToAccountID:   transfer1.ToAccountID,
		Amount:        utils.RandomMoney(),
	}

	transfer2, err := testQueries.UpdateTransfer(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)
	require.Equal(t, args.ID, transfer2.ID)
	require.Equal(t, args.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, args.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, args.Amount, transfer2.Amount)
}

func TestListTransfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomTransfer(t)
	}

	args := ListTransfersParams{
		Limit:  5,
		Offset: 5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), args)
	require.NoError(t, err)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}

func TestDeleteTransfers(t *testing.T) {
	transfer1 := CreateRandomTransfer(t)
	err := testQueries.DeleteTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)

	transfer2, err2 := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.Empty(t, transfer2)
	require.Error(t, err2)
	require.EqualError(t, err2, sql.ErrNoRows.Error())
}
