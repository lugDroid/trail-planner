package controller

import "lugdroid/trailPlanner/webapp/src/model"

var (
	routesController Routes
	uploadController Upload
)

func StartUp(storage *model.DbStorage) {
	routesController.storage = *storage
	uploadController.storage = *storage

	routesController.RegisterRoutes()
	uploadController.RegisterRoutes()
}
