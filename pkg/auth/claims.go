package auth

import "time"

type ClaimType string

const (
	UserClaim   ClaimType = "user"
	DeviceClaim ClaimType = "device"
)

type Claims struct {
	ID         int64
	Type       ClaimType
	Expiration time.Time
	JTI        string
}
