package locations

import (
	"errors"
)

type PlaceService struct {
	placeRepo *PlaceRepository
}

func NewPlaceService(placeRepo *PlaceRepository) *PlaceService {
	return &PlaceService{placeRepo}
}

func (s *PlaceService) ListPlaces() ([]Place, error) {
	return s.placeRepo.ListPlaces()
}

func (s *PlaceService) GetPlace(id int64) (*Place, error) {
	return s.placeRepo.GetPlace(id)
}

func (s *PlaceService) CreatePlace(request *CreatePlaceRequest) (*PlaceResponse, error) {
	// create
	place, err := s.placeRepo.CreatePlace(&Place{
		Name:      request.Name,
		Note:      request.Note,
		Latitude:  request.Latitude,
		Longitude: request.Longitude,
		Radius:    request.Radius,
	})
	if err != nil {
		return nil, errors.New("PlaceService.CreatePlace: failed to create place\n" + err.Error())
	}

	return &PlaceResponse{
		ID:        place.ID,
		Name:      place.Name,
		Note:      place.Note,
		Latitude:  place.Latitude,
		Longitude: place.Longitude,
		Radius:    place.Radius,
		Created:   place.Created,
		Updated:   place.Updated,
	}, nil
}

func (s *PlaceService) UpdateHistory(data *Place) (*Place, error) {
	return s.placeRepo.UpdatePlace(data)
}

func (s *PlaceService) DeleteHistory(id int64) error {
	return s.placeRepo.DeletePlace(id)
}
