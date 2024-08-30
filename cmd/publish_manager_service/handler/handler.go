package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ray31245/seo_cluster/cmd/publish_manager_service/model"
	publishManager "github.com/ray31245/seo_cluster/service/publish_manager"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	publisher *publishManager.PublishManager
}

func NewHandler(publisher *publishManager.PublishManager) *Handler {
	return &Handler{
		publisher: publisher,
	}
}

func (p *Handler) AveragePublishHandler(c *gin.Context) {
	// get data body from request
	req := model.PublishArticleRequest{}

	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	// check data
	if req.Title == "" || req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", "data is not complete"),
		})

		return
	}

	err = p.publisher.AveragePublish(c, req.ToZBlogAPI())
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
	req := model.PublishArticleRequest{}

	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	// check data
	if req.Title == "" || req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", "data is not complete"),
		})

		return
	}

	err = p.publisher.PrePublish(req.ToZBlogAPI())
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
	req := model.PublishArticleRequest{}

	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	// check data
	if req.Title == "" || req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", "data is not complete"),
		})

		return
	}

	err = p.publisher.AveragePublish(c, req.ToZBlogAPI())
	if errors.Is(err, publishManager.ErrNoCategoryNeedToBePublished) {
		err = p.publisher.PrePublish(req.ToZBlogAPI())
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

func (p *Handler) AddSiteHandler(c *gin.Context) {
	// get data body from request
	req := model.AddSiteRequest{}

	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	// check data
	if req.URL == "" || req.UserName == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", "data is not complete"),
		})

		return
	}

	err = p.publisher.AddSite(c, req.URL, req.UserName, req.Password)
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

func (p *Handler) ListSitesHandler(c *gin.Context) {
	sites, err := p.publisher.ListSites()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	data := model.ListSitesResponse{}
	data.FromDBSites(sites)

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"data":    data.Sites,
	})
}
