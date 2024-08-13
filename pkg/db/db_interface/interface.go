package dbinterface

import (
	"goTool/pkg/db/model"
)

type DBInterface interface{}

type DAOInterface interface {
	CreateCategory(category *model.Category) error
	CreateSite(site *model.Site) (model.Site, error)
	FirstPublishedCategory() (*model.Category, error)
	MarkPublished(categoryID string) error
	AddArticleToCache(article model.ArticleCache) error
}
