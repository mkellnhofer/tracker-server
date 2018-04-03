package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"kellnhofer.com/tracker/data"
	"kellnhofer.com/tracker/model"
)

type LocationRepo struct {
	db *sql.DB
}

func NewLocationRepo(db *sql.DB) *LocationRepo {
	return &LocationRepo{db}
}

// --- Public methods ---

func (r LocationRepo) GetLocations() ([]*model.Location, error) {
	rows, err := r.db.Query("SELECT time, lat, lng FROM location")
	if err != nil {
		log.Print(err)
		e := fmt.Sprintf("Failed to query locations! (%s)", err)
		return nil, errors.New(e)
	}
	defer rows.Close()

	locs, err := r.scanLocationRows(rows)
	if err != nil {
		return nil, err
	}

	return locs, nil
}

func (r LocationRepo) AddLocation(loc *model.Location) error {
	time := data.FormatTime(loc.Time)
	lat := loc.Lat
	lng := loc.Lng

	_, err := r.db.Exec("INSERT INTO location (time, lat, lng) VALUES (?, ?, ?)", time, lat, lng)
	if err != nil {
		log.Print(err)
		e := fmt.Sprintf("Failed to insert location! (%s)", err)
		return errors.New(e)
	}

	return nil
}

// --- Private methods ---

func (r LocationRepo) scanLocationRows(rows *sql.Rows) ([]*model.Location, error) {
	var locs []*model.Location
	for rows.Next() {
		var time string
		var lat float32
		var lng float32

		err := rows.Scan(&time, &lat, &lng)
		if err != nil {
			log.Print(err)
			e := fmt.Sprintf("Failed to query locations! (%s)", err)
			return nil, errors.New(e)
		}

		loc := &model.Location{data.ParseTime(time), lat, lng}
		locs = append(locs, loc)
	}

	if err := rows.Err(); err != nil {
		log.Print(err)
		e := fmt.Sprintf("Failed to query locations! (%s)", err)
		return nil, errors.New(e)
	}

	return locs, nil
}
