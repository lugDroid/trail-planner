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
