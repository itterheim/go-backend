package core

import (
	"errors"
	"time"
)

type EventService struct {
	repo *EventRepository
}

func NewEventService(repo *EventRepository) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) ListEvents(from, to time.Time) ([]Event, error) {
	return s.repo.ListEvents()
}

func (s *EventService) GetEvent(id int64) (*Event, error) {
	return s.repo.GetEvent(id)
}

func (s *EventService) CreateEvent(event *CreateEventRequest) (*Event, error) {
	return nil, errors.New("Not implemented")
}

func (s *EventService) UpdateEvent(event *UpdateEventRequest) (*Event, error) {
	return nil, errors.New("Not implemented")
}

func (s *EventService) DeleteEvent(id int64) error {
	return s.repo.DeleteEvent(id)
}
