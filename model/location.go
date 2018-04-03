package model

import "time"

type Location struct {
	Time time.Time `json:"time"`
	Lat  float32   `json:"lat"`
	Lng  float32   `json:"lng"`
}
