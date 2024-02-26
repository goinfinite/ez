package infraHelper

import (
	"math"
)

func RoundFloat(val float64) float64 {
	ratio := math.Pow(10, float64(1))
	return math.Ceil(val*ratio) / ratio
}
