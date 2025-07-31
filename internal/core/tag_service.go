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

func (s *TagService) GetTag(id int64) (*Tag, error) {
	return s.repo.GetTag(id)
}

func (s *TagService) CreateTag(data *CreateTagRequest) (*Tag, error) {
	// TODO: find parent based on tag string
	tag, err := s.repo.CreateTag(&Tag{
		Tag:         data.Tag,
		Description: data.Description,
	})
	if err != nil {
		return nil, err
	}

	return tag, nil
}

func (s *TagService) UpdateTag(data *UpdateTagRequest) (*Tag, error) {
	// TODO: find and update parent based on tag string
	tag, err := s.repo.UpdateTag(&Tag{
		ID:          data.ID,
		Tag:         data.Tag,
		Description: data.Description,
	})
	if err != nil {
		return nil, err
	}

	return tag, nil
}

func (s *TagService) DeleteTag(id int64) error {
	return s.repo.DeleteTag(id)
}
