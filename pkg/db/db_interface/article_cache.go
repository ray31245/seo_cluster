package dbinterface

import (
	"github.com/ray31245/seo_cluster/pkg/db/model"
)

type ArticleCacheDAOInterface interface {
	AddArticleToCache(article model.ArticleCache) error
	ListArticleCacheByLimit(limit int) ([]model.ArticleCache, error)
	DeleteArticleCache(id string) error
	CountArticleCache() (int64, error)
}
