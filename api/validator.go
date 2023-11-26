package api

import (
	"github.com/aradwann/eenergy/util"

	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if unit, ok := fl.Field().Interface().(string); ok {
		// check if unit is supported
		return util.IsSupportedUnit(unit)
	}
	return false

}
