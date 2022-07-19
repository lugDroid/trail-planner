package controller

import (
	"encoding/json"
	"fmt"
	"lugdroid/trailPlanner/webapp/src/gpx"
	"lugdroid/trailPlanner/webapp/src/model"
	"net/http"
	"regexp"
	"strconv"
)

type Routes struct {
	storage model.DbStorage
}

func (r Routes) StartUp(st *model.DbStorage) {
	r.storage = *st
	r.registerRoutes()
}

func (r Routes) registerRoutes() {
	http.HandleFunc("/routes", r.handleRoutes)
	http.HandleFunc("/routes/", r.handleRoutes)
}

func (r Routes) handleRoutes(w http.ResponseWriter, rq *http.Request) {
	idPattern, _ := regexp.Compile(`/routes/(\d+)`)
	idMatches := idPattern.FindStringSubmatch(rq.URL.Path)
	if len(idMatches) > 0 {
		routeId, _ := strconv.Atoi(idMatches[1])

		switch rq.Method {
		case http.MethodGet:
			r.handleDetail(w, rq, routeId)
		case http.MethodDelete:
			r.handleDelete(w, rq, routeId)
		case http.MethodPut:
			// TO-DO: Modify route to be implemented when needed
		}
	} else {
		switch rq.Method {
		case http.MethodGet:
			r.handleList(w, rq)
		case http.MethodPost:
			r.handleUpload(w, rq)
		}
	}
}

func (r Routes) handleDetail(w http.ResponseWriter, rq *http.Request, routeId int) {
	route, err := r.storage.GetRouteById(routeId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json, err := json.Marshal(route)
	if err != nil {
		fmt.Println("Could not marshal into json object", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (r Routes) handleDelete(w http.ResponseWriter, rq *http.Request, routeId int) {
	route, err := r.storage.DeleteRoute(routeId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json, err := json.Marshal(route)
	if err != nil {
		fmt.Println("Could not marshal into json object", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (r Routes) handleList(w http.ResponseWriter, rq *http.Request) {
	routeData, err := r.storage.GetAllRoutes()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(routeData)
	if err != nil {
		fmt.Println("Could not marshal into json object", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (r Routes) handleUpload(w http.ResponseWriter, rq *http.Request) {
	rq.ParseMultipartForm(32 << 20)
	//var buf bytes.Buffer
	file, header, err := rq.FormFile("file")
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
	route, err := r.storage.AddRoute(routeData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	json, err := json.Marshal(route)
	if err != nil {
		fmt.Println("Could not marshal response into json object", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
