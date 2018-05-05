package model

import "time"

type Location struct {
	Id      int64     `json:"id"`
	Name    string    `json:"name"`
	Time    time.Time `json:"time"`
	Lat     float32   `json:"lat"`
	Lng     float32   `json:"lng"`
	Persons []*Person `json:"persons"`
}
