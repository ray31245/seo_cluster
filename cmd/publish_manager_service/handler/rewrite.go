package handler

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/ray31245/seo_cluster/cmd/publish_manager_service/model"
	"github.com/ray31245/seo_cluster/pkg/util"
	rewritemanager "github.com/ray31245/seo_cluster/service/rewrite_manager"
)

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

	res := struct {
		Title   string `json:"Title"`
		Content string `json:"Content"`
	}{}

	art, err := r.rewritemanager.DefaultRewriteUntil(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	if utf8.RuneCountInString(art) < minArtLength {
		extArt, err := r.rewritemanager.DefaultExtendRewriteUntil(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("error: %v", err),
			})

			return
		}

		art = extArt
	}

	title, err := r.rewritemanager.DefaultMakeTitleUntil(c, art)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	res.Title = title

	art = string(util.MdToHTML([]byte(art)))

	imgDiv, err := util.GenImageListEncodeDiv(req)
	if err != nil {
		log.Printf("error: %v", err)
	} else {
		art += imgDiv
	}

	res.Content = art

	c.JSON(http.StatusOK, res)
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

func (r *RewriteHandler) GetDefaultMakeTitleSystemPromptHandler(c *gin.Context) {
	systemPrompt, err := r.rewritemanager.GetDefaultMakeTitleSystemPrompt()
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

func (r *RewriteHandler) SetDefaultMakeTitleSystemPromptHandler(c *gin.Context) {
	f := func(req model.SetDefaultMakeTitleSystemPromptRequest) error {
		return r.rewritemanager.SetDefaultMakeTitleSystemPrompt(req.Prompt)
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

func (r *RewriteHandler) GetDefaultMakeTitlePromptHandler(c *gin.Context) {
	prompt, err := r.rewritemanager.GetDefaultMakeTitlePrompt()
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

func (r *RewriteHandler) SetDefaultMakeTitlePromptHandler(c *gin.Context) {
	f := func(req model.SetDefaultMakeTitlePromptRequest) error {
		return r.rewritemanager.SetDefaultMakeTitlePrompt(req.Prompt)
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
