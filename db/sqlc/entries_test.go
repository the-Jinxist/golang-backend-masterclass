package backend_masterclass

import (
	"backend_masterclass/util"
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateEntry(t *testing.T) {

	arg := CreateEntryParams{
		AccountID: sql.NullInt64{
			Int64: util.RandomInt(1, 10),
			Valid: true,
		},
		Amount: util.RandomMoney(),
	}

	entries, err := testQueries.CreateEntry(context.Background(), arg)
	require.Empty(t, err)
	require.NotEmpty(t, entries)
	require.Equal(t, entries.AccountID, arg.AccountID)
}

func TestGetEntry(t *testing.T) {
	createEntryParams := CreateEntryParams{
		AccountID: sql.NullInt64{
			Int64: 20,
			Valid: true,
		},
		Amount: util.RandomMoney(),
	}

	entries, err := testQueries.CreateEntry(context.Background(), createEntryParams)
	require.Empty(t, err)
	require.NotEmpty(t, entries)
	require.Equal(t, entries.AccountID, createEntryParams.AccountID)

	getEntry, err := testQueries.GetEntry(context.Background(), entries.AccountID.Int64)
	require.Empty(t, err)
	require.NotEmpty(t, getEntry)
	require.Equal(t, entries.AccountID, getEntry.AccountID)
	require.Equal(t, entries.Amount, getEntry.Amount)
}

func TestListEntries(t *testing.T) {
	for i := 0; i < 10; i++ {
		createEntryParams := CreateEntryParams{
			AccountID: sql.NullInt64{
				Int64: util.RandomInt(1, 100),
				Valid: true,
			},
			Amount: util.RandomMoney(),
		}

		entries, err := testQueries.CreateEntry(context.Background(), createEntryParams)
		require.Empty(t, err)
		require.NotEmpty(t, entries)
		require.Equal(t, entries.AccountID, createEntryParams.AccountID)
	}

	listEntryParams := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	entries, err := testQueries.ListEntries(context.Background(), listEntryParams)
	require.Empty(t, err)
	require.NotEmpty(t, entries)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.NotEmpty(t, entry.Amount)
		require.NotEmpty(t, entry.AccountID)
	}
}
