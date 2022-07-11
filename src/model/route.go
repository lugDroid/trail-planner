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
	AccumulatedAscent   float64
	AccumulatedDescent  float64
	Coordinates         Coord
}

type Coord struct {
	Lat float64
	Lon float64
}
