package handler

import (
	"errors"
	"fmt"
	"goTool/pkg/db"
	publishmanager "goTool/pkg/publish_manager"
	zblogapi "goTool/pkg/z_blog_api"
	zModel "goTool/pkg/z_blog_api/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	DAO  *db.DAO
	ZApi *zblogapi.ZblogAPI
}

func NewHandler(dao *db.DAO, zApi *zblogapi.ZblogAPI) *Handler {
	return &Handler{
		DAO:  dao,
		ZApi: zApi,
	}
}

func (p *Handler) AveragePublishHandler(c *gin.Context) {
	// get data body from request
	article := zModel.PostArticleRequest{}
	c.ShouldBind(&article)

	// check data
	if article.Title == "" || article.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", "data is not complete"),
		})
		return
	}

	publishmanager := publishmanager.NewPublishManager(p.ZApi, p.DAO)
	err := publishmanager.AveragePublish(article)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func (p *Handler) PrePublishHandler(c *gin.Context) {
	// get data body from request
	article := zModel.PostArticleRequest{}
	c.ShouldBind(&article)

	// check data
	if article.Title == "" || article.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", "data is not complete"),
		})
		return
	}

	publishmanager := publishmanager.NewPublishManager(p.ZApi, p.DAO)
	err := publishmanager.PrePublish(article)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func (p *Handler) FlexiblePublishHandler(c *gin.Context) {
	// get data body from request
	article := zModel.PostArticleRequest{}
	c.ShouldBind(&article)

	// check data
	if article.Title == "" || article.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", "data is not complete"),
		})
		return
	}

	publisher := publishmanager.NewPublishManager(p.ZApi, p.DAO)
	err := publisher.AveragePublish(article)
	if errors.Is(err, publishmanager.ErrNoCategoryNeedToBePublished) {
		err = publisher.PrePublish(article)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("error: %v", err),
			})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})
		return

	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

type AddSiteRequest struct {
	URL      string `json:"url"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

func (p *Handler) AddSiteHandler(c *gin.Context) {
	// get data body from request
	req := AddSiteRequest{}
	c.ShouldBind(&req)

	// check data
	if req.URL == "" || req.UserName == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", "data is not complete"),
		})
		return
	}

	pub := publishmanager.NewPublishManager(p.ZApi, p.DAO)
	err := pub.AddSite(req.URL, req.UserName, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
