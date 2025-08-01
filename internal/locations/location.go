package locations

import (
	"backend/internal/core"
)

type Location struct {
	Latitude  float64
	Longitude float64
	Accuracy  float64
	EventID   int64
}

func (l *Location) ToLocationResponse() *LocationResponse {
	return &LocationResponse{
		Latitude:  l.Latitude,
		Longitude: l.Longitude,
		Accuracy:  l.Accuracy,
	}
}

type LocationEvent struct {
	core.Event
	Extras Location
}

type LocationRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Accuracy  float64 `json:"accuracy"`
}

type CreateLocationEventRequest struct {
	core.CreateEventRequest

	Extras LocationRequest `json:"extras"`
}

type UpdateLocationEventRequest struct {
	core.UpdateEventRequest

	Extras LocationRequest `json:"extras"`
}

type LocationResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Accuracy  float64 `json:"accuracy"`
}

type LocationEventResponse struct {
	core.EventResponse

	Extras LocationResponse `json:"extras"`
}
