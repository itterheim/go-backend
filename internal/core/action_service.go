package core

import (
	"errors"
)

type ActionService struct {
	repo *ActionRepository
}

func NewActionService(repo *ActionRepository) *ActionService {
	return &ActionService{repo: repo}
}

func (s *ActionService) ListActions(userId int64) ([]Action, error) {
	return s.repo.ListActions()
}

func (s *ActionService) GetAction(id, userId int64) (*Action, error) {
	return s.repo.GetAction(id)
}

func (s *ActionService) DeleteAction(id, userId int64) error {
	return s.repo.DeleteAction(id)
}

func (s *ActionService) CreateAction(data *CreateActionRequest) (*Action, error) {
	return nil, errors.New("Not implemented")
}

func (s *ActionService) UpdateAction(data *UpdateActionRequest) (*Action, error) {
	return nil, errors.New("Not implemented")
}
