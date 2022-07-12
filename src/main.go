package main

import (
	"fmt"
	"lugdroid/trailPlanner/webapp/src/gpx"
)

func main() {
	gpxData := gpx.ReadFile()
	routeData := gpx.ParseData(gpxData)

	fmt.Println("Route name: ", routeData.Name)
	fmt.Println("Number of GPS points: ", len(routeData.Points))
	fmt.Printf("Distance: %.2fkm\n", routeData.Distance)
	fmt.Printf("Ascent: %+.2fm\n", routeData.Ascent)
	fmt.Printf("Descent: %.2fm\n", routeData.Descent)
	fmt.Printf("Min. Elevation: %.1fm\n", routeData.MinElev)
	fmt.Printf("Max. Elevation: %.1fm\n", routeData.MaxElev)
	fmt.Println()

	for _, rp := range routeData.Points {
		fmt.Printf("Elevation: %7.2fm\t Change: %+6.2fm\t Acc.Ascent: %+6.2fm\t Acc.Descent: %+6.2fm\t Acc. Distance: %7.2fkm\t Distance to Prev.: %8.3fkm", rp.Elevation, rp.ElevationChange, rp.AccumulatedAscent, rp.AccumulatedDescent, rp.AccumulatedDistance, rp.DistanceToPrev)
		fmt.Println()
	}
	fmt.Println()

	for i, c := range routeData.Climbs {
		fmt.Printf("Climb #%v %+.2fm - %.2fkm > %.2fkm \n", i+1, (c.EndElev - c.StartElev), c.StartKm, c.EndKm)
	}

	// fmt.Println(calculateDistance(0, 0, 0, 0))
	// fmt.Println(calculateDistance(51.5, 0, 38.8, -77.1))
	// fmt.Println(calculateDistance(40.281897, -5.897834, 40.281249, -5.889513))
	// fmt.Println(fix3dDistance(calculateDistance(40.281897, -5.897834, 40.281249, -5.889513), 10))
}
