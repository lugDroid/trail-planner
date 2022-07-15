package model

import (
	"database/sql"
	"fmt"
)

type DbStorage struct {
	db *sql.DB
}

func NewDbStorage(database *sql.DB) DbStorage {
	return DbStorage{
		db: database,
	}
}

func (s *DbStorage) AddRoute(nr Route) Route {
	err := s.db.QueryRow(`
		INSERT INTO route (name, ascent, descent, min_elev ,max_elev)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, nr.Name, nr.Ascent, nr.Descent, nr.MinElev, nr.MaxElev).Scan(&nr.Id)
	if err != nil {
		fmt.Println("AddRoute query failed - ", err)
	}

	s.addRoutePoints(&nr)
	s.addRouteClimbs(&nr)

	return nr
	// TO-DO return error codes so handlefunc is able to sent proper response status
}

func (s *DbStorage) addRoutePoints(nr *Route) {
	for i, p := range (*nr).Points {
		err := s.db.QueryRow(`
			INSERT INTO point (route_id, distance_to_prev, acc_distance, elevation, elevation_change, acc_ascent, acc_descent, lat, lon)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING id
		`, nr.Id, p.DistanceToPrev, p.AccumulatedDistance, p.Elevation, p.ElevationChange, p.AccumulatedAscent, p.AccumulatedDescent, p.Coordinates.Lat, p.Coordinates.Lon).Scan(&nr.Points[i].Id)
		if err != nil {
			fmt.Println("AddRoutePoints query failed - ", err)
		}
	}
	// TO-DO return error code
}

func (s *DbStorage) addRouteClimbs(nr *Route) {
	for i, c := range (*nr).Climbs {
		err := s.db.QueryRow(`
			INSERT INTO climb (route_id, start_km, end_km, start_elev, end_elev)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`, nr.Id, c.StartKm, c.EndKm, c.StartElev, c.EndElev).Scan(&nr.Climbs[i].Id)
		if err != nil {
			fmt.Println("AddRouteClimbs query failed - ", err)
		}
	}
	// TO-DO return error code
}
