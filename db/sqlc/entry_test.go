package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/HeavenAQ/simple-bank/utils"
	"github.com/stretchr/testify/require"
)

func CreateRandomEntry(t *testing.T) Entry {
	account := CreateRandomAccount(t)
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    utils.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, entry.AccountID, arg.AccountID)
	require.Equal(t, entry.Amount, arg.Amount)
	require.NotEmpty(t, entry.Created)
	require.NotEmpty(t, entry.ID)
	return entry
}

func TestCreateEntry(t *testing.T) {
	CreateRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	entry1 := CreateRandomEntry(t)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
}

func TestListEntries(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomEntry(t)
	}

	arg := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entries)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}

func TestUpdateEntry(t *testing.T) {
	entry1 := CreateRandomEntry(t)
	arg := UpdateEntryParams{
		ID:        entry1.ID,
		AccountID: entry1.AccountID,
		Amount:    utils.RandomMoney(),
	}

	entry, err := testQueries.UpdateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, entry.ID, arg.ID)
	require.Equal(t, entry.AccountID, arg.AccountID)
	require.Equal(t, entry.Amount, arg.Amount)
}

func TestDeleteEntry(t *testing.T) {
	entry := CreateRandomEntry(t)
	err := testQueries.DeleteEntry(context.Background(), entry.ID)
	require.NoError(t, err)

	entry2, err2 := testQueries.GetEntry(context.Background(), entry.ID)
	require.Empty(t, entry2)
	require.Error(t, err2)
	require.EqualError(t, err2, sql.ErrNoRows.Error())
}
