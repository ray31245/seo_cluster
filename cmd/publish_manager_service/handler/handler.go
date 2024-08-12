package handler

import (
	"fmt"
	"goTool/pkg/db"
	"goTool/pkg/db/model"
	zblogapi "goTool/pkg/z_blog_api"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	DAO *db.DAO
	// TODO: ZblogAPIPoll should be sync.Map
	// key is site id
	ZApiPool map[uuid.UUID]*zblogapi.ZblogAPI
}

func NewHandler(dao *db.DAO) *Handler {
	return &Handler{
		DAO:      dao,
		ZApiPool: make(map[uuid.UUID]*zblogapi.ZblogAPI),
	}
}

func (p *Handler) AveragePublishHandler(c *gin.Context) {
	// find first publiched category
	cate, err := p.DAO.FirstPublishedCategory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})
		return
	}

	// get data body from request
	article := zblogapi.PostArticleRequest{}
	c.ShouldBind(&article)

	// check data
	if article.Title == "" || article.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", "data is not complete"),
		})
		return
	}

	// set category id
	article.CateID = cate.ZblogID

	// get zblog api
	api, ok := p.ZApiPool[cate.SiteID]
	if !ok {
		api = zblogapi.NewZblogAPI(cate.Site.URL, cate.Site.UserName, cate.Site.Password)
		p.ZApiPool[cate.SiteID] = api
	}

	// post article
	err = api.PostArticle(article)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})
		return
	}

	// mark last published
	err = p.DAO.MarkPublished(cate.ID.String())
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

	// check site is valid
	api := zblogapi.NewZblogAPI(req.URL, req.UserName, req.Password)
	err := api.Login()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})
		return
	}

	// add site
	site, err := p.DAO.CreateSite(&model.Site{URL: req.URL, UserName: req.UserName, Password: req.Password})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})
		return
	}

	// list category of site
	categories, err := api.ListCategory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})
		return
	}

	// add category
	errs := make([]error, 0)
	for _, cate := range categories {
		cateID, err := strconv.Atoi(cate.ID)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		p.DAO.CreateCategory(&model.Category{
			SiteID:  site.ID,
			ZblogID: uint32(cateID),
		})
	}
	if len(errs) > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", errs),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
