package gpx

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"lugdroid/trailPlanner/webapp/src/model"
	"math"
	"strconv"
)

func ReadFile(file io.Reader) (model.Gpx, error) {
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return model.Gpx{}, errors.New("could not read file")
	}

	var gpxData model.Gpx
	err = xml.Unmarshal(bytes, &gpxData)
	if err != nil {
		return model.Gpx{}, errors.New("could not unmarshal file content")
	}

	return gpxData, nil
}

func ParseData(gpxData model.Gpx) model.Route {
	var routeData model.Route

	// first pass will add individual points data
	if len(gpxData.Rte.Rtept) > 0 {
		routeData.Points = extractPointData(gpxData.Rte.Rtept)
		routeData.Name = gpxData.Rte.Name
	} else {
		routeData.Points = extractPointData(gpxData.Trk.Trkseg.Trkpt)
		routeData.Name = gpxData.Trk.Name
	}

	// second pass calculates data between points
	calculatePointsData(&routeData.Points)

	routeData.MinElev = getMinElev(&routeData.Points)
	routeData.MaxElev = getMaxElev(&routeData.Points)
	routeData.Distance = routeData.Points[len(routeData.Points)-1].AccumulatedDistance
	routeData.Ascent = routeData.Points[len(routeData.Points)-1].AccumulatedAscent
	routeData.Descent = routeData.Points[len(routeData.Points)-1].AccumulatedDescent

	routeData.Climbs = indentifyClimbs(&routeData.Points)

	return routeData
}

func extractPointData(gpxPoints []model.GpxPoint) []model.Point {
	var points []model.Point

	for _, p := range gpxPoints {
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

		points = append(points, rp)
	}

	return points
}

func calculatePointsData(points *[]model.Point) {
	for i := 0; i < len(*points); i++ {
		currentPoint := (*points)[i]

		if i == 0 {
			currentPoint.ElevationChange = 0
			currentPoint.DistanceToPrev = 0
			currentPoint.AccumulatedDistance = 0
			continue
		}

		prevPoint := (*points)[i-1]

		currentPoint.ElevationChange = currentPoint.Elevation - prevPoint.Elevation

		flatDistance := calculateDistance(currentPoint.Coordinates.Lat, currentPoint.Coordinates.Lon, prevPoint.Coordinates.Lat, prevPoint.Coordinates.Lon)
		currentPoint.DistanceToPrev = fix3dDistance(flatDistance, currentPoint.ElevationChange)
		currentPoint.AccumulatedDistance = prevPoint.AccumulatedDistance + currentPoint.DistanceToPrev
		if currentPoint.ElevationChange > 0 {
			currentPoint.AccumulatedAscent = prevPoint.AccumulatedAscent + currentPoint.ElevationChange
			currentPoint.AccumulatedDescent = prevPoint.AccumulatedDescent
		} else {
			currentPoint.AccumulatedAscent = prevPoint.AccumulatedAscent
			currentPoint.AccumulatedDescent = prevPoint.AccumulatedDescent + currentPoint.ElevationChange
		}

		(*points)[i] = currentPoint
		(*points)[i-1] = prevPoint
	}
}

func indentifyClimbs(points *[]model.Point) []model.Climb {
	var climbs []model.Climb

	initialPoint := (*points)[0]

	for i := 2; i < len(*points); i++ {
		// a climb start/ends when there´s a change in the elevation change sign
		if math.Signbit((*points)[i].ElevationChange) != math.Signbit((*points)[i-1].ElevationChange) {
			endPoint := (*points)[i-1]

			newClimb := model.Climb{
				StartKm:   initialPoint.AccumulatedDistance,
				EndKm:     endPoint.AccumulatedDistance,
				StartElev: initialPoint.Elevation,
				EndElev:   endPoint.Elevation,
			}

			initialPoint = endPoint

			climbs = append(climbs, newClimb)
		}
	}

	// add final climb
	finalClimb := model.Climb{
		StartKm:   initialPoint.AccumulatedDistance,
		EndKm:     (*points)[len(*points)-1].AccumulatedDistance,
		StartElev: initialPoint.Elevation,
		EndElev:   (*points)[len(*points)-1].Elevation,
	}
	climbs = append(climbs, finalClimb)

	return climbs
}

func getMinElev(points *[]model.Point) float64 {
	minElev := (*points)[0].Elevation

	for _, p := range *points {
		if p.Elevation < minElev {
			minElev = p.Elevation
		}
	}

	return minElev
}

func getMaxElev(points *[]model.Point) float64 {
	maxElev := (*points)[0].Elevation

	for _, p := range *points {
		if p.Elevation > maxElev {
			maxElev = p.Elevation
		}
	}

	return maxElev
}

func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	earthRadiusKm := 6371.0

	// convert coordinates to radians and calculate increments
	rLat1 := degToRad(lat1)
	rLat2 := degToRad(lat2)
	rDiffLat := degToRad(lat2 - lat1)
	rDiffLon := degToRad(lon2 - lon1)

	a := math.Sin(rDiffLat/2.0)*math.Sin(rDiffLat/2.0) + math.Sin(rDiffLon/2.0)*math.Sin(rDiffLon/2.0)*math.Cos(rLat1)*math.Cos(rLat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}

func fix3dDistance(distance, elevation float64) float64 {
	// convert distance to meters
	distance = distance * 1000

	// return kilometers
	return math.Sqrt(distance*distance+elevation*elevation) / 1000
}
