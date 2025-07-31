package core

type Tag struct {
	ID          int64  `json:"id"`
	Tag         string `json:"tag"`
	Description string `json:"description,omitempty"`
	ParentID    *int64 `json:"parentId,omitempty"`
	Private     bool   `json:"private"`
}

type CreateTagRequest struct {
	Tag         string `json:"tag"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
}

type UpdateTagRequest struct {
	ID          int64  `json:"-"`
	Tag         string `json:"tag"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
}
