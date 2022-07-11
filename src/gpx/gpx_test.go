package gpx

import (
	"testing"
)

func TestCalcDist(t *testing.T) {
	t.Run("", testCalcDistFunc([]float64{0.0, 0.0, 0.0, 0.0}, 0.0))
	t.Run("", testCalcDistFunc([]float64{40.000000, 5.000000, 40.000000, 5.001000}, 0.0852))
	t.Run("", testCalcDistFunc([]float64{40.000000, -5.000000, 40.000001, -5.000001}, 0.0001))
	t.Run("", testCalcDistFunc([]float64{40.265899, -5.870435, 40.264269, -5.873421}, 0.3115))
}

func testCalcDistFunc(coords []float64, expected float64) func(*testing.T) {
	return func(t *testing.T) {
		result := toFixed(calculateDistance(coords[0], coords[1], coords[2], coords[3]), 4)

		if result != expected {
			t.Errorf("Expected %.4f, got %.4f", expected, result)
		}
	}
}

func TestDegToRad(t *testing.T) {
	t.Run("0.0", testDegToRadFund(0.0, 0.0))
	t.Run("1.0", testDegToRadFund(1.0, 0.01745))
	t.Run("-1.0", testDegToRadFund(-1.0, -0.01745))
	t.Run("5.0", testDegToRadFund(5.0, 0.08727))
	t.Run("-5.0", testDegToRadFund(-5.0, -0.08727))
	t.Run("40.265899", testDegToRadFund(40.265899, 0.70277))
	t.Run("-5.870435", testDegToRadFund(-5.870435, -0.10246))
}

func testDegToRadFund(input float64, expected float64) func(*testing.T) {
	return func(t *testing.T) {
		result := toFixed(degToRad(input), 5)

		if result != expected {
			t.Errorf("Expected %.5f, got %.5f", expected, result)
		}
	}
}
