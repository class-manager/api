package model

import "github.com/gofrs/uuid"

// Class represents a Classman class
type Class struct {
	BaseModel
	Account     Account   `gorm:"constraint:OnDelete:CASCADE"` // The account that created this class
	AccountID   uuid.UUID `gorm:"index:;not null"`
	Name        string    `gorm:"not null"`
	SubjectName string    `gorm:"not null"`
}
