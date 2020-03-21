package util

// RiskValueIsAllowed is a utility function to check if a particular
// value of risk given by user is allowed.
func RiskValueIsAllowed(risk float64) bool {
	allowedValues := map[float64]bool{
		0.6: true,
		1.0: true,
	}
	return (allowedValues[risk] == true)
}
