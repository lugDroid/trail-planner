package gpx

import (
	"fmt"
	"math"
)

func check(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

func degToRad(deg float64) float64 {
	return deg * math.Pi / 180
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))

	return float64(math.Round(num*output)) / output
}
