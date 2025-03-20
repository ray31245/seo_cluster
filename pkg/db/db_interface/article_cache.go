package dbinterface

import (
	"github.com/ray31245/seo_cluster/pkg/db/model"
)

type ArticleCacheDAOInterface interface {
	AddArticleToCache(article model.ArticleCache) error
	ListReadyToPublishArticleCacheByLimit(limit int) ([]model.ArticleCache, error)
	ListPublishLaterArticleCachePaginator(titleKeyword, contentKeyword string, op model.Operator, page int, limit int) ([]model.ArticleCache, int, int64, error)
	ListEditAbleArticleCachePaginator(titleKeyword, contentKeyword string, op model.Operator, page int, limit int) ([]model.ArticleCache, int, int64, error)
	GetArticleCacheByID(id string) (*model.ArticleCache, error)
	DeleteArticleCacheByIDs(ids []string) error
	CountArticleCache() (int64, error)
	EditArticleCache(id string, title string, content string) error
	UpdateArticleCacheStatusByIDs(ids []string, status model.ArticleCacheStatus) error
}
