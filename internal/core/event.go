package core

import (
	"errors"
	"time"
)

type EventType string

const (
	EventTypeMoment   EventType = "moment"
	EventTypeInterval EventType = "interval"
)

type EventValidation struct{}

type Event struct {
	ID         int64      `json:"id"`
	Type       EventType  `json:"type"`
	Timestamp  *time.Time `json:"timestamp,omitempty"`
	Until      *time.Time `json:"until,omitempty"`
	Tags       []string   `json:"tags,omitempty"`
	Note       string     `json:"note,omitempty"`
	Reference  string     `json:"reference,omitempty"`
	ProviderID *int64     `json:"providerId,omitempty"`
}

type EventRequest struct {
	Type       EventType  `json:"type"`
	Timestamp  *time.Time `json:"timestamp,omitempty"`
	Until      *time.Time `json:"until,omitempty"`
	Tags       []string   `json:"tags,omitempty"`
	Note       string     `json:"note,omitempty"`
	Reference  string     `json:"reference,omitempty"`
	ProviderID *int64     `json:"-"`
}

func (e *EventRequest) Validate() error {
	if e.Type == EventTypeInterval {
		if e.Timestamp == nil && e.Until == nil {
			return errors.New("EventRequest.Validate: missing timestamp or until")
		}
		return nil
	}

	if e.Type == EventTypeMoment {
		if e.Timestamp == nil {
			now := time.Now()
			e.Timestamp = &now
		}
		return nil
	}

	return errors.New("EventRequest.Validate: invalid type " + string(e.Type))
}

type CreateEventRequest struct {
	EventRequest
}

type UpdateEventRequest struct {
	EventRequest
	ID int64 `json:"id"`
}

type EventResponse struct {
	ID        int64      `json:"id"`
	Type      EventType  `json:"type"`
	Timestamp *time.Time `json:"timestamp,omitempty"`
	Until     *time.Time `json:"until,omitempty"`
	Tags      []string   `json:"tags,omitempty"`
	Note      string     `json:"note,omitempty"`
	Reference string     `json:"reference,omitempty"`
}
