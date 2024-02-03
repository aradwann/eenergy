package token

import (
	"testing"
	"time"

	"github.com/aradwann/eenergy/util"

	"github.com/stretchr/testify/require"
)

func TestPASETOMaker(t *testing.T) {
	maker, err := NewPASETOMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)
	role := util.GeneratorRole

	token, payload, err := maker.CreateToken(username, role, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.Equal(t, role, payload.Role)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)

}

func TestExpiredPASETOToken(t *testing.T) {
	maker, err := NewPASETOMaker(util.RandomString(32))
	require.NoError(t, err)

	role := util.GeneratorRole

	token, payload, err := maker.CreateToken(util.RandomOwner(), role, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)

}

func TestPASETOMakerInvalidKeySize(t *testing.T) {
	symmetricKey := "shortKey" // Invalid key size

	_, err := NewPASETOMaker(symmetricKey)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid key size")
}
