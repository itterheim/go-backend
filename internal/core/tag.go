package core

type Tag struct {
	ID          int64  `json:"id"`
	Tag         string `json:"tag"`
	Description string `json:"description"`
	ParentID    *int64 `json:"parentId"`
	UserID      int64  `json:"-"`
}
