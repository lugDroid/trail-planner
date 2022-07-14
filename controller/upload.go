package controller

import (
	"encoding/json"
	"fmt"
	"lugdroid/trailPlanner/webapp/src/gpx"
	"net/http"
)

type Upload struct {
}

func (u Upload) RegisterRoutes() {
	http.HandleFunc("/upload", u.handleUpload)
}

func (u Upload) handleUpload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	//var buf bytes.Buffer
	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error retrieving file", err)
	}

	defer file.Close()
	//name := strings.Split(header.Filename, ".")
	fmt.Printf("File name %s\n", header.Filename)

	gpxData := gpx.ReadFile(file)
	routeData := gpx.ParseData(gpxData)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	jsonResp, err := json.Marshal(routeData)
	if err != nil {
		fmt.Println("Could not marshal response into json object", err)
	}

	w.Write(jsonResp)
}
