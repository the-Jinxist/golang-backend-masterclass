package util

//Constants for all supported variables
const (
	USD = "USD"
	EUR = "EUR"
	NGN = "NGN"
	GBP = "GBP"
)

//function that checks if the currency used in the transaction is currently supported
func IsCurrencySupported(currency string) bool {
	switch currency {
	case USD, EUR, NGN, GBP:
		return true
	}

	return false
}
