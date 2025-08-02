package raw

import (
	"backend/internal/core"
	"encoding/json"
)

type Raw struct {
	EventID int64
	Data    json.RawMessage
}

type RawEvent struct {
	core.Event
	Extras Raw
}

type CreateRawEventRequest struct {
	core.CreateEventRequest
	Extras json.RawMessage `json:"extras"`
}

type UpdateRawEventRequest struct {
	core.UpdateEventRequest
	Extras json.RawMessage `json:"extras"`
}

type RawEventResponse struct {
	core.EventResponse
	Extras json.RawMessage `json:"extras"`
}
