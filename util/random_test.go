package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomInt(t *testing.T) {
	var min, max int64
	min = 1
	max = 100
	result := RandomInt(min, max)

	require.Greater(t, result, min)
	require.Less(t, result, max)
}

func TestRandomString(t *testing.T) {
	result := RandomString(10)
	require.Equal(t, 10, len(result))
}

func TestRandomAmount(t *testing.T) {
	result := RandomAmount()
	require.Greater(t, result, int64(0))
	require.Less(t, result, int64(1000))
}

func TestRandomUnit(t *testing.T) {
	result := RandomUnit()
	require.True(t, IsSupportedUnit(result))
}

func TestRandomOwner(t *testing.T) {
	result := RandomOwner()
	require.Equal(t, 6, len(result))
}

func TestRandomEmail(t *testing.T) {
	result := RandomEmail()

	require.Equal(t, result[len(result)-10:], "@email.com")
}
