package model

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	Base
	ZblogID       uint32    `json:"zblog_id"`
	SiteID        uuid.UUID `json:"site_id"`
	Site          Site      `json:"site"`
	LastPublished time.Time `json:"last_published"`
}
