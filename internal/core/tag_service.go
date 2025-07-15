package core

type TagService struct {
	repo *TagRepository
}

func NewTagService(repo *TagRepository) *TagService {
	return &TagService{repo}
}

func (s *TagService) ListTags() ([]Tag, error) {
	return s.repo.ListTags()
}
