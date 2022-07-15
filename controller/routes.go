package controller

import (
	"encoding/json"
	"fmt"
	"lugdroid/trailPlanner/webapp/src/gpx"
	"lugdroid/trailPlanner/webapp/src/model"
	"net/http"
	"os"
)

type Routes struct {
	storage model.DbStorage
}

func (r Routes) RegisterRoutes() {
	http.HandleFunc("/routes", r.handleRoutes)
}

func (r Routes) handleRoutes(w http.ResponseWriter, rq *http.Request) {
	gpxFile, err := os.Open("../example003.gpx")
	if err != nil {
		fmt.Println("Could not open file", err)
	}
	defer gpxFile.Close()

	gpxData, err := gpx.ReadFile(gpxFile)
	if err != nil {
		fmt.Println("Error reading file", err)
	}

	routeData := gpx.ParseData(gpxData)

	json, err := json.Marshal(routeData)
	if err != nil {
		fmt.Println("Could not marshal into json object", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
