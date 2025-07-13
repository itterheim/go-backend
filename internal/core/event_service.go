package core

import (
	"errors"
	"fmt"
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
	err := s.validateType(event.Type)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	if len(event.Status) == 0 {
		event.Status = EventStatusPending
	}

	err = s.validateStatus(event.Status)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	result, err := s.repo.CreateEvent(&Event{
		Type:       event.Type,
		Timestamp:  event.Timestamp,
		Until:      event.Until,
		Status:     event.Status,
		Tags:       event.Tags,
		Note:       event.Note,
		ProviderID: event.ProviderID,
		UserID:     event.UserID,
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return result, nil
}

func (s *EventService) UpdateEvent(event *UpdateEventRequest) (*Event, error) {
	err := s.validateType(event.Type)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	if len(event.Status) == 0 {
		event.Status = EventStatusPending
	}

	err = s.validateStatus(event.Status)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	result, err := s.repo.UpdateEvent(&Event{
		ID:         event.ID,
		Type:       event.Type,
		Timestamp:  event.Timestamp,
		Until:      event.Until,
		Status:     event.Status,
		Tags:       event.Tags,
		Note:       event.Note,
		ProviderID: event.ProviderID,
		UserID:     event.UserID,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *EventService) DeleteEvent(id int64) error {
	return s.repo.DeleteEvent(id)
}

func (s *EventService) validateType(eventType EventType) error {
	if eventType == EventTypeInterval || eventType == EventTypeMoment {
		return nil
	}
	return errors.New("EventService.validateType: invalid type " + string(eventType))
}

func (s *EventService) validateStatus(status EventStatus) error {
	if status == EventStatusLocked || status == EventStatusPending {
		return nil
	}
	return errors.New("EventService.validateType: invalid status " + string(status))
}
