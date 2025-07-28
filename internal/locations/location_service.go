package locations

import (
	"backend/internal/core"
	"errors"
	"time"
)

type LocationService struct {
	locationRepo *LocationRepository
	eventRepo    *core.EventRepository
}

func NewLocationService(locationRepo *LocationRepository, eventRepo *core.EventRepository) *LocationService {
	return &LocationService{
		locationRepo: locationRepo,
		eventRepo:    eventRepo,
	}
}

func (s *LocationService) ListHistory(from, to time.Time, userID int64) ([]GpsHistory, error) {
	return s.locationRepo.ListHistory(userID)
}

func (s *LocationService) GetHistory(id, userID int64) (*GpsHistory, error) {
	return s.locationRepo.GetHistory(id, userID)
}

func (s *LocationService) RegisterHistory(request *CreateGpsHistoryRequest) (*GpsHistoryResponse, error) {
	if request.Timestamp.IsZero() {
		request.Timestamp = time.Now()
	}

	// create gps history
	history, err := s.locationRepo.CreateHistory(&GpsHistory{
		Timestamp:  request.Timestamp,
		Latitude:   request.Latitude,
		Longitude:  request.Longitude,
		Accuracy:   request.Accuracy,
		ProviderID: request.ProviderID,
		UserID:     request.UserID,
	})
	if err != nil {
		return nil, errors.New("LocationService.RegisterHistory: failed to create gps history\n" + err.Error())
	}

	return &GpsHistoryResponse{
		ID:        history.ID,
		Timestamp: history.Timestamp,
		Latitude:  history.Latitude,
		Longitude: history.Longitude,
		Accuracy:  history.Accuracy,
	}, nil
}

func (s *LocationService) UpdateHistory(data *GpsHistory) (*GpsHistory, error) {
	if data.Timestamp.IsZero() {
		data.Timestamp = time.Now()
	}
	return s.locationRepo.UpdateHistory(data)
}

func (s *LocationService) DeleteHistory(id, updateID int64) error {
	return s.locationRepo.DeleteHistory(id, updateID)
}
