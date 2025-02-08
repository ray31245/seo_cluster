package model

import (
	"time"

	"github.com/ray31245/seo_cluster/pkg/util"
	wordpressModel "github.com/ray31245/seo_cluster/pkg/wordpress_api/model"
	zModel "github.com/ray31245/seo_cluster/pkg/z_blog_api/model"
)

type Article struct {
	Title   string `json:"Title"`
	IsTop   bool   `json:"IsTop"`
	Content string `json:"Content"`
	CateID  uint32 `json:"CateID"`
}

func (a *Article) ToZBlogCreateRequest() zModel.PostArticleRequest {
	var isTop uint8 = 0
	if a.IsTop {
		isTop = 1
	}

	return zModel.PostArticleRequest{
		Title:    a.Title,
		IsTop:    isTop,
		Content:  a.Content,
		CateID:   a.CateID,
		Intro:    a.Content,
		PostTime: &util.UnixTime{Time: time.Now()},
	}
}

func (a *Article) ToWordpressCreateArgs(status wordpressModel.ArticleStatus) wordpressModel.CreateArticleArgs {
	date := wordpressModel.Date{
		Time: time.Now(),
	}

	return wordpressModel.CreateArticleArgs{
		Title:      a.Title,
		Sticky:     a.IsTop,
		Content:    a.Content,
		Categories: []uint32{a.CateID},
		Status:     status,
		Date:       &date,
	}
}
