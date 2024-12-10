package zinterface

import (
	"context"

	"github.com/ray31245/seo_cluster/pkg/z_blog_api/model"

	"github.com/google/uuid"
)

type ZBlogAPI interface {
	GetClient(ctx context.Context, ID uuid.UUID, urlStr string, userName string, password string) (ZBlogAPIClient, error)
	UpdateClient(ctx context.Context, ID uuid.UUID, urlStr string, userName string, password string) (ZBlogAPIClient, error)
	DeleteClient(ID uuid.UUID)
	NewClient(ctx context.Context, urlStr string, userName string, password string) (ZBlogAPIClient, error)
	NewAnonymousClient(ctx context.Context, urlStr string) ZBlogAPIClient
}

type ZBlogAPIClient interface {
	ListCategory(ctx context.Context) ([]model.Category, error)
	GetArticle(ctx context.Context, id string) (model.Article, error)
	ListArticle(ctx context.Context, req model.ListArticleRequest) ([]model.Article, error)
	PostArticle(ctx context.Context, art model.PostArticleRequest) (model.Article, error)
	PostComment(ctx context.Context, comment model.PostCommentRequest) error
	GetCountOfArticle(ctx context.Context, req model.ListArticleRequest) (int, error)
	ListTag(ctx context.Context, req model.ListTagRequest) ([]model.Tag, error)
	ListTagAll(ctx context.Context) ([]model.Tag, error)
	PostTag(ctx context.Context, tag model.PostTagRequest) (model.Tag, error)
}
