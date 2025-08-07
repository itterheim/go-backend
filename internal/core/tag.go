package core

type Tag struct {
	Tag         string  `json:"tag"`
	Description *string `json:"description,omitempty"`
	Parent      *string `json:"parent,omitempty"`
	Private     bool    `json:"private"`
}

type CreateTagRequest struct {
	Tag         string  `json:"tag"`
	Description *string `json:"description"`
	Private     bool    `json:"private"`
}

type UpdateTagRequest struct {
	Tag         string  `json:"tag"`
	NewTag      *string `json:"newTag,omitempty"`
	Description *string `json:"description"`
	Private     bool    `json:"private"`
}
