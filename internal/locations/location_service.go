package locations

import (
	"backend/internal/core"
	"errors"
	"fmt"
	"time"
)

type LocationService struct {
	locationRepo *LocationRepository
	eventService *core.EventService
}

func NewLocationService(locationRepo *LocationRepository, eventService *core.EventService) *LocationService {
	return &LocationService{
		locationRepo: locationRepo,
		eventService: eventService,
	}
}

func (s *LocationService) ListHistory(from, to time.Time) ([]LocationEventResponse, error) {
	data, err := s.locationRepo.ListHistory()
	if err != nil {
		return nil, err
	}

	result := make([]LocationEventResponse, len(data))
	for i, event := range data {
		result[i] = LocationEventResponse{
			EventResponse: core.EventResponse{
				ID:        event.ID,
				Type:      event.Type,
				Timestamp: event.Timestamp,
				Until:     event.Until,
				Tags:      event.Tags,
				Note:      event.Note,
				Reference: event.Reference,
			},
			Extras: LocationResponse{
				ID:        event.Extras.ID,
				Latitude:  event.Extras.Latitude,
				Longitude: event.Extras.Longitude,
				Accuracy:  event.Extras.Accuracy,
			},
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
		EventResponse: core.EventResponse{
			ID:        data.ID,
			Type:      data.Type,
			Timestamp: data.Timestamp,
			Until:     data.Until,
			Tags:      data.Tags,
			Note:      data.Note,
			Reference: data.Reference,
		},
		Extras: LocationResponse{
			ID:        data.Extras.ID,
			Latitude:  data.Extras.Latitude,
			Longitude: data.Extras.Longitude,
			Accuracy:  data.Extras.Accuracy,
		},
	}, nil
}

func (s *LocationService) RegisterHistory(request *CreateLocationEventRequest) (*LocationEventResponse, error) {
	request.Reference = LocationGPSHistoryTable
	request.Tags = append(request.Tags, "module:locations")

	event, err := s.eventService.CreateEvent(&request.CreateEventRequest)
	if err != nil {
		return nil, errors.New("LocationService.RegisterHIstory: failed to create event\n" + err.Error())
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
		EventResponse: *event,
		Extras: LocationResponse{
			ID:        history.ID,
			Latitude:  history.Latitude,
			Longitude: history.Longitude,
			Accuracy:  history.Accuracy,
		}}, nil
}

func (s *LocationService) UpdateHistory(request *UpdateLocationEventRequest) (*LocationEventResponse, error) {
	request.Reference = LocationGPSHistoryTable

	// update event
	event, err := s.eventService.UpdateEvent(&request.UpdateEventRequest)
	if err != nil {
		return nil, errors.New("LocationService.UpdateHistory: failed to update event\n" + err.Error())
	}

	// update gps history
	history, err := s.locationRepo.UpdateHistory(&Location{
		ID:        request.Extras.ID,
		Latitude:  request.Extras.Latitude,
		Longitude: request.Extras.Longitude,
		Accuracy:  request.Extras.Accuracy,
		EventID:   event.ID,
	})
	if err != nil {
		return nil, errors.New("LocationService.UpdateHistory: failed to update location\n" + err.Error())
	}

	return &LocationEventResponse{
		EventResponse: *event,
		Extras: LocationResponse{
			ID:        history.ID,
			Latitude:  history.Latitude,
			Longitude: history.Longitude,
			Accuracy:  history.Accuracy,
		},
	}, nil
}

func (s *LocationService) DeleteHistory(id int64) error {
	return s.locationRepo.DeleteHistory(id)
}
