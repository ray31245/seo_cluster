package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ray31245/seo_cluster/cmd/publish_manager_service/helper"
	"github.com/ray31245/seo_cluster/cmd/publish_manager_service/model"
	dbErr "github.com/ray31245/seo_cluster/pkg/db/error"
	dbModel "github.com/ray31245/seo_cluster/pkg/db/model"
	articleCacheManager "github.com/ray31245/seo_cluster/service/article_cache_manager"
)

type articleCacheHandler struct {
	articleCacheManager *articleCacheManager.ArticleCacheManager
}

func NewArticleCacheHandler(articleCacheManager *articleCacheManager.ArticleCacheManager) *articleCacheHandler {
	return &articleCacheHandler{
		articleCacheManager: articleCacheManager,
	}
}

func (a *articleCacheHandler) GetArticleCacheHandler(c *gin.Context) {
	id := c.Param("id")

	article, err := a.articleCacheManager.GetArticleCache(id)
	if dbErr.IsNotfoundErr(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})

		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"article": article,
	})
}

func (a *articleCacheHandler) ListPublishLaterArticleCacheHandler(c *gin.Context) {
	titleKeywords, contentKeywords, operator := helper.ParseArticleCacheFuzzySearchQuery(c)

	page, pageSize, err := helper.ParsePagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	articles, totalPage, totalRows, err := a.articleCacheManager.ListPublishLaterArticleCache(titleKeywords, contentKeywords, dbModel.Operator(operator), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"articles":   articles,
		"total_page": totalPage,
		"total_rows": totalRows,
	})
}

func (a *articleCacheHandler) ListEditAbleArticleCacheHandler(c *gin.Context) {
	titleKeywords, contentKeywords, operator := helper.ParseArticleCacheFuzzySearchQuery(c)

	page, pageSize, err := helper.ParsePagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	articles, totalPage, totalRows, err := a.articleCacheManager.ListEditAbleArticleCache(titleKeywords, contentKeywords, dbModel.Operator(operator), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"articles":   articles,
		"total_page": totalPage,
		"total_rows": totalRows,
	})
}

func (a *articleCacheHandler) UpdateArticleCacheStatusHandler(c *gin.Context) {
	req := model.UpdateArticleCacheStatusRequest{}

	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	if len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "ids is empty",
		})

		return
	}

	err = a.articleCacheManager.UpdateArticleCacheStatus(req.IDs, dbModel.ArticleCacheStatus(req.Status))
	if errors.Is(err, dbErr.ErrInvalidArticleCacheStatus) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func (a *articleCacheHandler) EditArticleCacheHandler(c *gin.Context) {
	req := model.EditArticleCacheRequest{}

	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	if req.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id is empty",
		})

		return
	}

	if req.Title == "" || req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "data is not complete",
		})

		return
	}

	err = a.articleCacheManager.EditArticleCache(req.ID, req.Title, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func (a *articleCacheHandler) DeleteArticleCacheHandler(c *gin.Context) {
	req := model.DeleteArticleCacheRequest{}

	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	if len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "ids is empty",
		})

		return
	}

	err = a.articleCacheManager.DeleteArticleCache(req.IDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
