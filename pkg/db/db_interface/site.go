package dbinterface

import (
	"github.com/ray31245/seo_cluster/pkg/db/model"
)

type SiteDAOInterface interface {
	CreateCategory(category *model.Category) error
	DeleteCategory(categoryID string) error
	UpdateCategory(category *model.Category) error
	DeleteSiteCategories(siteID string) error
	CreateSite(site *model.Site) (model.Site, error)
	ListSites() ([]model.Site, error)
	ListSitesByCMSType(cmsType model.CMSType) ([]model.Site, error)
	ListSitesRandom() ([]model.Site, error)
	GetSite(siteID string) (*model.Site, error)
	DeleteSite(siteID string) error
	UpdateSite(site *model.Site) error
	FirstPublishedCategory() (*model.Category, error)
	LastPublishedCategoryByCMSType(cmsType model.CMSType) (*model.Category, error)
	MarkPublished(categoryID string) error
	IncreaseLackCount(siteID string, count int) error
	SumLackCount() (int, error)
}
