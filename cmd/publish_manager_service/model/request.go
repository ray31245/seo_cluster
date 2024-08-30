package model

import (
	zModel "github.com/ray31245/seo_cluster/pkg/z_blog_api/model"
)

type PublishArticleRequest struct {
	Title   string `json:"Title"`
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

type AddSiteRequest struct {
	URL      string `json:"url"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
}
