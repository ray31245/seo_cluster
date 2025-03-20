package articlecachemanager

import (
	dbInterface "github.com/ray31245/seo_cluster/pkg/db/db_interface"
	dbModel "github.com/ray31245/seo_cluster/pkg/db/model"
)

type ArticleCacheManager struct {
	dbInterface.ArticleCacheDAOInterface
}

func NewArticleCacheManager(dao dbInterface.ArticleCacheDAOInterface) *ArticleCacheManager {
	return &ArticleCacheManager{
		ArticleCacheDAOInterface: dao,
	}
}

func (a *ArticleCacheManager) ListPublishLaterArticleCache(titleKeyword, contentKeyword string, op dbModel.Operator, page, limit int) ([]dbModel.ArticleCache, int, int64, error) {
	return a.ArticleCacheDAOInterface.ListPublishLaterArticleCachePaginator(titleKeyword, contentKeyword, op, page, limit)
}

func (a *ArticleCacheManager) ListEditAbleArticleCache(titleKeyword, contentKeyword string, op dbModel.Operator, page, limit int) ([]dbModel.ArticleCache, int, int64, error) {
	return a.ArticleCacheDAOInterface.ListEditAbleArticleCachePaginator(titleKeyword, contentKeyword, op, page, limit)
}

func (a *ArticleCacheManager) GetArticleCache(id string) (*dbModel.ArticleCache, error) {
	return a.ArticleCacheDAOInterface.GetArticleCacheByID(id)
}

func (a *ArticleCacheManager) DeleteArticleCache(IDs []string) error {
	return a.ArticleCacheDAOInterface.DeleteArticleCacheByIDs(IDs)
}

func (a *ArticleCacheManager) UpdateArticleCacheStatus(IDs []string, status dbModel.ArticleCacheStatus) error {
	return a.ArticleCacheDAOInterface.UpdateArticleCacheStatusByIDs(IDs, status)
}

func (a *ArticleCacheManager) EditArticleCache(id, title, content string) error {
	return a.ArticleCacheDAOInterface.EditArticleCache(id, title, content)
}
