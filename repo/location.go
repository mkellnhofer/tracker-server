package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

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

func (r LocationRepo) ExistsLocation(id int64) (bool, error) {
	row := r.db.QueryRow("SELECT COUNT(*) FROM location WHERE id = ?", id)

	var n int
	err := row.Scan(&n)
	if err != nil {
		log.Print(err)
		e := fmt.Sprintf("Failed to query location! (%s)", err)
		return false, errors.New(e)
	}

	return n > 0, nil
}

func (r LocationRepo) GetLocations() ([]*model.Location, error) {
	rows, err := r.db.Query("SELECT id, chng_time, name, time, lat, lng FROM location ORDER BY " +
		"time ASC")
	return r.getLocationRows(rows, err)
}

func (r LocationRepo) GetLocationsByChangeTime(ct int64) ([]*model.Location, error) {
	rows, err := r.db.Query("SELECT id, chng_time, name, time, lat, lng FROM location WHERE "+
		"chng_time >= ? ORDER BY time ASC", ct)
	return r.getLocationRows(rows, err)
}

func (r LocationRepo) GetLocation(id int64) (*model.Location, error) {
	row := r.db.QueryRow("SELECT id, chng_time, name, time, lat, lng FROM location WHERE id = ?",
		id)

	loc, err := r.scanLocationRow(row)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		log.Print(err)
		e := fmt.Sprintf("Failed to query location! (%s)", err)
		return nil, errors.New(e)
	default:
	}

	pers, err := r.GetLocationPersons(loc.Id)
	if err != nil {
		return nil, err
	}
	loc.Persons = pers

	return loc, nil
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

func (r LocationRepo) AddLocation(loc *model.Location) (int64, int64, error) {
	ct := time.Now().Unix()
	name := loc.Name
	t := data.FormatTime(loc.Time)
	lat := loc.Lat
	lng := loc.Lng

	res, err := r.db.Exec("INSERT INTO location (chng_time, name, time, lat, lng) VALUES (?, ?, "+
		"?, ?, ?)", ct, name, t, lat, lng)
	if err != nil {
		log.Print(err)
		e := fmt.Sprintf("Failed to insert location! (%s)", err)
		return 0, 0, errors.New(e)
	}

	locId, err := res.LastInsertId()
	if err != nil {
		log.Print(err)
		e := fmt.Sprintf("Failed to insert location! (%s)", err)
		return 0, 0, errors.New(e)
	}

	for _, per := range loc.Persons {
		perId, err := r.GetPersonId(per.FirstName, per.LastName)
		if err != nil {
			return 0, 0, err
		}

		if perId == 0 {
			perId, err = r.CreatePerson(per.FirstName, per.LastName)
			if err != nil {
				return 0, 0, err
			}
		}

		_, err = r.db.Exec("INSERT INTO location_person (location_id, person_id) VALUES (?, ?)",
			locId, perId)
		if err != nil {
			log.Print(err)
			e := fmt.Sprintf("Failed to insert location! (%s)", err)
			return 0, 0, errors.New(e)
		}
	}

	return locId, ct, nil
}

func (r LocationRepo) DeleteLocation(id int64) error {
	_, err := r.db.Exec("DELETE FROM location WHERE id = ?", id)
	if err != nil {
		log.Print(err)
		e := fmt.Sprintf("Failed to delete location! (%s)", err)
		return errors.New(e)
	}

	dt := time.Now().Unix()

	_, err = r.db.Exec("INSERT INTO deleted_location (id, del_time) VALUES (?, ?)", id, dt)
	if err != nil {
		log.Print(err)
		e := fmt.Sprintf("Failed to insert location! (%s)", err)
		return errors.New(e)
	}

	return nil
}

func (r LocationRepo) GetDeletedLocationIdsByDeletionTime(dt int64) ([]int64, error) {
	rows, err := r.db.Query("SELECT id FROM deleted_location WHERE del_time >= ?", dt)
	if err != nil {
		log.Print(err)
		e := fmt.Sprintf("Failed to query deleted locations! (%s)", err)
		return nil, errors.New(e)
	}
	defer rows.Close()

	ids, err := r.scanDeletedLocationRows(rows)
	if err != nil {
		return nil, err
	}

	return ids, nil
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

func (r LocationRepo) getLocationRows(rows *sql.Rows, err error) ([]*model.Location, error) {
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

func (r LocationRepo) scanLocationRows(rows *sql.Rows) ([]*model.Location, error) {
	var locs []*model.Location
	for rows.Next() {
		loc, err := r.scanLocationRow(rows)
		if err != nil {
			log.Print(err)
			e := fmt.Sprintf("Failed to query locations! (%s)", err)
			return nil, errors.New(e)
		}
		locs = append(locs, loc)
	}

	if err := rows.Err(); err != nil {
		log.Print(err)
		e := fmt.Sprintf("Failed to query locations! (%s)", err)
		return nil, errors.New(e)
	}

	return locs, nil
}

func (r LocationRepo) scanLocationRow(scan Scanner) (*model.Location, error) {
	var id int64
	var ct int64
	var name string
	var t string
	var lat float32
	var lng float32

	err := scan.Scan(&id, &ct, &name, &t, &lat, &lng)
	if err != nil {
		return nil, err
	}

	return &model.Location{id, ct, name, data.ParseTime(t), lat, lng, nil}, nil
}

func (r LocationRepo) scanDeletedLocationRows(rows *sql.Rows) ([]int64, error) {
	var ids []int64
	for rows.Next() {
		id, err := r.scanDeletedLocationRow(rows)
		if err != nil {
			log.Print(err)
			e := fmt.Sprintf("Failed to query deleted locations! (%s)", err)
			return nil, errors.New(e)
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		log.Print(err)
		e := fmt.Sprintf("Failed to query deleted locations! (%s)", err)
		return nil, errors.New(e)
	}

	return ids, nil
}

func (r LocationRepo) scanDeletedLocationRow(scan Scanner) (int64, error) {
	var id int64

	err := scan.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
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
