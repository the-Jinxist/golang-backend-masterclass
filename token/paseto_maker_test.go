package token

import (
	"backend_masterclass/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.WithinDuration(t, payload.ExpiredAt, expiredAt, time.Second)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.Equal(t, payload.Username, username)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	token, payload, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrorExpiredToken.Error())
	require.Nil(t, payload)
}

// func TestInvalidPasetoToken(t *testing.T) {
// 	normalSecretToken := util.RandomString(32)
// 	maker, err := NewPasetoMaker(normalSecretToken)
// 	fakeMaker
// 	require.NoError(t, err)

// 	token, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, token)

// 	payload, err := maker.VerifyToken(token)
// 	require.Error(t, err)
// 	require.EqualError(t, err, ErrorExpiredToken.Error())
// 	require.Nil(t, payload)
// }
