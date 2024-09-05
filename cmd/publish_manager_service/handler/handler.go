package handler

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/ray31245/seo_cluster/cmd/publish_manager_service/model"
	aiAssistInterface "github.com/ray31245/seo_cluster/pkg/ai_assist/ai_assist_interface"
	"github.com/ray31245/seo_cluster/pkg/util"
	publishManager "github.com/ray31245/seo_cluster/service/publish_manager"
	sitemanager "github.com/ray31245/seo_cluster/service/site_manager"

	"github.com/gin-gonic/gin"
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

	req.Content, err = util.DecodeImageListDivFromHTMl([]byte(req.Content))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
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

func (p *PublishHandler) PrePublishHandler(c *gin.Context) {
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

	req.Content, err = util.DecodeImageListDivFromHTMl([]byte(req.Content))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
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

func (p *PublishHandler) FlexiblePublishHandler(c *gin.Context) {
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

	req.Content, err = util.DecodeImageListDivFromHTMl([]byte(req.Content))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
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

type SiteHandler struct {
	sitemanager *sitemanager.SiteManager
}

func NewSiteHandler(sitemanager *sitemanager.SiteManager) *SiteHandler {
	return &SiteHandler{
		sitemanager: sitemanager,
	}
}

func (s *SiteHandler) AddSiteHandler(c *gin.Context) {
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

	err = s.sitemanager.AddSite(c, req.URL, req.UserName, req.Password)
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

func (s *SiteHandler) ListSitesHandler(c *gin.Context) {
	sites, err := s.sitemanager.ListSites()
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

type RewriteHandler struct {
	aiAssist aiAssistInterface.AIAssistInterface
}

func NewRewriteHandler(aiAssist aiAssistInterface.AIAssistInterface) *RewriteHandler {
	return &RewriteHandler{
		aiAssist: aiAssist,
	}
}

func (r *RewriteHandler) RewriteHandler(c *gin.Context) {
	// get data body from request
	req, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	// check data
	if len(req) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", "data is not complete"),
		})

		return
	}

	art, err := r.aiAssist.Rewrite(c, req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	art.Content = string(util.MdToHTML([]byte(art.Content)))

	imgDiv, err := util.GenImageListEncodeDiv(req)
	if err != nil {
		log.Printf("error: %v", err)
	} else {
		art.Content += imgDiv
	}

	log.Printf("%s", art.Content)

	c.JSON(http.StatusOK, art)
}
