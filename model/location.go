package model

import "time"

type Location struct {
	Id   int32
	Name string
	Time time.Time
	Lat  float32
	Lng  float32
}
