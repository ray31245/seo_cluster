package db

import (
	"fmt"

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

	return articles, err
}

func (d *ArticleCacheDAO) DeleteArticleCache(id string) error {
	return d.db.Delete(&model.ArticleCache{}, fmt.Sprintf("id = '%s'", id)).Error
}
