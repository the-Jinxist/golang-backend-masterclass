package backend_masterclass

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	arg := CreateAccountParams{
		Owner:    "Favour",
		Balance:  20000,
		Currency: "NGN",
	}

	//We used the testify library to check the test result
	accounts, err := testQueries.CreateAccount(context.Background(), arg)

	//require keyword will end the test if the assetion doesn't pass
	//assert keyword will log the failed assertion but will still continue the test
	assert.Nil(t, err)
	assert.NotNil(t, accounts)
	require.NotEmpty(t, accounts)

	assert.Equal(t, accounts.Balance, int64(20000))
	assert.Equal(t, accounts.Currency, "NGN")
	assert.Equal(t, accounts.Owner, "Favour")

	require.NotZero(t, accounts.ID)
	require.NotZero(t, accounts.CreatedAt)
}
