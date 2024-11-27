package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ray31245/seo_cluster/cmd/publish_manager_service/model"
	"github.com/ray31245/seo_cluster/pkg/util"
	publishManager "github.com/ray31245/seo_cluster/service/publish_manager"
)

type PublishHandler struct {
	publisher *publishManager.PublishManager
}

func NewPublishHandler(publisher *publishManager.PublishManager) *PublishHandler {
	return &PublishHandler{
		publisher: publisher,
	}
}

func (p *PublishHandler) AveragePublishHandler(c *gin.Context) {
	// get data body from request
	req := model.PublishArticleRequest{}

	err := c.ShouldBind(&req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	// check data
	if req.Title == "" || req.Content == "" {
		log.Println("data is not complete")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", "data is not complete"),
		})

		return
	}

	req.Content, err = util.DecodeImageListDivFromHTMl([]byte(req.Content))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	err = p.publisher.AveragePublish(c, req.ToPublishManager())
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func (p *PublishHandler) PrePublishHandler(c *gin.Context) {
	// get data body from request
	req := model.PublishArticleRequest{}

	err := c.ShouldBind(&req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	// check data
	if req.Title == "" || req.Content == "" {
		log.Println("data is not complete")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", "data is not complete"),
		})

		return
	}

	req.Content, err = util.DecodeImageListDivFromHTMl([]byte(req.Content))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	err = p.publisher.PrePublish(req.ToPublishManager())
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func (p *PublishHandler) FlexiblePublishHandler(c *gin.Context) {
	// get data body from request
	req := model.PublishArticleRequest{}

	err := c.ShouldBind(&req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	// check data
	if req.Title == "" || req.Content == "" {
		log.Println("data is not complete")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", "data is not complete"),
		})

		return
	}

	req.Content, err = util.DecodeImageListDivFromHTMl([]byte(req.Content))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	err = p.publisher.AveragePublish(c, req.ToPublishManager())
	if errors.Is(err, publishManager.ErrNoCategoryNeedToBePublished) {
		err = p.publisher.PrePublish(req.ToPublishManager())
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("error: %v", err),
			})

			return
		}
	} else if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func (p *PublishHandler) GetArticleCacheCountHandler(c *gin.Context) {
	count, err := p.publisher.CountArticleCache()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": count,
	})
}
