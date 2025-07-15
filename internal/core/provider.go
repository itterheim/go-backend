package core

import "time"

type Provider struct {
	ID          int64      `json:"id"`
	Created     time.Time  `json:"created"`
	Updated     time.Time  `json:"updated"`
	UserID      int64      `json:"userId"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	JTI         *string    `json:"-"`
	Expiration  *time.Time `json:"expiration,omitempty"`
}
