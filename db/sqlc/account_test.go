package backend_masterclass

import (
	"backend_masterclass/util"
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
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

	assert.Equal(t, accounts.Balance, arg.Balance)
	assert.Equal(t, accounts.Currency, arg.Currency)
	assert.Equal(t, accounts.Owner, arg.Owner)

	require.NotZero(t, accounts.ID)
	require.NotZero(t, accounts.CreatedAt)
}

func TestGetAccount(t *testing.T) {
	//Unit tests should be independent from each other, so we'll have to create an account before then getting said account
	user := createRandomUser(t)

	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	//We used the testify library to check the test result
	var err error
	accounts, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	gottenAccount, err := testQueries.GetAccount(context.Background(), accounts.ID)

	require.NoError(t, err)
	require.NotEmpty(t, gottenAccount)

	require.Equal(t, accounts.Owner, gottenAccount.Owner)
	require.Equal(t, accounts.Balance, gottenAccount.Balance)
	require.Equal(t, accounts.Currency, gottenAccount.Currency)
}

func TestUpdateAccount(t *testing.T) {
	user := createRandomUser(t)

	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	var err error
	accounts, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	updateArg := UpdateAccountParams{
		ID:      accounts.ID,
		Balance: util.RandomMoney(),
	}

	updateAccount, err := testQueries.UpdateAccount(context.Background(), updateArg)
	require.NoError(t, err)
	require.NotEmpty(t, updateAccount)

	freshlyUpdatedAccount, err := testQueries.GetAccount(context.Background(), accounts.ID)
	require.NoError(t, err)
	require.NotEmpty(t, freshlyUpdatedAccount)
	require.Equal(t, updateArg.Balance, freshlyUpdatedAccount.Balance)
	require.Equal(t, accounts.Currency, freshlyUpdatedAccount.Currency)
	require.Equal(t, accounts.Owner, freshlyUpdatedAccount.Owner)
}

func TestDeleteAccount(t *testing.T) {
	user := createRandomUser(t)

	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	var err error
	accounts, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	err = testQueries.DeleteAccount(context.Background(), accounts.ID)
	require.NoError(t, err)

	deleteAccount, err := testQueries.GetAccount(context.Background(), accounts.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, deleteAccount)
}

func TestListAccount(t *testing.T) {
	var lastAccount Users
	for i := 0; i < 10; i++ {
		user := createRandomUser(t)
		lastAccount = user

		arg := CreateAccountParams{
			Owner:    user.Username,
			Balance:  util.RandomMoney(),
			Currency: util.RandomCurrency(),
		}

		testQueries.CreateAccount(context.Background(), arg)

	}

	arg := ListAccountsParams{
		Owner:  lastAccount.Username,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, account.Owner, lastAccount.Username)
	}
}
