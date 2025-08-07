package core

type TagService struct {
	tagRepo   *TagRepository
	eventRepo *EventRepository
}

func NewTagService(tagRepo *TagRepository, eventRepo *EventRepository) *TagService {
	return &TagService{tagRepo, eventRepo}
}

func (s *TagService) ListTags(private bool) ([]Tag, error) {
	return s.tagRepo.ListTags(private)
}

func (s *TagService) GetTag(tag string) (*Tag, error) {
	return s.tagRepo.GetTag(tag)
}

func (s *TagService) CreateTag(data *CreateTagRequest) (*Tag, error) {
	// TODO: find parent based on tag string
	tag, err := s.tagRepo.CreateTag(&Tag{
		Tag:         data.Tag,
		Description: data.Description,
		Private:     data.Private,
	})
	if err != nil {
		return nil, err
	}

	return tag, nil
}

func (s *TagService) UpdateTag(data *UpdateTagRequest) (*Tag, error) {
	// TODO: find and update parent based on tag string
	tag, err := s.tagRepo.UpdateTag(&Tag{
		Tag:         data.Tag,
		Description: data.Description,
		Private:     data.Private,
	}, data.NewTag)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

func (s *TagService) DeleteTag(tag string) error {
	return s.tagRepo.DeleteTag(tag)
}

func (s *TagService) SynchronizeTags() error {
	tags, err := s.eventRepo.UsedTags()
	if err != nil {
		return err
	}

	err = s.tagRepo.SynchronizeTags(tags)

	return err
}
