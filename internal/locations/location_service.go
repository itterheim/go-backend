package locations

type LocationService struct {
	repo *LocationRepository
}

func NewLocationService(repo *LocationRepository) *LocationService {
	return &LocationService{repo: repo}
}
