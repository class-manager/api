package model

import "github.com/gofrs/uuid"

// BehaviourNote represents a Classman behaviour note
type BehaviourNote struct {
	BaseModel
	Note      string    `gorm:"not null"`
	Lesson    Lesson    `gorm:"constraint:OnDelete:CASCADE"`
	LessonID  uuid.UUID `gorm:"index:; not null"`
	Student   Student   `gorm:"constraint:OnDelete:CASCADE"`
	StudentID uuid.UUID `gorm:"index:; not null"`
}
