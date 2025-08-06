package core

import (
	"fmt"
)

type EventService struct {
	repo *EventRepository
}

func NewEventService(repo *EventRepository) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) ListEvents(query *EventQueryBuilder) ([]EventResponse, error) {
	events, err := s.repo.ListEvents(query)
	if err != nil {
		return nil, err
	}

	result := make([]EventResponse, len(events))
	for i, event := range events {
		result[i] = *event.ToEventResponse()
	}

	return result, nil
}

func (s *EventService) GetEvent(id int64) (*EventResponse, error) {
	event, err := s.repo.GetEvent(id)
	if err != nil {
		return nil, err
	}

	return event.ToEventResponse(), nil
}

func (s *EventService) CreateEvent(request *CreateEventRequest) (*EventResponse, error) {
	err := request.Validate()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	event, err := s.repo.CreateEvent(request.ToEvent())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return event.ToEventResponse(), nil
}

func (s *EventService) UpdateEvent(request *UpdateEventRequest) (*EventResponse, error) {
	err := request.Validate()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	event, err := s.repo.UpdateEvent(request.ToEvent())
	if err != nil {
		return nil, err
	}

	return event.ToEventResponse(), nil
}

func (s *EventService) DeleteEvent(id int64) error {
	return s.repo.DeleteEvent(id)
}
