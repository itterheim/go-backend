package locations

import (
	"time"
)

type Place struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Note      string    `json:"note,omitempty"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Radius    float64   `json:"radius"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
}

type CreatePlaceRequest struct {
	Name      string  `json:"name"`
	Note      string  `json:"note,omitempty"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Radius    float64 `json:"radius"`
}

type PlaceResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Note      string    `json:"note,omitempty"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Radius    float64   `json:"radius"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
}
