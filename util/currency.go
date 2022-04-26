package util

const (
	USD = "USD"
	CAD = "CAD"
	NAR = "NAR"
)

// IsSupportedCurrency returns if a currency is supported or not
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, CAD, NAR:
		return true
	}
	return false
}
