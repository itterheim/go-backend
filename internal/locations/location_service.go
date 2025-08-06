package locations

import (
	"backend/internal/core"
	"errors"
	"fmt"
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

func (s *LocationService) ListHistory(query *core.EventQueryBuilder) ([]LocationEventResponse, error) {
	data, err := s.locationRepo.ListHistory(query)
	if err != nil {
		return nil, err
	}

	result := make([]LocationEventResponse, len(data))
	for i, event := range data {
		result[i] = LocationEventResponse{
			EventResponse: *event.ToEventResponse(),
			Extras:        *event.Extras.ToLocationResponse(),
		}
	}

	return result, nil
}

func (s *LocationService) GetHistory(id int64) (*LocationEventResponse, error) {
	data, err := s.locationRepo.GetHistory(id)
	if err != nil {
		return nil, fmt.Errorf("LocationService.GetHistory: failed to retrieve LocationEvent, %v", err)
	}

	return &LocationEventResponse{
		EventResponse: *data.ToEventResponse(),
		Extras:        *data.Extras.ToLocationResponse(),
	}, nil
}

func (s *LocationService) RegisterHistory(request *CreateLocationEventRequest) (*LocationEventResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, fmt.Errorf("LocationService.RegisterHistory: validation failed, %v", err)
	}

	request.Reference = LocationGPSHistoryTable
	request.Tags = append(request.Tags, "module:locations")

	event, err := s.eventRepo.CreateEvent(request.CreateEventRequest.ToEvent())
	if err != nil {
		return nil, errors.New("LocationService.RegisterHistory: failed to create event\n" + err.Error())
	}

	// create gps history
	history, err := s.locationRepo.CreateHistory(&Location{
		Latitude:  request.Extras.Latitude,
		Longitude: request.Extras.Longitude,
		Accuracy:  request.Extras.Accuracy,
		EventID:   event.ID,
	})
	if err != nil {
		return nil, errors.New("LocationService.RegisterHistory: failed to create gps history\n" + err.Error())
	}

	return &LocationEventResponse{
		EventResponse: *event.ToEventResponse(),
		Extras:        *history.ToLocationResponse(),
	}, nil
}

func (s *LocationService) UpdateHistory(request *UpdateLocationEventRequest) (*LocationEventResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, fmt.Errorf("LocationService.UpdateHistory: validation failed, %v", err)
	}

	request.Reference = LocationGPSHistoryTable

	// update event
	event, err := s.eventRepo.UpdateEvent(request.UpdateEventRequest.ToEvent())
	if err != nil {
		return nil, errors.New("LocationService.UpdateHistory: failed to update event\n" + err.Error())
	}

	// update gps history
	history, err := s.locationRepo.UpdateHistory(&Location{
		Latitude:  request.Extras.Latitude,
		Longitude: request.Extras.Longitude,
		Accuracy:  request.Extras.Accuracy,
		EventID:   event.ID,
	})
	if err != nil {
		return nil, errors.New("LocationService.UpdateHistory: failed to update location\n" + err.Error())
	}

	return &LocationEventResponse{
		EventResponse: *event.ToEventResponse(),
		Extras:        *history.ToLocationResponse(),
	}, nil
}

func (s *LocationService) DeleteHistory(id int64) error {
	return s.locationRepo.DeleteHistory(id)
}
