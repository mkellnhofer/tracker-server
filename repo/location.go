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
	rows, err := r.db.Query("SELECT id, name, time, lat, lng FROM location ORDER BY time ASC")
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

	for _, loc := range locs {
		pers, err := r.GetLocationPersons(loc.Id)
		if err != nil {
			return nil, err
		}
		loc.Persons = pers
	}

	return locs, nil
}

func (r LocationRepo) GetLocationPersons(id int64) ([]*model.Person, error) {
	rows, err := r.db.Query("SELECT p.id, p.first_name, p.last_name FROM location_person lp "+
		"INNER JOIN person p ON lp.person_id = p.id "+
		"WHERE lp.location_id = ?", id)
	if err != nil {
		log.Print(err)
		e := fmt.Sprintf("Failed to query location persons! (%s)", err)
		return nil, errors.New(e)
	}
	defer rows.Close()

	pers, err := r.scanLocationPersonRows(rows)
	if err != nil {
		return nil, err
	}

	return pers, nil
}

func (r LocationRepo) AddLocation(loc *model.Location) error {
	name := loc.Name
	time := data.FormatTime(loc.Time)
	lat := loc.Lat
	lng := loc.Lng

	res, err := r.db.Exec("INSERT INTO location (name, time, lat, lng) VALUES (?, ?, ?, ?)", name,
		time, lat, lng)
	if err != nil {
		log.Print(err)
		e := fmt.Sprintf("Failed to insert location! (%s)", err)
		return errors.New(e)
	}

	locId, err := res.LastInsertId()
	if err != nil {
		log.Print(err)
		e := fmt.Sprintf("Failed to insert location! (%s)", err)
		return errors.New(e)
	}

	for _, per := range loc.Persons {
		perId, err := r.GetPersonId(per.FirstName, per.LastName)
		if err != nil {
			return err
		}

		if perId == 0 {
			perId, err = r.CreatePerson(per.FirstName, per.LastName)
			if err != nil {
				return err
			}
		}

		_, err = r.db.Exec("INSERT INTO location_person (location_id, person_id) VALUES (?, ?)",
			locId, perId)
		if err != nil {
			log.Print(err)
			e := fmt.Sprintf("Failed to insert location! (%s)", err)
			return errors.New(e)
		}
	}

	return nil
}

func (r LocationRepo) GetPersonId(firstName string, lastName string) (int64, error) {
	row := r.db.QueryRow("SELECT id FROM person WHERE first_name LIKE ? AND last_name LIKE ?",
		firstName, lastName)

	var perId sql.NullInt64

	err := row.Scan(&perId)
	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		log.Print(err)
		e := fmt.Sprintf("Failed to query person ID! (%s)", err)
		return 0, errors.New(e)
	}

	if perId.Valid {
		return perId.Int64, nil
	} else {
		return 0, nil
	}
}

func (r LocationRepo) CreatePerson(firstName string, lastName string) (int64, error) {
	res, err := r.db.Exec("INSERT INTO person(first_name, last_name) VALUES(?, ?)", firstName,
		lastName)
	if err != nil {
		log.Print(err)
		e := fmt.Sprintf("Failed to create person! (%s)", err)
		return 0, errors.New(e)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Print(err)
		e := fmt.Sprintf("Failed to create person! (%s)", err)
		return 0, errors.New(e)
	}

	return id, nil
}

// --- Private methods ---

func (r LocationRepo) scanLocationRows(rows *sql.Rows) ([]*model.Location, error) {
	var locs []*model.Location
	for rows.Next() {
		var id int64
		var name string
		var time string
		var lat float32
		var lng float32

		err := rows.Scan(&id, &name, &time, &lat, &lng)
		if err != nil {
			log.Print(err)
			e := fmt.Sprintf("Failed to query locations! (%s)", err)
			return nil, errors.New(e)
		}

		loc := &model.Location{id, name, data.ParseTime(time), lat, lng, nil}
		locs = append(locs, loc)
	}

	if err := rows.Err(); err != nil {
		log.Print(err)
		e := fmt.Sprintf("Failed to query locations! (%s)", err)
		return nil, errors.New(e)
	}

	return locs, nil
}

func (r LocationRepo) scanLocationPersonRows(rows *sql.Rows) ([]*model.Person, error) {
	var pers []*model.Person
	for rows.Next() {
		var id int64
		var firstName string
		var lastName string

		err := rows.Scan(&id, &firstName, &lastName)
		if err != nil {
			log.Print(err)
			e := fmt.Sprintf("Failed to query location persons! (%s)", err)
			return nil, errors.New(e)
		}

		per := &model.Person{id, firstName, lastName}
		pers = append(pers, per)
	}

	if err := rows.Err(); err != nil {
		log.Print(err)
		e := fmt.Sprintf("Failed to query location persons! (%s)", err)
		return nil, errors.New(e)
	}

	return pers, nil
}
