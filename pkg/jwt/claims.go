package jwt

import "time"

type ClaimType string
type ClaimRole string

const (
	UserClaim     ClaimType = "user"
	ProviderClaim ClaimType = "provider"

	OwnerRole ClaimRole = "owner"
	GuestRole ClaimRole = "guest"
)

type Claims struct {
	UserID     int64
	ProviderID *int64
	Role       ClaimRole
	Type       ClaimType
	Expiration time.Time
	JTI        string
}
