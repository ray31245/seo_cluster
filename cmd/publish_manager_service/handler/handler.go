package handler

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/ray31245/seo_cluster/cmd/publish_manager_service/model"
	aiAssistInterface "github.com/ray31245/seo_cluster/pkg/ai_assist/ai_assist_interface"
	aiAssistModel "github.com/ray31245/seo_cluster/pkg/ai_assist/model"
	"github.com/ray31245/seo_cluster/pkg/util"
	commentbot "github.com/ray31245/seo_cluster/service/comment_bot"
	sitemanager "github.com/ray31245/seo_cluster/service/site_manager"
	usermanager "github.com/ray31245/seo_cluster/service/user_manager"

	"github.com/gin-gonic/gin"
)

const (
	retryLimit = 50
	retryDelay = 100 * time.Millisecond
)

// TODO: error handling on middleware

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
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	// check data
	if req.URL == "" || req.UserName == "" || req.Password == "" {
		log.Println("data is not complete")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", "data is not complete"),
		})

		return
	}

	err = s.sitemanager.AddSite(c, req.CMSType, req.URL, req.UserName, req.Password, req.ExpectCategoryNum)
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

func (s *SiteHandler) DeleteSiteHandler(c *gin.Context) {
	id := c.Param("siteID")

	err := s.sitemanager.DeleteSite(id)
	if err != nil {
		log.Println(err)

		errCode := http.StatusInternalServerError
		if errors.Is(err, sitemanager.ErrSiteNotFound) {
			errCode = http.StatusNotFound
		}

		c.JSON(errCode, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func (s *SiteHandler) UpdateSiteHandler(c *gin.Context) {
	// get data body from request
	req := model.UpdateSiteRequest{}

	err := c.ShouldBind(&req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	// check data
	if req.SiteID == "" {
		log.Println("data is not complete")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", "data is not complete"),
		})

		return
	}

	err = s.sitemanager.UpdateSite(c, req.SiteID, req.URL, req.UserName, req.Password)
	if err != nil {
		log.Println(err)

		errCode := http.StatusInternalServerError
		if errors.Is(err, sitemanager.ErrSiteNotFound) {
			errCode = http.StatusNotFound
		}

		c.JSON(errCode, gin.H{
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
		log.Println(err)
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

func (s *SiteHandler) GetSiteHandler(c *gin.Context) {
	id := c.Param("siteID")

	site, err := s.sitemanager.GetSite(id)
	if err != nil {
		log.Println(err)

		errCode := http.StatusInternalServerError
		if errors.Is(err, sitemanager.ErrSiteNotFound) {
			errCode = http.StatusNotFound
		}

		c.JSON(errCode, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	data := model.GetSiteResponse{}
	data.FromDBSite(*site)

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"data":    data,
	})
}

func (s *SiteHandler) SyncCategoryFromAllSiteHandler(c *gin.Context) {
	err := s.sitemanager.SyncCategoryFromAllSite(c)
	if err != nil {
		log.Println(err)

		errCode := http.StatusMultiStatus
		c.JSON(errCode, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func (s *SiteHandler) SyncCategoryFromSiteHandler(c *gin.Context) {
	id := c.Param("siteID")

	err := s.sitemanager.SyncCategoryFromSite(c, id)
	if err != nil {
		log.Println(err)

		errCode := http.StatusInternalServerError
		if errors.Is(err, sitemanager.ErrSiteNotFound) {
			errCode = http.StatusNotFound
		}

		c.JSON(errCode, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func (s *SiteHandler) IncreaseLackCountHandler(c *gin.Context) {
	req := model.IncreaseLackCountRequest{}

	err := c.ShouldBind(&req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	err = s.sitemanager.IncreaseLackCount(req.SiteID, req.Count)
	if err != nil {
		log.Println(err)

		errCode := http.StatusInternalServerError
		if errors.Is(err, sitemanager.ErrSiteNotFound) {
			errCode = http.StatusNotFound
		}

		c.JSON(errCode, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
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
		log.Println("data is not complete")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", "data is not complete"),
		})

		return
	}

	art, err := r.rewriteUntil(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	if utf8.RuneCountInString(art.Content) < 500 {
		extArt, err := r.extendRewriteUntil(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("error: %v", err),
			})

			return
		}

		art.Content = extArt.Content
		art.Title = extArt.Title
	}

	art.Content = string(util.MdToHTML([]byte(art.Content)))

	imgDiv, err := util.GenImageListEncodeDiv(req)
	if err != nil {
		log.Printf("error: %v", err)
	} else {
		art.Content += imgDiv
	}

	// log.Printf("%s", art.Content)

	c.JSON(http.StatusOK, art)
}

func (r *RewriteHandler) rewriteUntil(c *gin.Context, req []byte) (art aiAssistModel.RewriteResponse, err error) {
	log.Println("rewriting...")

	r.aiAssist.Lock()
	defer r.aiAssist.Unlock()

	for range retryLimit {
		art, err = r.aiAssist.Rewrite(c, req)
		if err == nil {
			break
		}

		log.Println("retrying...")
		<-time.After(retryDelay)
	}

	if err != nil {
		return art, fmt.Errorf("rewriteUntil: retry limit: %w", err)
	}

	return
}

func (r *RewriteHandler) extendRewriteUntil(c *gin.Context, req []byte) (art aiAssistModel.ExtendRewriteResponse, err error) {
	log.Println("extending rewriting...")

	r.aiAssist.Lock()
	defer r.aiAssist.Unlock()

	for range retryLimit {
		art, err = r.aiAssist.ExtendRewrite(c, req)
		if err == nil {
			break
		}

		log.Println("retrying...")
		<-time.After(retryDelay)
	}

	if err != nil {
		return art, fmt.Errorf("extendRewriteUntil: retry limit: %w", err)
	}

	return
}

type UserHandler struct {
	usermanager *usermanager.UserManager
}

func NewUserHandler(usermanager *usermanager.UserManager) *UserHandler {
	return &UserHandler{
		usermanager: usermanager,
	}
}

func (u *UserHandler) AddFirstAdminUser(c *gin.Context) {
	// get data body from request
	req := model.AddFirstAdminUserRequest{}

	err := c.ShouldBind(&req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	// check data
	if req.UserName == "" || req.Password == "" {
		log.Println("data is not complete")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", "data is not complete"),
		})

		return
	}

	err = u.usermanager.CreateFirstAdminUser(c, req.UserName, req.Password)
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

type CommentBotHandler struct {
	commentBot *commentbot.CommentBot
}

func NewCommentBotHandler(commentBot *commentbot.CommentBot) *CommentBotHandler {
	return &CommentBotHandler{
		commentBot: commentBot,
	}
}

func (b *CommentBotHandler) StartAutoCommentHandler(c *gin.Context) {
	err := b.commentBot.StartAutoComment(c)
	if err != nil {
		log.Println(err)

		errCode := http.StatusInternalServerError
		c.JSON(errCode, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func (b *CommentBotHandler) StopAutoCommentHandler(c *gin.Context) {
	err := b.commentBot.StopAutoComment(c)
	if err != nil {
		log.Println(err)

		errCode := http.StatusInternalServerError
		c.JSON(errCode, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func (b *CommentBotHandler) GetStopAutoCommentStatusHandler(c *gin.Context) {
	status, err := b.commentBot.IsAutoCommentStopped()
	if err != nil {
		log.Println(err)

		errCode := http.StatusInternalServerError
		c.JSON(errCode, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": status,
	})
}
