package db

import (
	"errors"
	"fmt"
	dbErr "goTool/pkg/db/error"
	"goTool/pkg/db/model"
	"time"

	"gorm.io/gorm"
)

type SiteDAO struct {
	db *gorm.DB
}

func (d *DB) NewSiteDAO() (*SiteDAO, error) {
	err := d.db.AutoMigrate(&model.Site{}, &model.Category{})
	if err != nil {
		return nil, err
	}
	return &SiteDAO{db: d.db}, nil
}

func (d *SiteDAO) CreateSite(site *model.Site) (model.Site, error) {
	err := d.db.Create(site).Error
	if err != nil {
		return model.Site{}, err
	}
	res := model.Site{}
	err = d.db.Where("url = ?", site.URL).First(&res).Error
	if err != nil {
		return model.Site{}, err
	}
	return res, nil
}

func (d *SiteDAO) ListSites() ([]model.Site, error) {
	var sites []model.Site
	err := d.db.Find(&sites).Error
	return sites, err
}

func (d *SiteDAO) CreateCategory(category *model.Category) error {
	return d.db.Create(category).Error
}

func (d *SiteDAO) UpdateSite(site *model.Site) error {
	return d.db.Save(site).Error
}

func (d *SiteDAO) UpdateCategory(category *model.Category) error {
	return d.db.Save(category).Error
}

func (d *SiteDAO) DeleteSite(siteID string) error {
	return d.db.Delete(&model.Site{}, siteID).Error
}

func (d *SiteDAO) DeleteCategory(categoryID string) error {
	return d.db.Delete(&model.Category{}, categoryID).Error
}

func (d *SiteDAO) GetSite(siteID string) (*model.Site, error) {
	var site model.Site
	err := d.db.Preload("Categories").First(&site, siteID).Error
	return &site, err
}

func (d *SiteDAO) GetCategory(categoryID string) (*model.Category, error) {
	var category model.Category
	err := d.db.First(&category, fmt.Sprintf("id = '%s'", categoryID)).Error
	return &category, err
}

func (d *SiteDAO) FirstPublishedCategory() (*model.Category, error) {
	var category model.Category
	err := d.db.
		Where("exists (select 1 from sites where sites.id = categories.site_id and sites.lack_count != 0)").
		Preload("Site").Order("last_published").First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = dbErr.ErrNoCategoryNeedToBePublished
		}
		return nil, fmt.Errorf("FirstPublishedCategory: %w", err)
	}
	return &category, nil
}

func (d *SiteDAO) LastPublishedCategory() (*model.Category, error) {
	var category model.Category
	err := d.db.Preload("Site").Order("last_published desc").First(&category).Error
	if err != nil {
		return nil, fmt.Errorf("LastPublishedCategory: %w", err)
	}
	return &category, nil
}

func (d *SiteDAO) MarkPublished(categoryID string) error {
	cate := &model.Category{}
	err := d.db.Model(cate).Where("id = ?", categoryID).Update("last_published", time.Now()).Error
	if err != nil {
		return err
	}
	cate, err = d.GetCategory(categoryID)
	if err != nil {
		return err
	}
	err = d.db.Model(&model.Site{}).Where("id = ?", cate.SiteID).Update("lack_count", gorm.Expr("lack_count - 1")).Error
	return err
}

func (d *SiteDAO) IncreaseLackCount(siteID string, count int) error {
	return d.db.Model(&model.Site{}).Where("id = ?", siteID).Update("lack_count", gorm.Expr("lack_count + ?", count)).Error
}
