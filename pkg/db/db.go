package db

import (
	"errors"
	"fmt"
	dbErr "goTool/pkg/db/error"
	"goTool/pkg/db/model"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	db *gorm.DB
}

// note: dsn is "publish_manager.db"

func NewDB(dsn string) (*DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &DB{db: db}, nil
}

func (d *DB) Close() error {
	db, err := d.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func (d *DB) Migrate() error {
	return d.db.AutoMigrate(
		&model.Site{},
		&model.Category{},
		&model.ArticleCache{},
	)
}

type DAO struct {
	db *gorm.DB
}

func (d *DB) NewDAO() *DAO {
	return &DAO{db: d.db}
}

func (d *DAO) CreateSite(site *model.Site) (model.Site, error) {
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

func (d *DAO) ListSites() ([]model.Site, error) {
	var sites []model.Site
	err := d.db.Find(&sites).Error
	return sites, err
}

func (d *DAO) CreateCategory(category *model.Category) error {
	return d.db.Create(category).Error
}

func (d *DAO) UpdateSite(site *model.Site) error {
	return d.db.Save(site).Error
}

func (d *DAO) UpdateCategory(category *model.Category) error {
	return d.db.Save(category).Error
}

func (d *DAO) DeleteSite(siteID string) error {
	return d.db.Delete(&model.Site{}, siteID).Error
}

func (d *DAO) DeleteCategory(categoryID string) error {
	return d.db.Delete(&model.Category{}, categoryID).Error
}

func (d *DAO) GetSite(siteID string) (*model.Site, error) {
	var site model.Site
	err := d.db.Preload("Categorys").First(&site, siteID).Error
	return &site, err
}

func (d *DAO) GetCategory(categoryID string) (*model.Category, error) {
	var category model.Category
	err := d.db.First(&category, fmt.Sprintf("id = '%s'", categoryID)).Error
	return &category, err
}

func (d *DAO) FirstPublishedCategory() (*model.Category, error) {
	var category model.Category
	err := d.db.
		Where("exists (select 1 from sites where sites.id = categories.site_id and sites.lack_count != 0)"). // nolint:lll
		Preload("Site").Order("last_published").First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = dbErr.ErrNoCategoryNeedToBePublished
		}
		return nil, fmt.Errorf("FirstPublishedCategory: %w", err)
	}
	return &category, nil
}

func (d *DAO) LastPublishedCategory() (*model.Category, error) {
	var category model.Category
	err := d.db.Preload("Site").Order("last_published desc").First(&category).Error
	if err != nil {
		return nil, fmt.Errorf("LastPublishedCategory: %w", err)
	}
	return &category, nil
}

func (d *DAO) MarkPublished(categoryID string) error {
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

func (d *DAO) IncreaseLackCount(siteID string, count int) error {
	return d.db.Model(&model.Site{}).Where("id = ?", siteID).Update("lack_count", gorm.Expr("lack_count + ?", count)).Error
}

func (d *DAO) AddArticleToCache(article model.ArticleCache) error {
	return d.db.Create(&article).Error
}

func (d *DAO) ListArticleCacheByLimit(limit int) ([]model.ArticleCache, error) {
	var articles []model.ArticleCache
	err := d.db.Limit(limit).Order("created_at").Find(&articles).Error
	return articles, err
}

func (d *DAO) DeleteArticleCache(id string) error {
	return d.db.Delete(&model.ArticleCache{}, fmt.Sprintf("id = '%s'", id)).Error
}
