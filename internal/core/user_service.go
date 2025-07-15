package core

type UserService struct {
	repo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
	return &UserService{repo}
}

func (s *UserService) ListUsers() ([]User, error) {
	return s.repo.ListUsers()
}

func (s *UserService) GetUser(id int64) (*User, error) {
	return s.repo.GetUser(id)
}
