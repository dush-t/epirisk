package util

// HealthStatusValueIsAllowed is a utility function to check if a particular
// value of risk given by user is allowed.
func HealthStatusValueIsAllowed(risk float64) bool {
	allowedValues := map[float64]bool{
		0.0: true, // perfectly fine
		0.9: true, // has symptoms but unconfirmed
		1.0: true, // confirmed case
	}
	return (allowedValues[risk] == true)
}

// RiskShouldBeRecalculated decides whether to update the risk of surrounding
// nodes when user updates his health status
func RiskShouldBeRecalculated(prevHealthStatus float64, newHealthStatus float64) bool {
	return newHealthStatus > prevHealthStatus
}
