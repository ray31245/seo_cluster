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

	var rewriteF rewriteF = func(req string) (string, error) {
		art, err := r.rewritemanager.DefaultRewriteUntil(c, []byte(req))
		if err != nil {
			return "", err
		}

		return art, nil
	}

	var extendRewriteF extendRewriteF = func(req string) (string, error) {
		extArt, err := r.rewritemanager.DefaultExtendRewriteUntil(c, []byte(req))
		if err != nil {
			return "", err
		}

		return extArt, nil
	}

	var makeTitleF makeTitleF = func(req string) (string, error) {
		title, err := r.rewritemanager.DefaultMakeTitleUntil(c, req)
		if err != nil {
			return "", err
		}

		return title, nil
	}

	res, err := r.rewriteWorkFlow(string(req), rewriteF, extendRewriteF, makeTitleF, true)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	c.JSON(http.StatusOK, res)
}

func (r *RewriteHandler) RewriteTestHandler(c *gin.Context) {
	req := model.RewriteTestRequest{}

	err := c.ShouldBind(&req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	steps := []string{}

	var rewriteF rewriteF = func(originArtContext string) (string, error) {
		steps = append(steps, "rewrite")

		art, err := r.rewritemanager.CustomRewrite(c, req.SystemPrompt, req.Prompt, []byte(originArtContext))
		if err != nil {
			return "", err
		}

		steps = append(steps, "rewrite done")

		return art, nil
	}

	var extendRewriteF extendRewriteF = func(originArtContext string) (string, error) {
		steps = append(steps, "extend rewrite")

		art, err := r.rewritemanager.CustomRewrite(c, req.ExtendSystemPrompt, req.ExtendPrompt, []byte(originArtContext))
		if err != nil {
			return "", err
		}

		steps = append(steps, "extend rewrite done")

		return art, nil
	}

	var makeTitleF makeTitleF = func(originArtContext string) (string, error) {
		steps = append(steps, "make title")

		title, err := r.rewritemanager.CustomRewrite(c, req.MakeTitleSystemPrompt, req.MakeTitlePrompt, []byte(originArtContext))
		if err != nil {
			return "", err
		}

		steps = append(steps, "make title done")

		return title, nil
	}

	newArt, err := r.rewriteWorkFlow(req.Content, rewriteF, extendRewriteF, makeTitleF, true)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
			"steps":   steps,
		})

		return
	}

	newArt.Content, err = util.DecodeImageListDivFromHTMl([]byte(newArt.Content))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
			"steps":   steps,
		})
	}

	res := model.RewriteTestResponse{
		RewriteResponse: newArt,
		Steps:           steps,
	}

	c.JSON(http.StatusOK, res)
}

func (r *RewriteHandler) MultiSectionsRewriteHandler(c *gin.Context) {
	// get data body from request
	req, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	var rewriteF rewriteF = func(req string) (string, error) {
		art, err := r.rewritemanager.DefaultMultiSectionsRewrite(c, req)
		if err != nil {
			return "", err
		}

		return art, nil
	}

	var extendRewriteF extendRewriteF = func(req string) (string, error) {
		extArt, err := r.rewritemanager.DefaultExtendRewrite(c, []byte(req))
		if err != nil {
			return "", err
		}

		return extArt, nil
	}

	var makeTitleF makeTitleF = func(req string) (string, error) {
		title, err := r.rewritemanager.DefaultMakeTitle(c, req)
		if err != nil {
			return "", err
		}

		return title, nil
	}

	res, err := r.rewriteWorkFlow(string(req), rewriteF, extendRewriteF, makeTitleF, false)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	c.JSON(http.StatusOK, res)
}

func (r *RewriteHandler) MultiSectionsRewriteTestHandler(c *gin.Context) {
	req := model.MultiSectionsRewriteTestRequest{}

	err := c.ShouldBind(&req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	steps := []string{}

	var rewriteF rewriteF = func(originArtContext string) (string, error) {
		steps = append(steps, "rewrite")

		art, err := r.rewritemanager.DefaultMultiSectionsRewrite(c, originArtContext)
		if err != nil {
			return "", err
		}

		steps = append(steps, "rewrite done")

		return art, nil
	}

	var extendRewriteF extendRewriteF = func(originArtContext string) (string, error) {
		steps = append(steps, "extend rewrite")

		art, err := r.rewritemanager.DefaultExtendRewrite(c, []byte(originArtContext))
		if err != nil {
			return "", err
		}

		steps = append(steps, "extend rewrite done")

		return art, nil
	}

	var makeTitleF makeTitleF = func(originArtContext string) (string, error) {
		steps = append(steps, "make title")

		title, err := r.rewritemanager.CustomRewrite(c, req.MakeTitleSystemPrompt, req.MakeTitlePrompt, []byte(originArtContext))
		if err != nil {
			return "", err
		}

		steps = append(steps, "make title done")

		return title, nil
	}

	newArt, err := r.rewriteWorkFlow(req.Content, rewriteF, extendRewriteF, makeTitleF, false)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
			"steps":   steps,
		})

		return
	}

	res := model.RewriteTestResponse{
		RewriteResponse: newArt,
		Steps:           steps,
	}

	c.JSON(http.StatusOK, res)
}

