package handler

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"unicode/utf8"

	"github.com/ray31245/seo_cluster/cmd/publish_manager_service/model"
	"github.com/ray31245/seo_cluster/pkg/util"
	commentbot "github.com/ray31245/seo_cluster/service/comment_bot"
	rewritemanager "github.com/ray31245/seo_cluster/service/rewrite_manager"
	sitemanager "github.com/ray31245/seo_cluster/service/site_manager"
	usermanager "github.com/ray31245/seo_cluster/service/user_manager"

	"github.com/gin-gonic/gin"
)

const (
	minSrcLength = 2000
	minArtLength = 500
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
	rewritemanager *rewritemanager.RewriteManager
}

func NewRewriteHandler(rewritemanager *rewritemanager.RewriteManager) *RewriteHandler {
	return &RewriteHandler{
		rewritemanager: rewritemanager,
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
	if utf8.RuneCount(req) < minSrcLength {
		log.Println("data is not complete")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", "data is not complete"),
		})

		return
	}

	art, err := r.rewritemanager.DefaultRewriteUntil(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	if utf8.RuneCountInString(art.Content) < minArtLength {
		extArt, err := r.rewritemanager.DefaultExtendRewriteUntil(c, req)
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

func (r *RewriteHandler) GetDefaultSystemPromptHandler(c *gin.Context) {
	systemPrompt, err := r.rewritemanager.GetDefaultSystemPrompt()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"data":    systemPrompt,
	})
}

func (r *RewriteHandler) SetDefaultSystemPromptHandler(c *gin.Context) {
	f := func(req model.SetDefaultSystemPromptRequest) error {
		return r.rewritemanager.SetDefaultSystemPrompt(req.Prompt)
	}

	code, err := setDefaultRewritePrompt(c, f)
	if err != nil {
		log.Println(err)
		c.JSON(code, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func (r *RewriteHandler) GetDefaultPromptHandler(c *gin.Context) {
	prompt, err := r.rewritemanager.GetDefaultPrompt()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"data":    prompt,
	})
}

func (r *RewriteHandler) SetDefaultPromptHandler(c *gin.Context) {
	f := func(req model.SetDefaultPromptRequest) error {
		return r.rewritemanager.SetDefaultPrompt(req.Prompt)
	}

	code, err := setDefaultRewritePrompt(c, f)
	if err != nil {
		log.Println(err)
		c.JSON(code, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func (r *RewriteHandler) GetDefaultExtendSystemPromptHandler(c *gin.Context) {
	systemPrompt, err := r.rewritemanager.GetDefaultExtendSystemPrompt()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"data":    systemPrompt,
	})
}

func (r *RewriteHandler) SetDefaultExtendSystemPromptHandler(c *gin.Context) {
	f := func(req model.SetDefaultExtendSystemPromptRequest) error {
		return r.rewritemanager.SetDefaultExtendSystemPrompt(req.Prompt)
	}

	code, err := setDefaultRewritePrompt(c, f)
	if err != nil {
		log.Println(err)
		c.JSON(code, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func (r *RewriteHandler) GetDefaultExtendPromptHandler(c *gin.Context) {
	prompt, err := r.rewritemanager.GetDefaultExtendPrompt()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"data":    prompt,
	})
}

func (r *RewriteHandler) SetDefaultExtendPromptHandler(c *gin.Context) {
	f := func(req model.SetDefaultExtendPromptRequest) error {
		return r.rewritemanager.SetDefaultExtendPrompt(req.Prompt)
	}

	code, err := setDefaultRewritePrompt(c, f)
	if err != nil {
		log.Println(err)
		c.JSON(code, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

type setPromptRequest interface {
	// model.SetDefaultSystemPromptRequest | model.SetDefaultPromptRequest | model.SetDefaultExtendSystemPromptRequest | model.SetDefaultExtendPromptRequest
	GetPrompt() string
}

type setPromptF[T setPromptRequest] func(T) error

func setDefaultRewritePrompt[T setPromptRequest](c *gin.Context, f setPromptF[T]) (int, error) {
	var req T

	err := c.ShouldBind(&req)
	if err != nil {
		log.Println(err)

		return http.StatusBadRequest, fmt.Errorf("setPromptRequest: %w", err)
	}

	if req.GetPrompt() == "" {
		log.Println("data is not complete")

		return http.StatusBadRequest, fmt.Errorf("setPromptRequest: %w", errors.New("data is not complete"))
	}

	err = f(req)
	if err != nil {
		log.Println(err)

		return http.StatusInternalServerError, fmt.Errorf("setPromptRequest: %w", err)
	}

	return http.StatusOK, nil
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
		"status": status,
	})
}
