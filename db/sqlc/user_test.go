package backend_masterclass

import (
	"backend_masterclass/util"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	createRandomUser(t)

}

func TestGetUser(t *testing.T) {
	randomUser := createRandomUser(t)

	retrievedUser, err := testQueries.GetUser(context.Background(), randomUser.Username)

	require.NoError(t, err)
	require.NotEmpty(t, retrievedUser)

	require.Equal(t, retrievedUser.Username, randomUser.Username)
	require.Equal(t, retrievedUser.FullName, randomUser.FullName)
	require.Equal(t, retrievedUser.HashedPassword, randomUser.HashedPassword)
	require.Equal(t, retrievedUser.Email, randomUser.Email)
}

func createRandomUser(t *testing.T) Users {
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: "secret",
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	assert.Nil(t, err)
	assert.NotNil(t, user)
	require.NotEmpty(t, user)

	assert.Equal(t, user.Email, arg.Email)
	assert.Equal(t, user.HashedPassword, arg.HashedPassword)
	assert.Equal(t, user.FullName, arg.FullName)
	assert.Equal(t, user.Username, arg.Username)

	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())

	return user
}