type (
	rewriteF       func(string) (string, error)
	extendRewriteF func(string) (string, error)
	makeTitleF     func(string) (string, error)
)

func (r *RewriteHandler) rewriteWorkFlow(
	articleContent string,
	rewriteF func(string) (string, error),
	extendRewriteF func(string) (string, error),
	makeTitleF func(string) (string, error),
	isGenImageList bool,
) (model.RewriteResponse, error) {
	originalArticle, err := util.HTMLToMd(articleContent)
	if err != nil {
		log.Println(err)

		return model.RewriteResponse{}, fmt.Errorf("rewriteWorkFlow: %w", err)
	}

	// check data
	if utf8.RuneCount([]byte(originalArticle)) < minSrcLength {
		log.Println("data is not complete")

		return model.RewriteResponse{}, fmt.Errorf("rewriteWorkFlow: %w", errors.New("data is not complete"))
	}

	res := struct {
		Title   string `json:"Title"`
		Content string `json:"Content"`
	}{}

	art, err := rewriteF(originalArticle)
	if err != nil {
		return model.RewriteResponse{}, fmt.Errorf("rewriteWorkFlow: %w", err)
	}

	if utf8.RuneCountInString(art) < minArtLength {
		extArt, err := extendRewriteF(articleContent)
		if err != nil {
			return model.RewriteResponse{}, fmt.Errorf("rewriteWorkFlow: %w", err)
		}

		art = extArt
	}

	title, err := makeTitleF(art)
	if err != nil {
		return model.RewriteResponse{}, fmt.Errorf("rewriteWorkFlow: %w", err)
	}

	res.Title = title

	art = string(util.MdToHTML([]byte(art)))

	if isGenImageList {
		imgDiv, err := util.GenImageListEncodeDiv([]byte(articleContent))
		if err != nil {
			log.Printf("error: %v", err)
		} else {
			art += imgDiv
		}
	}

	res.Content = art

	return res, nil
}

func (r *RewriteHandler) CreateRewriteTestCaseHandler(c *gin.Context) {
	req := model.CreateRewriteTestCaseRequest{}

	err := c.ShouldBind(&req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	if req.Name == "" || req.Content == "" {
		log.Println("data is not complete")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", "data is not complete"),
		})

		return
	}

	testCase, err := r.rewritemanager.CreateRewriteTestCase(req.Name, req.Source, req.Content)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    testCase,
		"message": "ok",
	})
}

func (r *RewriteHandler) ListRewriteTestCaseHandler(c *gin.Context) {
	testCases, err := r.rewritemanager.ListRewriteTestCases()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"data":    testCases,
	})
}

func (r *RewriteHandler) UpdateRewriteTestCaseHandler(c *gin.Context) {
	req := model.UpdateRewriteTestCaseRequest{}

	err := c.ShouldBind(&req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})

		return
	}

	if req.Name == "" || req.Content == "" {
		log.Println("data is not complete")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", "data is not complete"),
		})

		return
	}

	ID := c.Param("id")

	err = r.rewritemanager.UpdateRewriteTestCase(ID, req.Name, req.Source, req.Content)
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

func (r *RewriteHandler) DeleteRewriteTestCaseHandler(c *gin.Context) {
	ID := c.Param("id")

	err := r.rewritemanager.DeleteRewriteTestCase(ID)
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

func (r *RewriteHandler) GetDefaultMultiSectionsSystemPromptHandler(c *gin.Context) {
	systemPrompt, err := r.rewritemanager.GetDefaultMultiSectionsSystemPrompt()
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

func (r *RewriteHandler) SetDefaultMultiSectionsSystemPromptHandler(c *gin.Context) {
	f := func(req model.SetDefaultMultiSectionsSystemPromptRequest) error {
		return r.rewritemanager.SetDefaultMultiSectionsSystemPrompt(req.Prompt)
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
