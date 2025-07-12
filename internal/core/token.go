package core

import (
	"time"
)

type Token struct {
	ID         int64
	Created    time.Time
	Updated    time.Time
	UserID     int64
	JTI        string
	Expiration time.Time
	Blocked    bool
}
