package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsSupportedUnit(t *testing.T) {
	t.Run("supported unit", func(t *testing.T) {
		result := IsSupportedUnit(KWH)
		require.True(t, result)
	})

	t.Run("not supported unit", func(t *testing.T) {
		result := IsSupportedUnit("notSupportedUnit")
		require.False(t, result)
	})
}
