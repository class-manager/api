package model

import (
	"time"

	"github.com/gofrs/uuid"
)

// BaseModel is the base override model for gorm's Model
type BaseModel struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time `gorm:"autoCreateTime:nano"`
	// UpdatedAt time.Time `gorm:"autoUpdateTime:nano"`
	// DeletedAt gorm.DeletedAt `gorm:"index"`
}
