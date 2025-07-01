package utils

func ComputeValue(start, end int, percentComplete float64) int {
	return start + int(float64(end-start)*percentComplete)
}

func GetFahrenheit(celsius float64) float64 {
	return (celsius * 9 / 5) + 32
}
