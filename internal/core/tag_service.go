package core

type TagService struct {
	repo *TagRepository
}

func NewTagService(repo *TagRepository) *TagService {
	return &TagService{repo}
}

func (s *TagService) ListTags(private bool) ([]Tag, error) {
	return s.repo.ListTags(private)
}

func (s *TagService) GetTag(tag string) (*Tag, error) {
	return s.repo.GetTag(tag)
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
		Tag:         data.Tag,
		Description: data.Description,
	}, data.NewTag)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

func (s *TagService) DeleteTag(tag string) error {
	return s.repo.DeleteTag(tag)
}
