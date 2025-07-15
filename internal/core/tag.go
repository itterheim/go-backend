package core

type Tag struct {
	ID          int64  `json:"id"`
	Tag         string `json:"tag"`
	Description string `json:"description"`
	ParentID    *int64 `json:"parentId"`
	UserID      int64  `json:"-"`
}

type CreateTagRequest struct {
	Tag         string `json:"tag"`
	Description string `json:"description"`
	UserID      int64  `json:"-"`
}

type UpdateTagRequest struct {
	ID          int64  `json:"-"`
	Tag         string `json:"tag"`
	Description string `json:"description"`
	UserID      int64  `json:"-"`
}
