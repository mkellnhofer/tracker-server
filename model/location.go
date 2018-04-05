package model

import "time"

type Location struct {
	Id      int64
	Name    string
	Time    time.Time
	Lat     float32
	Lng     float32
	Persons []*Person
}
