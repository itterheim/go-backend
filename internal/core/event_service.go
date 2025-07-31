package core

import (
	"fmt"
	"time"
)

type EventService struct {
	repo *EventRepository
}

func NewEventService(repo *EventRepository) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) ListEvents(from, to time.Time) ([]EventResponse, error) {
	events, err := s.repo.ListEvents()
	if err != nil {
		return nil, err
	}

	result := make([]EventResponse, len(events))
	for i, event := range events {
		result[i] = EventResponse{
			ID:        event.ID,
			Type:      event.Type,
			Timestamp: event.Timestamp,
			Until:     event.Until,
			Tags:      event.Tags,
			Note:      event.Note,
			Reference: event.Reference,
		}
	}

	return result, nil
}

func (s *EventService) GetEvent(id int64) (*EventResponse, error) {
	event, err := s.repo.GetEvent(id)
	if err != nil {
		return nil, err
	}

	result := &EventResponse{
		ID:        event.ID,
		Type:      event.Type,
		Timestamp: event.Timestamp,
		Until:     event.Until,
		Tags:      event.Tags,
		Note:      event.Note,
		Reference: event.Reference,
	}

	return result, nil
}

func (s *EventService) CreateEvent(request *CreateEventRequest) (*EventResponse, error) {
	// err := s.validateType(event.Type)
	err := request.Validate()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if request.Type == EventTypeMoment && request.Timestamp == nil {
		now := time.Now()
		request.Timestamp = &now
	}

	event, err := s.repo.CreateEvent(&Event{
		Type:       request.Type,
		Timestamp:  request.Timestamp,
		Until:      request.Until,
		Tags:       request.Tags,
		Note:       request.Note,
		ProviderID: request.ProviderID,
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	result := &EventResponse{
		ID:        event.ID,
		Type:      event.Type,
		Timestamp: event.Timestamp,
		Until:     event.Until,
		Tags:      event.Tags,
		Note:      event.Note,
		Reference: event.Reference,
	}

	return result, nil
}

func (s *EventService) UpdateEvent(request *UpdateEventRequest) (*EventResponse, error) {
	// err := s.validateType(event.Type)
	err := request.Validate()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if request.Timestamp.IsZero() {
		now := time.Now()
		request.Timestamp = &now
	}

	event, err := s.repo.UpdateEvent(&Event{
		ID:         request.ID,
		Type:       request.Type,
		Timestamp:  request.Timestamp,
		Until:      request.Until,
		Tags:       request.Tags,
		Note:       request.Note,
		ProviderID: request.ProviderID,
	})
	if err != nil {
		return nil, err
	}

	return &EventResponse{
		ID:        event.ID,
		Type:      event.Type,
		Timestamp: event.Timestamp,
		Until:     event.Until,
		Tags:      event.Tags,
		Note:      event.Note,
	}, nil
}

func (s *EventService) DeleteEvent(id int64) error {
	return s.repo.DeleteEvent(id)
}

// func (s *EventService[T]) validateType(eventType EventType) error {
// 	if eventType == EventTypeInterval {
// 		return nil
// 	}

// 	if eventType == EventTypeMoment {
// 		return nil
// 	}

// 	return errors.New("EventService.validateType: invalid type " + string(eventType))
// }
