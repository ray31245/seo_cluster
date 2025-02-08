package model

import (
	zModel "github.com/ray31245/seo_cluster/pkg/z_blog_api/model"
	publishModel "github.com/ray31245/seo_cluster/service/publish_manager/model"
)

type PublishArticleRequest struct {
	Title   string `json:"Title"`
	IsTop   bool   `json:"IsTop"`
	Content string `json:"Content"`
	Intro   string `json:"Intro"`
	CateID  uint32 `json:"CateID"`
}

func (p *PublishArticleRequest) ToZBlogAPI() zModel.PostArticleRequest {
	return zModel.PostArticleRequest{
		Title:   p.Title,
		Content: p.Content,
		Intro:   p.Intro,
		CateID:  p.CateID,
	}
}

func (p *PublishArticleRequest) ToPublishManager() publishModel.Article {
	return publishModel.Article{
		Title:   p.Title,
		IsTop:   p.IsTop,
		Content: p.Content,
		CateID:  p.CateID,
	}
}

type AddSiteRequest struct {
	URL               string `json:"url"`
	CMSType           string `json:"cms_type"`
	UserName          string `json:"user_name"`
	Password          string `json:"password"`
	ExpectCategoryNum uint8  `json:"expect_category_num"`
}

type UpdateSiteRequest struct {
	SiteID   string `json:"site_id"`
	URL      string `json:"url"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type IncreaseLackCountRequest struct {
	SiteID string `json:"site_id"`
	Count  int    `json:"count"`
}

type AddFirstAdminUserRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type SetUnconfigCategoryNameRequest struct {
	Name string `json:"name"`
}

type SetConfigTagBlackListRequest struct {
	Tags []string `json:"tags"`
}

type UpdateArticleCacheStatusRequest struct {
	IDs    []string `json:"ids"`
	Status string   `json:"status"`
}

type EditArticleCacheRequest struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type DeleteArticleCacheRequest struct {
	IDs []string `json:"ids"`
}
