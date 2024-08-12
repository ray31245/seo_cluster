package zinterface

import (
	"goTool/pkg/z_blog_api/model"

	"github.com/google/uuid"
)

type ZBlogApi interface {
	GetClient(siteID uuid.UUID, urlStr string, userName string, password string) (ZBlogApiClient, error)
	NewClient(urlStr string, userName string, password string) (ZBlogApiClient, error)
}

type ZBlogApiClient interface {
	ListCategory() ([]model.Category, error)
	PostArticle(art model.PostArticleRequest) error
}
