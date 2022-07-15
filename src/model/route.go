package model

type Route struct {
	Id       int
	Name     string
	Distance float64
	Ascent   float64
	Descent  float64
	MinElev  float64
	MaxElev  float64
	Points   []Point
	Climbs   []Climb
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

type Climb struct {
	StartKm   float64
	EndKm     float64
	StartElev float64
	EndElev   float64
}
