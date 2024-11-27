package wordpressinterface

import (
	"context"

	"github.com/google/uuid"
	"github.com/ray31245/seo_cluster/pkg/wordpress_api/model"
)

type WordpressAPI interface {
	GetClient(ctx context.Context, ID uuid.UUID, urlStr string, userName string, password string) (WordpressClient, error)
	UpdateClient(ctx context.Context, ID uuid.UUID, urlStr string, userName string, password string) (WordpressClient, error)
	DeleteClient(ID uuid.UUID)
	NewClient(ctx context.Context, urlStr string, userName string, password string) (WordpressClient, error)
	NewAnonymousClient(ctx context.Context, urlStr string) WordpressClient
}

type WordpressClient interface {
	RetrieveUserMe(ctx context.Context) (model.RetrieveUserMeResponse, error)
	ListTag(ctx context.Context, args model.ListTagArgs) (model.ListTagResponse, error)
	CreateTag(ctx context.Context, args model.CreateTagArgs) (model.CreateTagResponse, error)
	ListCategory(ctx context.Context, args model.ListCategoryArgs) (model.ListCategoryResponse, error)
	ListArticle(ctx context.Context, args model.ListArticleArgs) (model.ListArticleResponse, error)
	CreateArticle(ctx context.Context, args model.CreateArticleArgs) (model.CreateArticleResponse, error)
	RetrieveArticle(ctx context.Context, args model.RetrieveArticleArgs) (model.RetrieveArticleResponse, error)
	CreateComment(ctx context.Context, args model.CreateCommentArgs) (model.CreateCommentResponse, error)
}
