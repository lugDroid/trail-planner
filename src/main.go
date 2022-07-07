package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"lugdroid/trailPlanner/webapp/src/model"
	"os"
	"strconv"
)

func main() {
	gpxData := readFile()
	routeData := parseGpxData(gpxData)

	fmt.Println("Route name: ", routeData.Name)
	fmt.Println("Number of points: ", len(routeData.Points))

	for _, rp := range routeData.Points {
		fmt.Printf("Elevation: %7.2fm - Change: %+6.2fm", rp.Elevation, rp.ElevationChange)
		fmt.Println()
	}
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
		routeData.Points = append(routeData.Points, rp)
	}

	// second pass calculates data between points
	for i := 0; i < len(routeData.Points); i++ {
		if i == 0 {
			routeData.Points[i].ElevationChange = 0
			continue
		}

		routeData.Points[i].ElevationChange = routeData.Points[i].Elevation - routeData.Points[i-1].Elevation
	}

	return routeData
}

func check(e error) {
	if e != nil {
		fmt.Println(e)
	}
}
