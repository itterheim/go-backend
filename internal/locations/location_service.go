package locations

import (
	"backend/internal/core"
	"errors"
	"fmt"
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

func (s *LocationService) ListHistory(from, to time.Time) ([]GpsHistory, error) {
	return s.locationRepo.ListHistory()
}

func (s *LocationService) GetHistory(id int64) (*GpsHistory, error) {
	return s.locationRepo.GetHistory(id)
}

func (s *LocationService) RegisterHistory(request *CreateGpsHistoryRequest) (*GpsHistoryResponse, error) {
	if request.Timestamp.IsZero() {
		request.Timestamp = time.Now()
	}

	// create event
	event, err := s.eventRepo.CreateEvent(&core.Event{
		Type:       core.EventTypeMoment,
		Timestamp:  request.Timestamp,
		Status:     core.EventStatusPending,
		Tags:       request.Tags,
		Note:       request.Note,
		Reference:  LocationGPSHistoryTable,
		ProviderID: request.ProviderID,
		UserID:     request.UserID,
	})
	if err != nil {
		return nil, errors.New("LocationService.RegisterHistory: failed to create event\n" + err.Error())
	}

	// create gps history
	history, err := s.locationRepo.CreateHistory(&GpsHistory{
		Latitude:  request.Latitude,
		Longitude: request.Longitude,
		Accuracy:  request.Accuracy,
	})
	if err != nil {
		return nil, errors.New("LocationService.RegisterHistory: failed to create gps history\n" + err.Error())
	}

	return &GpsHistoryResponse{
		Event:     event,
		Latitude:  history.Latitude,
		Longitude: history.Longitude,
		Accuracy:  history.Accuracy,
		Created:   history.Created,
	}, nil
}

func (s *LocationService) UpdateHistory(userId int64, data *GpsHistory) (*GpsHistory, error) {
	history, err := s.GetHistory(data.ID)
	if err != nil {
		return nil, err
	}

	// TODO: check if user actually owns this history
	fmt.Println(history)

	return s.locationRepo.UpdateHistory(data)
}

func (s *LocationService) DeleteHistory(id int64) error {
	return s.locationRepo.DeleteHistory(id)
}
