package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"lugdroid/trailPlanner/webapp/src/model"
	"math"
	"os"
	"strconv"
)

func main() {
	gpxData := readFile()
	routeData := parseGpxData(gpxData)

	fmt.Println("Route name: ", routeData.Name)
	fmt.Println("Number of points: ", len(routeData.Points))

	for _, rp := range routeData.Points {
		fmt.Printf("Elevation: %7.2fm - Change: %+6.2fm \t Acc. Distance: %7.2fkm \t Distance to Prev.: %7.2fkm", rp.Elevation, rp.ElevationChange, rp.AccumulatedDistance, rp.DistanceToPrev)
		fmt.Println()
	}

	fmt.Println(calculateDistance(0, 0, 0, 0))
	fmt.Println(calculateDistance(51.5, 0, 38.8, -77.1))
}

func readFile() model.Gpx {
	gpxFile, error := os.Open("../example.gpx")
	check(error)

	fmt.Println("GPX file successfully opened")
	defer gpxFile.Close()

	bytes, err := ioutil.ReadAll(gpxFile)
	check(err)

	var gpxData model.Gpx
	xml.Unmarshal(bytes, &gpxData)

	return gpxData
}

func parseGpxData(gpxData model.Gpx) model.Route {
	var routeData model.Route

	routeData.Name = gpxData.Rte.Name

	// first pass will add individual points data
	for _, p := range gpxData.Rte.Rtept {
		var rp model.Point

		eleValue, err := strconv.ParseFloat(p.Ele, 32)
		if err != nil {
			fmt.Println("Error while parsing point elevation: ", err)
		}

		rp.Elevation = eleValue

		latValue, err := strconv.ParseFloat(p.Lat, 32)
		if err != nil {
			fmt.Println("Error while parsing point latitude: ", err)
		}
		rp.Coordinates.Lat = latValue

		lonValue, err := strconv.ParseFloat(p.Lon, 32)
		if err != nil {
			fmt.Println("Error while parsing point latitude: ", err)
		}
		rp.Coordinates.Lon = lonValue

		routeData.Points = append(routeData.Points, rp)
	}

	// second pass calculates data between points
	for i := 0; i < len(routeData.Points); i++ {
		currentPoint := routeData.Points[i]

		if i == 0 {
			currentPoint.ElevationChange = 0
			currentPoint.DistanceToPrev = 0
			currentPoint.AccumulatedDistance = 0
			continue
		}

		prevPoint := routeData.Points[i-1]

		currentPoint.ElevationChange = currentPoint.Elevation - prevPoint.Elevation
		currentPoint.DistanceToPrev = calculateDistance(currentPoint.Coordinates.Lat, currentPoint.Coordinates.Lon, prevPoint.Coordinates.Lat, prevPoint.Coordinates.Lon)
		currentPoint.AccumulatedDistance = prevPoint.AccumulatedDistance + currentPoint.DistanceToPrev

		routeData.Points[i] = currentPoint
		routeData.Points[i-1] = prevPoint
	}

	return routeData
}

func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	earthRadiusKm := 6371.0

	// convert coordinates to radians and calculate increments
	lat1 = degToRad(lat1)
	lat2 = degToRad(lat2)
	dLat := degToRad(lat2 - lat1)
	dLon := degToRad(lon2 - lon1)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}

func degToRad(deg float64) float64 {
	return deg * math.Pi / 180
}

func check(e error) {
	if e != nil {
		fmt.Println(e)
	}
}
