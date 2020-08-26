package internal

func FloatInSlice(value float64, slice []float64) bool {
	for _, i := range slice {
		if i == value {
			return true
		}
	}
	return false
}

func StringInSlice(value string, slice []string) bool {
	for _, i := range slice {
		if i == value {
			return true
		}
	}
	return false
}
