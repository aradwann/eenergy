package api

import (
	"testing"

	"github.com/aradwann/eenergy/util"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

var unitItems = []struct {
	have string
	want bool
}{
	{util.KWH, true},
	{"1234567890QWERTYUIOP", false},
}

// TestValidateUnit
func TestValidateUnit(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("unit", validUnit)

	for _, item := range unitItems {
		err := validate.Var(item.have, "unit")
		if item.want {
			require.Nil(t, err)
		} else {
			require.Error(t, err)
		}
	}
}
