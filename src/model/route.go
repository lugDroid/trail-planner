package model

type Route struct {
	Name   string
	Points []Point
}

type Point struct {
	DistanceToPrev      float64
	AccumulatedDistance float64
	Elevation           float64
	ElevationChange     float64
}
