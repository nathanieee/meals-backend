package utmath

import "math"

func Round(val float64, roundOn float64, places int) (newVal float64) {
	var (
		round float64
	)

	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func BMICalculation(weight, height float64) float64 {
	height = height / 100

	return weight / (height * height)
}
