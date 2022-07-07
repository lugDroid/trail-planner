package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"lugdroid/trailPlanner/webapp/src/model"
	"os"
)

func main() {
	readFile()
}

func readFile() {
	gpxFile, error := os.Open("../example.gpx")
	check(error)

	fmt.Println("GPX file successfully opened")
	defer gpxFile.Close()

	bytes, err := ioutil.ReadAll(gpxFile)
	check(err)

	var gpxData model.Gpx
	xml.Unmarshal(bytes, &gpxData)

	fmt.Println(gpxData.Rte.Name)
	fmt.Println(len(gpxData.Rte.Rtept))
}

func check(e error) {
	if e != nil {
		fmt.Println(e)
	}
}
