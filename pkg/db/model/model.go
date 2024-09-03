package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (b *Base) BeforeCreate(tx *gorm.DB) error { //nolint: revive
	// generate uuid if not exist
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	// b.CreatedAt = time.Now()
	// b.UpdatedAt = time.Now()
	return nil
}
