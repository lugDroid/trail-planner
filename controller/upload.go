package controller

import (
	"encoding/json"
	"fmt"
	"lugdroid/trailPlanner/webapp/src/gpx"
	"lugdroid/trailPlanner/webapp/src/model"
	"net/http"
)

type Upload struct {
	storage model.DbStorage
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
	fmt.Printf("Received file name %s\n", header.Filename)

	gpxData, err := gpx.ReadFile(file)
	if err != nil {
		fmt.Println("Error reading file", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	routeData := gpx.ParseData(gpxData)

	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)

	u.storage.AddRoute(routeData)

	w.WriteHeader(http.StatusOK)
	resp["Status"] = "Status OK"
	resp["ReceivedFile"] = header.Filename

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("Could not marshal response into json object", err)
	}

	w.Write(jsonResp)
}
