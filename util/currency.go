package util

var supportedCurrenciesMap = map[string]struct{}{"USD": {}, "EUR": {}, "CAD": {}} // mimics a set (which is not supported by Go) - maps string to a struct{} which takes up no memory

func IsSupportedCurrency(currency string) bool {
	_, found := supportedCurrenciesMap[currency]
	return found
}
