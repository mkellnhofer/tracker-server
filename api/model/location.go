package model

import "time"

type Location struct {
	Id          int64     `json:"id"`
	ChangeTime  int64     `json:"changeTime"`
	Name        string    `json:"name"`
	Time        time.Time `json:"time"`
	Lat         float32   `json:"lat"`
	Lng         float32   `json:"lng"`
	Description string    `json:"description"`
	Persons     []*Person `json:"persons"`
}
