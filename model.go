package main

import "time"

type location struct {
	Lat  float32   `json:"lat"`
	Lng  float32   `json:"lng"`
	Time time.Time `json:"time"`
}
