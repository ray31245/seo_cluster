package db

import (
	"fmt"
	"log"
	"slices"

	dbErr "github.com/ray31245/seo_cluster/pkg/db/error"
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
	if article.Status == "" {
		article.Status = model.ArticleCacheStatusDefault
	}

	if !slices.Contains(model.ArticleCacheStatues, article.Status) {
		return fmt.Errorf("AddArticleToCache: %w", dbErr.ErrInvalidArticleCacheStatus)
	}

	return d.db.Create(&article).Error
}

func (d *ArticleCacheDAO) ListReadyToPublishArticleCacheByLimit(limit int) ([]model.ArticleCache, error) {
	var articles []model.ArticleCache
	err := d.db.Where("status != ?", model.ArticleCacheStatusReserved).Limit(limit).Order("created_at").Find(&articles).Error

	if len(articles) < limit {
		log.Println("ArticleCacheDAO: ListArticleCacheByLimit: less than limit")
	}

	return articles, err
}

func (d *ArticleCacheDAO) ListPublishLaterArticleCachePaginator(titleKeyword, contentKeyword string, op model.Operator, page int, limit int) ([]model.ArticleCache, int, int64, error) {
	var (
		articles  []model.ArticleCache
		totalPage int
	)

	q := d.db.Where("status = ?", model.ArticleCacheStatusDefault)

	q = articleFuzzySearch(q, titleKeyword, contentKeyword, op)

	articles, totalPage, totalRows, err := paginator(model.ArticleCache{}, q, page, limit)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("ListPublishLaterArticleCachePaginator: %w", err)
	}

	return articles, totalPage, totalRows, err
}

func (d *ArticleCacheDAO) ListEditAbleArticleCachePaginator(titleKeyword, contentKeyword string, op model.Operator, page int, limit int) ([]model.ArticleCache, int, int64, error) {
	var articles []model.ArticleCache

	q := d.db.Where("status = ?", model.ArticleCacheStatusReserved)

	q = articleFuzzySearch(q, titleKeyword, contentKeyword, op)

	articles, totalPage, totalRows, err := paginator(model.ArticleCache{}, q, page, limit)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("ListEditAbleArticleCachePaginator: %w", err)
	}

	return articles, totalPage, totalRows, err
}

func (d *ArticleCacheDAO) GetArticleCacheByID(id string) (*model.ArticleCache, error) {
	article := model.ArticleCache{}

	err := d.db.Where("id = ?", id).First(&article).Error
	if err != nil {
		return nil, fmt.Errorf("GetArticleCacheByID: %w", err)
	}

	return &article, nil
}

func (d *ArticleCacheDAO) DeleteArticleCacheByIDs(ids []string) error {
	return d.db.Where("id IN ?", ids).Delete(&model.ArticleCache{}).Error
}

func (d *ArticleCacheDAO) CountArticleCache() (int64, error) {
	var count int64
	err := d.db.Model(&model.ArticleCache{}).Count(&count).Error

	return count, err
}

func (d *ArticleCacheDAO) EditArticleCache(id string, title string, content string) error {
	return d.db.Model(&model.ArticleCache{}).Where("id = ?", id).Updates(map[string]interface{}{"title": title, "content": content}).Error
}

func (d *ArticleCacheDAO) UpdateArticleCacheStatusByIDs(ids []string, status model.ArticleCacheStatus) error {
	if !slices.Contains(model.ArticleCacheStatues, status) {
		return fmt.Errorf("AddArticleToCache: %w", dbErr.ErrInvalidArticleCacheStatus)
	}

	return d.db.Model(&model.ArticleCache{}).Where("id IN ?", ids).Update("status", status).Error
}
