package model

import "github.com/gofrs/uuid"

// TaskResult represents a Classman task result
type TaskResult struct {
	BaseModel
	Mark      float64   `gorm:"not null"`
	Student   Student   `gorm:"constraint:OnDelete:CASCADE"`
	StudentID uuid.UUID `gorm:"not null"`
	Task      Task      `gorm:"constraint:OnDelete:CASCADE"`
	TaskID    uuid.UUID `gorm:"not null"`
}
