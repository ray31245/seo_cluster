package model

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	Base
	ZBlogID       uint32    `json:"z_blog_id"`
	WordpressID   uint32    `json:"wordpress_id"`
	SiteID        uuid.UUID `json:"site_id"`
	Site          Site      `json:"site"`
	LastPublished time.Time `json:"last_published"`
}
