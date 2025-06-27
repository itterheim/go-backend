package models

import "time"

type Device struct {
	ID             int64
	Created        time.Time
	Updated        time.Time
	Name           string
	Description    string
	Identification string
	JTI            string
	Expiration     time.Time
}

type EnviDevice struct {
	DeviceId int64
	Sensors  []string
}
