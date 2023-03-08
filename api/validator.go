package api

import (
	"bank/util"

	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	currency, ok := fieldLevel.Field().Interface().(string) // convert to string
	if !ok {
		return false
	}

	// check if currency is supported
	return util.IsSupportedCurrency(currency)
}
