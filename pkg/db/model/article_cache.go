package model

type ArticleCacheStatus string

const (
	ArticleCacheStatusDefault  ArticleCacheStatus = "default"
	ArticleCacheStatusReserved ArticleCacheStatus = "reserved"
	ArticleCacheStatusInBuffer ArticleCacheStatus = "in_buffer"
)

var ArticleCacheStatues = []ArticleCacheStatus{ArticleCacheStatusDefault, ArticleCacheStatusReserved, ArticleCacheStatusInBuffer}

type ArticleCache struct {
	Base
	Title   string             `json:"title"`
	Content string             `json:"content"`
	Status  ArticleCacheStatus `json:"status" gorm:"default:defualt"`
}
