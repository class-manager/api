package model

import (
	"time"

	"github.com/gofrs/uuid"
)

// Student represents a Classman student
type Student struct {
	BaseModel
	FirstName       string    `gorm:"not null"`
	LastName        string    `gorm:"not null"`
	DOB             time.Time `gorm:"type:date"`
	GraduatingClass uint32    `gorm:"not null"`
	GeneralNote     string
	StudentNumber   string
	CreatedBy       Account   // The account that created this student
	CreatedByID     uuid.UUID `gorm:"index:; not null"`
	Classes         []*Class  `gorm:"many2many:students_classes"`
}
