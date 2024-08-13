package model

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	log.Println("before create")
	b.ID = uuid.New()
	// b.CreatedAt = time.Now()
	// b.UpdatedAt = time.Now()
	return
}
