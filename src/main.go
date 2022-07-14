package main

import (
	"lugdroid/trailPlanner/webapp/controller"
	"net/http"
)

func main() {
	var routesController controller.Routes
	var uploadController controller.Upload

	routesController.RegisterRoutes()
	uploadController.RegisterRoutes()

	http.ListenAndServe(":3000", nil)
}
