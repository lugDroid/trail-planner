package controller

import (
	"encoding/json"
	"fmt"
	"lugdroid/trailPlanner/webapp/src/model"
	"net/http"
	"regexp"
	"strconv"
)

type Routes struct {
	storage model.DbStorage
}

func (r Routes) RegisterRoutes() {
	http.HandleFunc("/routes", r.handleRoutes)
	http.HandleFunc("/routes/", r.handleRoutes)
}

func (r Routes) handleRoutes(w http.ResponseWriter, rq *http.Request) {
	idPattern, _ := regexp.Compile(`/routes/(\d+)`)
	idMatches := idPattern.FindStringSubmatch(rq.URL.Path)
	if len(idMatches) > 0 {
		routeId, _ := strconv.Atoi(idMatches[1])
		r.handleDetail(w, rq, routeId)
		return
	}

	routeData := r.storage.GetAllRoutes()
	json, err := json.Marshal(routeData)
	if err != nil {
		fmt.Println("Could not marshal into json object", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (r Routes) handleDetail(w http.ResponseWriter, rq *http.Request, routeId int) {
	route := r.storage.GetRouteById(routeId)
	json, err := json.Marshal(route)
	if err != nil {
		fmt.Println("Could not marshal into json object", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
