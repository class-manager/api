package model

import (
	"time"

	"github.com/gofrs/uuid"
)

// RefreshToken represents a Classman refresh token.
type RefreshToken struct {
	BaseModel
	Token      string    `go:"index:,unique; not null"`
	Expiration time.Time `go:"index:; not null"`
	Account    Account   `go:"constraint:OnDelete:CASCADE"`
	AccountID  uuid.UUID `go:"index;: not null"`
}
