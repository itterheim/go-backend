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
	ID         int64
	Type       EventType
	Timestamp  *time.Time
	Until      *time.Time
	Tags       []string
	Note       string
	Reference  string
	ProviderID *int64
}

func (e *Event) ToEventResponse() *EventResponse {
	return &EventResponse{
		ID:        e.ID,
		Type:      e.Type,
		Timestamp: e.Timestamp,
		Until:     e.Until,
		Tags:      e.Tags,
		Note:      e.Note,
	}
}

type EventRequest struct {
	Type       EventType  `json:"type"`
	Timestamp  *time.Time `json:"timestamp,omitempty"`
	Until      *time.Time `json:"until,omitempty"`
	Tags       []string   `json:"tags,omitempty"`
	Note       string     `json:"note,omitempty"`
	Reference  string     `json:"-"`
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

func (e *EventRequest) ToEvent() *Event {
	return &Event{
		Type:       e.Type,
		Timestamp:  e.Timestamp,
		Until:      e.Until,
		Tags:       e.Tags,
		Note:       e.Note,
		Reference:  e.Reference,
		ProviderID: e.ProviderID,
	}
}

type CreateEventRequest struct {
	EventRequest
}

type UpdateEventRequest struct {
	EventRequest
	ID int64 `json:"id"`
}

func (e *UpdateEventRequest) ToEvent() *Event {
	event := e.EventRequest.ToEvent()
	event.ID = e.ID

	return event
}

type EventResponse struct {
	ID        int64      `json:"id"`
	Type      EventType  `json:"type"`
	Timestamp *time.Time `json:"timestamp,omitempty"`
	Until     *time.Time `json:"until,omitempty"`
	Tags      []string   `json:"tags,omitempty"`
	Note      string     `json:"note,omitempty"`
}
