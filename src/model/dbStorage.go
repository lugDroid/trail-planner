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

func (s *DbStorage) GetAllRoutes() []Route {
	rows, err := s.db.Query(`
		SELECT id, name, ascent, descent, min_elev, max_elev
		FROM route
	`)
	if err != nil {
		fmt.Println("GetAllRoutes query failed", err)
	}
	defer rows.Close()

	var routes []Route
	for rows.Next() {
		r := Route{}
		err := rows.Scan(&r.Id, &r.Name, &r.Ascent, &r.Descent, &r.MinElev, &r.MaxElev)
		if err != nil {
			fmt.Println("Error scanning GetAllRoutes query results", err)
		}

		routes = append(routes, r)
	}

	for i := range routes {
		s.getRoutePoints(&routes[i])
		s.getRouteClimbs(&routes[i])
	}

	return routes
}

func (s *DbStorage) getRoutePoints(nr *Route) {
	rows, err := s.db.Query(`
		SELECT id, distance_to_prev, acc_distance, elevation, elevation_change, acc_ascent, acc_descent, lat, lon
		FROM point
		WHERE route_id = $1
	`, nr.Id)
	if err != nil {
		fmt.Println("GetRoutePoints query failed", err)
	}
	defer rows.Close()

	var points []Point
	for rows.Next() {
		p := Point{}
		err := rows.Scan(&p.Id, &p.DistanceToPrev, &p.AccumulatedDistance, &p.Elevation, &p.ElevationChange, &p.AccumulatedAscent, &p.AccumulatedDescent, &p.Coordinates.Lat, &p.Coordinates.Lon)
		if err != nil {
			fmt.Println("Error scanning GetAllPoints query results", err)
		}

		points = append(points, p)
	}

	nr.Points = points
}

func (s *DbStorage) getRouteClimbs(nr *Route) {
	rows, err := s.db.Query(`
		SELECT id, start_km, end_km, start_elev, end_elev
		FROM climb
		WHERE route_id = $1
	`, nr.Id)
	if err != nil {
		fmt.Println("GetRouteClimbs query failed", err)
	}
	defer rows.Close()

	var climbs []Climb
	for rows.Next() {
		c := Climb{}
		err := rows.Scan(&c.Id, &c.StartKm, &c.EndKm, &c.StartElev, &c.EndElev)
		if err != nil {
			fmt.Println("Error scanning GetAllPoints query results", err)
		}

		climbs = append(climbs, c)
	}

	nr.Climbs = climbs
}

func (s *DbStorage) GetRouteById(routeId int) Route {
	r := Route{}

	row := s.db.QueryRow(`
		SELECT id, name, ascent, descent, min_elev, max_elev
		FROM route
		WHERE id = $1
	`, routeId)
	err := row.Scan(&r.Id, &r.Name, &r.Ascent, &r.Descent, &r.MinElev, &r.MaxElev)
	if err != nil {
		fmt.Println("GetRouteById query failed", err)
	}

	s.getRoutePoints(&r)
	s.getRouteClimbs(&r)

	return r
}
