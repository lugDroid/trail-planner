package controller

import (
	"encoding/json"
	"fmt"
	"lugdroid/trailPlanner/webapp/src/model"
	"net/http"
)

type Routes struct {
	storage model.DbStorage
}

func (r Routes) RegisterRoutes() {
	http.HandleFunc("/routes", r.handleRoutes)
}

func (r Routes) handleRoutes(w http.ResponseWriter, rq *http.Request) {
	routeData := r.storage.GetAllRoutes()

	json, err := json.Marshal(routeData)
	if err != nil {
		fmt.Println("Could not marshal into json object", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
