package zinterface

import (
	"context"

	"goTool/pkg/z_blog_api/model"

	"github.com/google/uuid"
)

type ZBlogAPI interface {
	GetClient(ctx context.Context, siteID uuid.UUID, urlStr string, userName string, password string) (ZBlogAPIClient, error)
	NewClient(ctx context.Context, urlStr string, userName string, password string) (ZBlogAPIClient, error)
}

type ZBlogAPIClient interface {
	ListCategory(ctx context.Context) ([]model.Category, error)
	PostArticle(ctx context.Context, art model.PostArticleRequest) error
	GetCountOfArticle(ctx context.Context, req model.ListArticleRequest) (int, error)
}
