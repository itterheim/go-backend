package locations

import (
	"backend/internal/core"
)

type Location struct {
	ID        int64   `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Accuracy  float64 `json:"accuracy"`
	EventID   int64   `json:"eventId"`
}

type LocationEvent struct {
	core.Event
	Extras Location
}

type CreateLocationRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Accuracy  float64 `json:"accuracy"`
}

type CreateLocationEventRequest struct {
	core.CreateEventRequest

	Extras CreateLocationRequest `json:"extras"`
}

type UpdateLocationRequest struct {
	ID        int64   `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Accuracy  float64 `json:"accuracy"`
}

type UpdateLocationEventRequest struct {
	core.UpdateEventRequest

	Extras UpdateLocationRequest `json:"extras"`
}

type LocationResponse struct {
	ID        int64   `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Accuracy  float64 `json:"accuracy"`
}

type LocationEventResponse struct {
	core.EventResponse

	Extras LocationResponse `json:"extras"`
}
