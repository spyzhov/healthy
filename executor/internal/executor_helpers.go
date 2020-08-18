package internal

func FloatInSlice(value float64, slice []float64) bool {
	for _, i := range slice {
		if i == value {
			return true
		}
	}
	return false
}
