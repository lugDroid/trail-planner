package main

import (
	"database/sql"
	"log"
	"lugdroid/trailPlanner/webapp/controller"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	db := connectToDatabase()
	defer db.Close()

	var routesController controller.Routes
	var uploadController controller.Upload

	routesController.RegisterRoutes()
	uploadController.RegisterRoutes()

	http.ListenAndServe(":3000", nil)
}

func connectToDatabase() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("TP_DATABASE_URL"))
	if err != nil {
		log.Fatal("Unable to connect to database: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Connection to database failed: ", err)
	}

	return db
}
