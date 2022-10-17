package api

import (
	"backend_masterclass/util"

	"github.com/go-playground/validator/v10"
)

//This methods creates a new function on the validator.Func interface that allows us to validate a field
//in our JSON validator
var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	currency, ok := fieldLevel.Field().Interface().(string)
	if ok {
		return util.IsCurrencySupported(currency)
	}

	return false
}
