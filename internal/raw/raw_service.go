package raw

import (
	"backend/internal/core"
	"fmt"
)

type RawService struct {
	rawRepo   *RawRepository
	eventRepo *core.EventRepository
}

func NewRawService(rawRepo *RawRepository, eventRepo *core.EventRepository) *RawService {
	return &RawService{rawRepo, eventRepo}
}

func (s *RawService) ListRawEvents(query *core.EventQueryBuilder) ([]RawEventResponse, error) {
	data, err := s.rawRepo.ListRawEvents(query)
	if err != nil {
		return nil, fmt.Errorf("RawService.ListRawEvents: %v", err)
	}

	result := make([]RawEventResponse, len(data))
	for i, raw := range data {
		result[i] = RawEventResponse{
			EventResponse: *raw.ToEventResponse(),
			Extras:        raw.Extras.Data,
		}
	}

	return result, nil
}

func (s *RawService) GetRawEvent(eventID int64) (*RawEventResponse, error) {
	data, err := s.rawRepo.GetRawEvent(eventID)
	if err != nil {
		return nil, fmt.Errorf("RawService.GetRawEvent: %v", err)
	}

	return &RawEventResponse{
		EventResponse: *data.ToEventResponse(),
		Extras:        data.Extras.Data,
	}, nil
}

func (s *RawService) RegisterRawEvent(request *CreateRawEventRequest) (*RawEventResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, fmt.Errorf("RawService.RegisterRawEvent: validation failed, %v", err)
	}

	request.Reference = RawTable
	request.Tags = append(request.Tags, "module:raw")

	event, err := s.eventRepo.CreateEvent(request.CreateEventRequest.ToEvent())
	if err != nil {
		return nil, fmt.Errorf("RawService.RegisterRawEvent: failed to create event, %v", err)
	}

	data, err := s.rawRepo.CreateRaw(&Raw{EventID: event.ID, Data: request.Extras})
	if err != nil {
		return nil, fmt.Errorf("RawService.RegisterEvent: failed to create raw data, %v", err)
	}

	return &RawEventResponse{
		EventResponse: *event.ToEventResponse(),
		Extras:        data.Data,
	}, nil
}

func (s *RawService) UpdateRawEvent(request *UpdateRawEventRequest) (*RawEventResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, fmt.Errorf("RawService.UpdateRawEvent: validation failed, %v", err)
	}

	request.Reference = RawTable

	event, err := s.eventRepo.UpdateEvent(request.UpdateEventRequest.ToEvent())
	if err != nil {
		return nil, fmt.Errorf("RawService.UpdateRawEvent: failed to update event, %v", err)
	}

	data, err := s.rawRepo.UpdateRaw(&Raw{EventID: event.ID, Data: request.Extras})
	if err != nil {
		return nil, fmt.Errorf("RawService.UpdateRawEvent: faile to update raw data, %v", err)
	}

	return &RawEventResponse{
		EventResponse: *event.ToEventResponse(),
		Extras:        data.Data,
	}, nil
}

func (s *RawService) DeleteRawEvent(eventID int64) error {
	return s.rawRepo.DeleteRawEvent(eventID)
}
