package backend_masterclass

import (
	"backend_masterclass/util"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	testStore := NewStore(testDB)

	//require keyword will end the test if the assetion doesn't pass
	//assert keyword will log the failed assertion but will still continue the test
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 5
	amount := int64(10)

	//channels for passing the result from the inner goroutine in the forloop
	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		//separate go routine

		//because we're running this transaction in a separate goroutine, we have to find another way
		//to send it's result to the main goroutine `TestTransferTx`

		//As it stands, using testify inside this inner goroutine will not `definitely` stop the execution
		//of the outer goroutine
		go func() {
			result, err := testStore.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()

	}

	//Checking result here
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		//Checking transfer
		transfer := result.Transfers
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccount)
		require.Equal(t, account2.ID, transfer.ToAccount)
		require.Equal(t, amount, transfer.Amout)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = testStore.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//Check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//Check accounts balance

	}

}

func createRandomAccount(t *testing.T) Accounts {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	accounts, err := testQueries.CreateAccount(context.Background(), arg)

	assert.Nil(t, err)
	assert.NotNil(t, accounts)
	require.NotEmpty(t, accounts)

	assert.Equal(t, accounts.Balance, arg.Balance)
	assert.Equal(t, accounts.Currency, arg.Currency)
	assert.Equal(t, accounts.Owner, arg.Owner)

	require.NotZero(t, accounts.ID)
	require.NotZero(t, accounts.CreatedAt)

	return accounts
}
