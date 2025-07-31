package jwt

import "time"

type ClaimType string

const (
	UserClaim     ClaimType = "user"
	ProviderClaim ClaimType = "provider"
)

type Claims struct {
	UserID     int64
	ProviderID *int64
	Type       ClaimType
	Expiration time.Time
	JTI        string
}
