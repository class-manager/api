package model

import (
	"time"

	"github.com/gofrs/uuid"
)

// Task represents a Classman task
type Task struct {
	BaseModel
	Name        string `gorm:"not null"`
	Description *string
	Open        bool
	OpenDate    time.Time
	DueDate     time.Time `gorm:"not null"`
	MaxMark     int32     `gorm:"not null"`
	Class       Class     `gorm:"constraint:OnDelete:CASCADE"`
	ClassID     uuid.UUID `gorm:"index:; not null"`
}
