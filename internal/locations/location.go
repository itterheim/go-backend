package locations

import (
	"backend/internal/core"
	"time"
)

type GpsHistory struct {
	ID        int64     `json:"id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Accuracy  float64   `json:"accuracy"`
	Created   time.Time `json:"created"`
}

type CreateGpsHistoryRequest struct {
	Timestamp  time.Time `json:"timestmap"`
	Tags       []string  `json:"tags"`
	Note       *string   `json:"note"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	Accuracy   float64   `json:"accuracy"`
	UserID     int64
	ProviderID *int64
}

type GpsHistoryResponse struct {
	Event     *core.Event          `json:"event"`
	Action    *core.ActionResponse `json:"action"`
	Latitude  float64              `json:"latitude"`
	Longitude float64              `json:"longitude"`
	Accuracy  float64              `json:"accuracy"`
	Created   time.Time            `json:"created"`
}
