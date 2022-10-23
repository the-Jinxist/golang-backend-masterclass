package backend_masterclass

import (
	"backend_masterclass/util"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateEntry(t *testing.T) {

	user := createRandomUser(t)

	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	//We used the testify library to check the test result
	accounts, err := testQueries.CreateAccount(context.Background(), arg)

	//require keyword will end the test if the assetion doesn't pass
	//assert keyword will log the failed assertion but will still continue the test
	assert.Nil(t, err)
	assert.NotNil(t, accounts)
	require.NotEmpty(t, accounts)

	createEntryArg := CreateEntryParams{
		AccountID: accounts.ID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), createEntryArg)

	require.Empty(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, entry.AccountID, accounts.ID)
}

func TestGetEntry(t *testing.T) {

	user := createRandomUser(t)

	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	//We used the testify library to check the test result
	accounts, err := testQueries.CreateAccount(context.Background(), arg)

	//require keyword will end the test if the assetion doesn't pass
	//assert keyword will log the failed assertion but will still continue the test
	assert.Nil(t, err)
	assert.NotNil(t, accounts)
	require.NotEmpty(t, accounts)

	createEntryArg := CreateEntryParams{
		AccountID: accounts.ID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), createEntryArg)

	require.Empty(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, entry.AccountID, accounts.ID)

	gottenEntry, err := testQueries.GetEntry(context.Background(), entry.ID)

	require.Empty(t, err)
	require.NotEmpty(t, gottenEntry)
	require.Equal(t, entry.Amount, gottenEntry.Amount)
	require.Equal(t, entry.AccountID, gottenEntry.AccountID)
}

func TestListEntries(t *testing.T) {

	user := createRandomUser(t)

	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	//We used the testify library to check the test result
	accounts, err := testQueries.CreateAccount(context.Background(), arg)

	//require keyword will end the test if the assetion doesn't pass
	//assert keyword will log the failed assertion but will still continue the test
	assert.Nil(t, err)
	assert.NotNil(t, accounts)
	require.NotEmpty(t, accounts)

	for i := 0; i < 10; i++ {
		createEntryArg := CreateEntryParams{
			AccountID: accounts.ID,
			Amount:    util.RandomMoney(),
		}

		testQueries.CreateEntry(context.Background(), createEntryArg)
	}

	listEntryArg := ListEntriesParams{
		AccountID: accounts.ID,
		Limit:     5,
		Offset:    5,
	}

	listEntries, err := testQueries.ListEntries(context.Background(), listEntryArg)
	require.Empty(t, err)
	require.NotEmpty(t, listEntries)
	require.Len(t, listEntries, 5)

	for _, entry := range listEntries {
		require.NotEmpty(t, entry)
	}

}
