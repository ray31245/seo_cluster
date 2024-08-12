package model

import (
	"time"

	"github.com/google/uuid"
)

type Base struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
