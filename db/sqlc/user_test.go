package backend_masterclass

import (
	"backend_masterclass/util"
	"context"
	"database/sql"
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

func TestUpdateUserOnlyFullName(t *testing.T) {

	randomUser := createRandomUser(t)
	newFullName := util.RandomString(10)

	updateArg := UpdateUserParams{
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
		Username: randomUser.Username,
	}

	updatedUser, err := testQueries.UpdateUser(context.Background(), updateArg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	require.NotEqual(t, updatedUser.FullName, randomUser.FullName)
	require.Equal(t, updatedUser.FullName, newFullName)
	require.Equal(t, updatedUser.Username, randomUser.Username)
	require.Equal(t, updatedUser.Email, randomUser.Email)

}

func TestUpdateUserOnlyEmail(t *testing.T) {

	randomUser := createRandomUser(t)
	newEmail := util.RandomEmail()

	updateArg := UpdateUserParams{
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		Username: randomUser.Username,
	}

	updatedUser, err := testQueries.UpdateUser(context.Background(), updateArg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	require.NotEqual(t, updatedUser.Email, randomUser.Email)
	require.Equal(t, updatedUser.Email, newEmail)
	require.Equal(t, updatedUser.Username, randomUser.Username)
	require.Equal(t, updatedUser.FullName, randomUser.FullName)

}

func TestUpdateUserOnlyPassword(t *testing.T) {

	randomUser := createRandomUser(t)
	newPassword := util.RandomString(10)

	newHashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)
	require.NotEmpty(t, newHashedPassword)

	updateArg := UpdateUserParams{
		HashedPassword: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
		Username: randomUser.Username,
	}

	updatedUser, err := testQueries.UpdateUser(context.Background(), updateArg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	require.NotEqual(t, updatedUser.HashedPassword, randomUser.HashedPassword)
	require.Equal(t, updatedUser.HashedPassword, newHashedPassword)
	require.Equal(t, updatedUser.Username, randomUser.Username)
	require.Equal(t, updatedUser.FullName, randomUser.FullName)

}

func TestUpdateUserAllFields(t *testing.T) {

	randomUser := createRandomUser(t)
	newPassword := util.RandomString(10)

	newEmail := util.RandomEmail()
	newFullName := util.RandomString(10)
	newHashedPassword, err := util.HashPassword(newPassword)

	require.NoError(t, err)
	require.NotEmpty(t, newHashedPassword)

	updateArg := UpdateUserParams{
		HashedPassword: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		Username: randomUser.Username,
	}

	updatedUser, err := testQueries.UpdateUser(context.Background(), updateArg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	require.NotEqual(t, updatedUser.HashedPassword, randomUser.HashedPassword)
	require.NotEqual(t, updatedUser.FullName, randomUser.FullName)
	require.NotEqual(t, updatedUser.Email, randomUser.Email)

	require.Equal(t, updatedUser.HashedPassword, newHashedPassword)
	require.Equal(t, updatedUser.FullName, newFullName)
	require.Equal(t, updatedUser.Email, newEmail)

}

func createRandomUser(t *testing.T) Users {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	assert.Nil(t, err)
	require.NotEmpty(t, hashedPassword)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
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
