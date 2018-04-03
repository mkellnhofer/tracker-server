package model

import "time"

type Location struct {
	Name string    `json:"name"`
	Time time.Time `json:"time"`
	Lat  float32   `json:"lat"`
	Lng  float32   `json:"lng"`
}
