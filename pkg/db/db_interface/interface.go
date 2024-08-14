package dbinterface

import (
	"goTool/pkg/db/model"
)

type DBInterface interface{}

type DAOInterface interface {
	CreateCategory(category *model.Category) error
	CreateSite(site *model.Site) (model.Site, error)
	ListSites() ([]model.Site, error)
	FirstPublishedCategory() (*model.Category, error)
	LastPublishedCategory() (*model.Category, error)
	MarkPublished(categoryID string) error
	IncreaseLackCount(siteID string, count int) error
	AddArticleToCache(article model.ArticleCache) error
	ListArticleCacheByLimit(limit int) ([]model.ArticleCache, error)
	DeleteArticleCache(id string) error
}
