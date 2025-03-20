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

func (p *PublishHandler) DirectPublishHandler(c *gin.Context) {
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

	cateID := c.Param("cateID")

	err = p.publisher.DirectPublish(c, cateID, req.ToPublishManager())
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

func (p *PublishHandler) BroadcastPublishHandler(c *gin.Context) {
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

	err = p.publisher.BroadcastPublish(c, req.ToPublishManager())
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

func (p *PublishHandler) SpecifyPublishHandler(c *gin.Context) {
	req := model.SpecifyPublishRequest{}

	err := c.ShouldBind(&req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	// check data
	if req.ArticleID == "" || req.CateID == "" {
		log.Println("data is not complete")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", "data is not complete"),
		})

		return
	}

	err = p.publisher.SpecifyPublish(c, req.CateID, req.ArticleID)
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

func (p *PublishHandler) SetConfigUnCategoryNameHandler(c *gin.Context) {
	// get data body from request
	req := model.SetUnconfigCategoryNameRequest{}

	err := c.ShouldBind(&req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	// check data
	if req.Name == "" {
		log.Println("data is not complete")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", "data is not complete"),
		})

		return
	}

	err = p.publisher.SetConfigUnCateName(c, req.Name)
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

func (p *PublishHandler) GetConfigUnCateNameHandler(c *gin.Context) {
	name, err := p.publisher.GetConfigUnCateName()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name": name,
	})
}

func (p *PublishHandler) SetConfigTagBlackList(c *gin.Context) {
	// get data body from request
	req := model.SetConfigTagBlackListRequest{}

	err := c.ShouldBind(&req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	// check data
	if req.Tags == nil {
		log.Println("data is not complete")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", "data is not complete"),
		})

		return
	}

	err = p.publisher.SetTagsBlockList(c, req.Tags)
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

func (p *PublishHandler) GetConfigTagBlackList(c *gin.Context) {
	tags, err := p.publisher.GetTagsBlockList()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tags": tags,
	})
}

func (p *PublishHandler) StopAutoPublishHandler(c *gin.Context) {
	err := p.publisher.StopAutoPublish()
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

func (p *PublishHandler) StartAutoPublishHandler(c *gin.Context) {
	err := p.publisher.StartAutoPublish()
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

func (p *PublishHandler) GetStopAutoPublishStatusHandler(c *gin.Context) {
	status, err := p.publisher.StopAutoPublishStatus()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": status,
	})
}
