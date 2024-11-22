package db

import (
	"fmt"
	"log"

	"github.com/ray31245/seo_cluster/pkg/db/model"

	"gorm.io/gorm"
)

type ArticleCacheDAO struct {
	db *gorm.DB
}

func (d *DB) NewArticleCacheDAO() (*ArticleCacheDAO, error) {
	err := d.db.AutoMigrate(&model.ArticleCache{})
	if err != nil {
		return nil, fmt.Errorf("NewArticleCacheDAO: %w", err)
	}

	return &ArticleCacheDAO{db: d.db}, nil
}

func (d *ArticleCacheDAO) AddArticleToCache(article model.ArticleCache) error {
	return d.db.Create(&article).Error
}

func (d *ArticleCacheDAO) ListArticleCacheByLimit(limit int) ([]model.ArticleCache, error) {
	var articles []model.ArticleCache
	err := d.db.Limit(limit).Order("created_at").Find(&articles).Error

	if len(articles) < limit {
		log.Println("ArticleCacheDAO: ListArticleCacheByLimit: less than limit")
	}

	return articles, err
}

func (d *ArticleCacheDAO) DeleteArticleCache(id string) error {
	return d.db.Delete(&model.ArticleCache{}, fmt.Sprintf("id = '%s'", id)).Error
}

func (d *ArticleCacheDAO) CountArticleCache() (int64, error) {
	var count int64
	err := d.db.Model(&model.ArticleCache{}).Count(&count).Error

	return count, err
}
