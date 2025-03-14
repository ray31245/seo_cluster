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

type SetPromptRequest struct {
	Prompt string `json:"prompt"`
}

func (s SetPromptRequest) GetPrompt() string {
	return s.Prompt
}

type SetDefaultSystemPromptRequest struct {
	SetPromptRequest
}

type SetDefaultPromptRequest struct {
	SetPromptRequest
}

type SetDefaultExtendSystemPromptRequest struct {
	SetPromptRequest
}

type SetDefaultExtendPromptRequest struct {
	SetPromptRequest
}

type SetDefaultMakeTitleSystemPromptRequest struct {
	SetPromptRequest
}

type SetDefaultMakeTitlePromptRequest struct {
	SetPromptRequest
}

type SetDefaultMultiSectionsSystemPromptRequest struct {
	SetPromptRequest
}

type CreateRewriteTestCaseRequest struct {
	Name    string `json:"name"`
	Source  string `json:"source"`
	Content string `json:"content"`
}

type UpdateRewriteTestCaseRequest struct {
	Name    string `json:"name"`
	Source  string `json:"source"`
	Content string `json:"content"`
}

type RewriteTestRequest struct {
	Content               string `json:"content"`
	SystemPrompt          string `json:"system_prompt"`
	Prompt                string `json:"prompt"`
	ExtendSystemPrompt    string `json:"extend_system_prompt"`
	ExtendPrompt          string `json:"extend_prompt"`
	MakeTitleSystemPrompt string `json:"make_title_system_prompt"`
	MakeTitlePrompt       string `json:"make_title_prompt"`
}

type MultiSectionsRewriteTestRequest struct {
	Content               string `json:"content"`
	SystemPrompt          string `json:"system_prompt"`
	MakeTitleSystemPrompt string `json:"make_title_system_prompt"`
	MakeTitlePrompt       string `json:"make_title_prompt"`
}
