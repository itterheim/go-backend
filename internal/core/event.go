package core

import "time"

type EventStatus string

const (
	EventStatusPending EventStatus = "Pending"
	EventStatusLocked  EventStatus = "Locked"
)

type Event struct {
	ID         int64     `json:"id"`
	Type       string    `json:"type"`
	Timestamp  time.Time `json:"timestamp"`
	Until      time.Time `json:"until"`
	Status     string    `json:"status"`
	Tags       []string  `json:"tags"`
	Note       *string   `json:"note"`
	ProviderID *int64    `json:"-"`
	UserID     int64     `json:"-"`
}

type CreateEventRequest struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Until     time.Time `json:"until"`
	Status    string    `json:"status"`
	Tags      []string  `json:"tags"`
	Note      *string   `json:"note"`
	UserID    int64     `json:"-"`
}

type UpdateEventRequest struct {
	ID         int64     `json:"id"`
	Type       string    `json:"type"`
	Timestamp  time.Time `json:"timestamp"`
	Until      time.Time `json:"until"`
	Status     string    `json:"status"`
	Tags       []string  `json:"tags"`
	Note       *string   `json:"note"`
	ProviderID *int64    `json:"-"`
	UserId     int64     `json:"-"`
}
