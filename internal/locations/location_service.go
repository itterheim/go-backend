package locations

import "backend/internal/core"

type LocationService struct {
	locationRepo *LocationRepository
	eventRepo    *core.EventRepository
	actionRepo   *core.ActionRepository
}

func NewLocationService(locationRepo *LocationRepository, eventRepo *core.EventRepository, actionRepo *core.ActionRepository) *LocationService {
	return &LocationService{
		locationRepo: locationRepo,
		eventRepo:    eventRepo,
		actionRepo:   actionRepo,
	}
}
