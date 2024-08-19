package dbinterface

import (
	"github.com/ray31245/seo_cluster/pkg/db/model"
)

type DBInterface interface{}

type SiteDAOInterface interface {
	CreateCategory(category *model.Category) error
	CreateSite(site *model.Site) (model.Site, error)
	ListSites() ([]model.Site, error)
	FirstPublishedCategory() (*model.Category, error)
	LastPublishedCategory() (*model.Category, error)
	MarkPublished(categoryID string) error
	IncreaseLackCount(siteID string, count int) error
}

type ArticleCacheDAOInterface interface {
	AddArticleToCache(article model.ArticleCache) error
	ListArticleCacheByLimit(limit int) ([]model.ArticleCache, error)
	DeleteArticleCache(id string) error
}
