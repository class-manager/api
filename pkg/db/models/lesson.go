package model

import (
	"time"

	"github.com/gofrs/uuid"
)

// Lesson represents a Classman lesson
type Lesson struct {
	BaseModel
	Name        string `gorm:"not null"`
	Description string
	StartTime   time.Time `gorm:"not null"`
	EndTime     time.Time `gorm:"not null"`
	Class       Class     `gorm:"constraint:OnDelete:CASCADE"`
	ClassID     uuid.UUID
}
