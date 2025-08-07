package core

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
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
		ID:         e.ID,
		Type:       e.Type,
		Timestamp:  e.Timestamp,
		Until:      e.Until,
		Tags:       e.Tags,
		Note:       e.Note,
		Reference:  e.Reference,
		ProviderID: e.ProviderID,
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
	ID         int64      `json:"id"`
	Type       EventType  `json:"type"`
	Timestamp  *time.Time `json:"timestamp,omitempty"`
	Until      *time.Time `json:"until,omitempty"`
	Tags       []string   `json:"tags,omitempty"`
	Note       string     `json:"note,omitempty"`
	Reference  string     `json:"reference"`
	ProviderID *int64     `json:"providerId,omitempty"`
}

type EventQueryBuilder struct {
	Type    EventType
	From    time.Time
	To      time.Time
	Private bool
	Tags    []string
}

func (b *EventQueryBuilder) FromRequest(r *http.Request) error {
	b.Private = r.URL.Query().Has("private")

	if r.URL.Query().Has("type") {
		b.Type = EventType(r.URL.Query().Get("type"))
	}

	b.Tags = []string{}
	if r.URL.Query().Has("tags") {
		query := r.URL.Query().Get("tags")
		b.Tags = strings.Split(query, ",")
	}

	var err error = nil
	if r.URL.Query().Has("from") {
		b.From, err = time.Parse(time.RFC3339, r.URL.Query().Get("from"))
		if err != nil {
			return err
		}
	}

	if r.URL.Query().Has("to") {
		b.To, err = time.Parse(time.RFC3339, r.URL.Query().Get("to"))
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *EventQueryBuilder) Build() (string, []any) {
	params := make([]any, 0)
	and := make([]string, 0)

	if len(b.Type) > 0 {
		params = append(params, b.Type)
		where := fmt.Sprintf("events.type = $%v", len(params))
		and = append(and, "("+where+")")
	}

	if !b.From.IsZero() && !b.To.IsZero() {
		params = append(params, b.From, b.To)
		where := fmt.Sprintf("(events.type = 'interval' AND events.timestamp <= $%[2]v AND events.until >= $%[1]v) OR (events.type = 'moment' AND events.timestamp >= $%[1]v AND events.timestamp < $%[2]v)", len(params)-1, len(params))
		and = append(and, "("+where+")")
	}

	if !b.Private {
		where := `NOT EXISTS (
			SELECT 1 FROM UNNEST(events.tags) AS event_tag JOIN tags ON event_tag = tags.tag WHERE tags.private = TRUE
		)`
		and = append(and, "("+where+")")
	}

	if len(b.Tags) > 0 {
		params = append(params, b.Tags)
		where := fmt.Sprintf("events.tags @> $%v", len(params))
		and = append(and, "("+where+")")
	}

	if len(and) == 0 {
		return "", []any{}
	} else {
		where := "WHERE " + strings.Join(and, " AND ")

		return where, params
	}

}
