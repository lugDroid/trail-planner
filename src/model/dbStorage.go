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
	`, nr.Name, nr.Ascent, nr.Descent, nr.MinElev, nr.MaxElev).Scan((&nr.Id))
	if err != nil {
		fmt.Println("AddRoute query failed - ", err)
	}

	return nr
}
