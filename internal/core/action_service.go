package core

import (
	"errors"
)

type ActionService struct {
	eventRepository  *EventRepository
	actionRepository *ActionRepository
}

func NewActionService(actionRepository *ActionRepository, eventRepository *EventRepository) *ActionService {
	return &ActionService{
		actionRepository: actionRepository,
		eventRepository:  eventRepository,
	}
}

func (s *ActionService) ListActions(userId int64) ([]ActionResponse, error) {
	return s.actionRepository.ListActions()
}

func (s *ActionService) GetAction(id, userId int64) (*ActionResponse, error) {
	return s.actionRepository.GetAction(id)
}

func (s *ActionService) DeleteAction(id, userId int64) error {
	return s.actionRepository.DeleteAction(id)
}

func (s *ActionService) CreateAction(userId int64, action *CreateActionRequest) (*ActionResponse, error) {
	event, err := s.eventRepository.GetEvent(action.EventID)
	if err != nil {
		return nil, errors.New("ActionService.CreateAction: failed to fetch event")
	}
	if event == nil {
		return nil, errors.New("ActionService.CreateAction: event not found")
	}
	if event.UserID != userId {
		return nil, errors.New("ActionService.CreateAction: user not authorized")
	}

	result, err := s.actionRepository.CreateAction(action)
	if err != nil {
		return nil, errors.New("ActionService.CreateAction: failed to create action")
	}

	return result, nil
}

func (s *ActionService) UpdateAction(userId int64, action *UpdateActionRequest) (*ActionResponse, error) {
	originalAction, err := s.actionRepository.GetAction(action.ID)
	if err != nil {
		return nil, errors.New("ActionService.UpdateAction: action not forund")
	}

	event, err := s.eventRepository.GetEvent(originalAction.EventID)
	if err != nil {
		return nil, errors.New("ActionService.CreateAction: failed to fetch event for the stored action")
	}
	if event == nil {
		return nil, errors.New("ActionService.CreateAction: event not found for the stored action")
	}
	if event.UserID != userId {
		return nil, errors.New("ActionService.CreateAction: user not authorized to update the original action")
	}

	event, err = s.eventRepository.GetEvent(action.EventID)
	if err != nil {
		return nil, errors.New("ActionService.CreateAction: failed to fetch event")
	}
	if event == nil {
		return nil, errors.New("ActionService.CreateAction: event not found")
	}
	if event.UserID != userId {
		return nil, errors.New("ActionService.CreateAction: user not authorized")
	}

	result, err := s.actionRepository.UpdateAction(action)
	if err != nil {
		return nil, errors.New("ActionService.CreateAction: failed to create action")
	}

	return result, nil
}
