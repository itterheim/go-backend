package core

import "time"

type EventStatus string
type EventType string

const (
	EventStatusPending   EventStatus = "Pending"
	EventStatusApproved  EventStatus = "Approved"
	EventStatusRejected  EventStatus = "Rejected"
	EventStatusDraft     EventStatus = "Draft"
	EventStatusScheduled EventStatus = "Scheduled"
	// TODO: add the rest of statuses

	EventTypeMoment   EventType = "Moment"
	EventTypeInterval EventType = "Interval"
)

var EventStatuses = map[EventStatus]bool{
	EventStatusPending:   true,
	EventStatusApproved:  true,
	EventStatusRejected:  true,
	EventStatusDraft:     true,
	EventStatusScheduled: true,
}

type Event struct {
	ID         int64       `json:"id"`
	Type       EventType   `json:"type"`
	Timestamp  time.Time   `json:"timestamp"`
	Until      *time.Time  `json:"until,omitempty"`
	Status     EventStatus `json:"status"`
	Tags       []string    `json:"tags"`
	Note       string      `json:"note,omitempty"`
	Reference  string      `json:"reference,omitempty"`
	ProviderID *int64      `json:"-"`
	UserID     int64       `json:"-"`
}

type CreateEventRequest struct {
	Type       EventType   `json:"type"`
	Timestamp  time.Time   `json:"timestamp"`
	Until      *time.Time  `json:"until"`
	Status     EventStatus `json:"status"`
	Tags       []string    `json:"tags"`
	Note       string      `json:"note"`
	UserID     int64       `json:"-"`
	ProviderID *int64      `json:"-"`
}

type UpdateEventRequest struct {
	ID         int64       `json:"id"`
	Type       EventType   `json:"type"`
	Timestamp  time.Time   `json:"timestamp"`
	Until      *time.Time  `json:"until"`
	Status     EventStatus `json:"status"`
	Tags       []string    `json:"tags"`
	Note       string      `json:"note"`
	UserID     int64       `json:"-"`
	ProviderID *int64      `json:"-"`
}
