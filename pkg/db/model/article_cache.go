package model

type ArticleCache struct {
	Base
	Title   string `json:"title"`
	Content string `json:"content"`
}
