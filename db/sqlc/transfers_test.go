package backend_masterclass

import (
	"backend_masterclass/util"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {

	user := createRandomUser(t)

	account1Params := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account1, err := testQueries.CreateAccount(context.Background(), account1Params)
	assert.Nil(t, err)
	assert.NotNil(t, account1)
	require.NotEmpty(t, account1)

	user2 := createRandomUser(t)

	account2Params := CreateAccountParams{
		Owner:    user2.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account2, err := testQueries.CreateAccount(context.Background(), account2Params)
	assert.Nil(t, err)
	assert.NotNil(t, account2)
	require.NotEmpty(t, account2)

	createTransferArg := CreateTransferParams{
		FromAccount: account1.ID,
		ToAccount:   account2.ID,
		Amout:       util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), createTransferArg)
	assert.Nil(t, err)
	assert.NotNil(t, transfer)
	require.NotEmpty(t, transfer)
	require.Equal(t, transfer.FromAccount, account1.ID)
	require.Equal(t, transfer.ToAccount, account2.ID)
	require.Equal(t, transfer.Amout, createTransferArg.Amout)

}

func TestGetTransfers(t *testing.T) {

	user := createRandomUser(t)

	account1Params := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account1, err := testQueries.CreateAccount(context.Background(), account1Params)
	assert.Nil(t, err)
	assert.NotNil(t, account1)
	require.NotEmpty(t, account1)

	user2 := createRandomUser(t)

	account2Params := CreateAccountParams{
		Owner:    user2.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account2, err := testQueries.CreateAccount(context.Background(), account2Params)
	assert.Nil(t, err)
	assert.NotNil(t, account2)
	require.NotEmpty(t, account2)

	for i := 0; i < 10; i++ {
		createTransferArg := CreateTransferParams{
			FromAccount: account1.ID,
			ToAccount:   account2.ID,
			Amout:       util.RandomMoney(),
		}

		testQueries.CreateTransfer(context.Background(), createTransferArg)
	}
	arg := GetTransfersParams{
		FromAccount: account1.ID,
		ToAccount:   account2.ID,
		Limit:       5,
		Offset:      5,
	}

	accounts, err := testQueries.GetTransfers(context.Background(), arg)
	assert.Nil(t, err)
	require.NotEmpty(t, accounts)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}

}

func TestGetTransfer(t *testing.T) {

	user := createRandomUser(t)

	account1Params := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account1, err := testQueries.CreateAccount(context.Background(), account1Params)
	assert.Nil(t, err)
	assert.NotNil(t, account1)
	require.NotEmpty(t, account1)

	user2 := createRandomUser(t)

	account2Params := CreateAccountParams{
		Owner:    user2.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account2, err := testQueries.CreateAccount(context.Background(), account2Params)
	assert.Nil(t, err)
	assert.NotNil(t, account2)
	require.NotEmpty(t, account2)

	createTransferArg := CreateTransferParams{
		FromAccount: account1.ID,
		ToAccount:   account2.ID,
		Amout:       util.RandomMoney(),
	}

	createdTransfer, err := testQueries.CreateTransfer(context.Background(), createTransferArg)
	assert.Nil(t, err)
	assert.NotNil(t, createdTransfer)
	require.NotEmpty(t, createdTransfer)
	require.Equal(t, createdTransfer.FromAccount, account1.ID)
	require.Equal(t, createdTransfer.ToAccount, account2.ID)
	require.Equal(t, createdTransfer.Amout, createTransferArg.Amout)

	gottenTransfer, err := testQueries.GetTransfer(context.Background(), createdTransfer.ID)
	assert.Nil(t, err)
	assert.NotNil(t, gottenTransfer)
	require.NotEmpty(t, gottenTransfer)
	require.Equal(t, createdTransfer.FromAccount, gottenTransfer.FromAccount)
	require.Equal(t, createdTransfer.ToAccount, gottenTransfer.ToAccount)
	require.Equal(t, createdTransfer.Amout, gottenTransfer.Amout)
}
