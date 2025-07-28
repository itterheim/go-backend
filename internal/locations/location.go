package locations

import (
	"time"
)

type GpsHistory struct {
	ID         int64     `json:"id"`
	Timestamp  time.Time `json:"timestamp"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	Accuracy   float64   `json:"accuracy"`
	ProviderID *int64    `json:"-"`
	UserID     int64     `json:"-"`
}

type CreateGpsHistoryRequest struct {
	Timestamp  time.Time `json:"timestamp"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	Accuracy   float64   `json:"accuracy"`
	ProviderID *int64    `json:"-"`
	UserID     int64     `json:"-"`
}

type GpsHistoryResponse struct {
	ID        int64     `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Accuracy  float64   `json:"accuracy"`
}
