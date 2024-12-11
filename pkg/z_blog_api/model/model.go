package model

import (
	"github.com/ray31245/seo_cluster/pkg/util"
)

// ----response----
type BasicResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Env       string `json:"env"`
	Zbp       string `json:"zbp"`
	AppCenter string `json:"appcenter"`
}

type LoginResponse struct {
	BasicResponse
	Data struct {
		Token      string        `json:"token"`
		ExpireTime util.UnixTime `json:"expire_time"`
	} `json:"data"`
}

type ListMemberResponse struct {
	BasicResponse
	Data Data[Member] `json:"data"`
}

type Member struct {
	ID     util.NumberString `json:"ID"`
	Level  string            `json:"Level,omitempty"`
	Status string            `json:"Status,omitempty"`
	Name   string            `json:"Name,omitempty"`
	Alias  string            `json:"Alias,omitempty"`
	Email  string            `json:"Email,omitempty"`
}

type Article struct {
	ID         util.NumberString `json:"ID"`
	CateID     util.NumberString `json:"CateID"`
	AuthorID   util.NumberString `json:"AuthorID"`
	Title      string            `json:"Title"`
	Content    string            `json:"Content"`
	CommNums   util.StringNumber `json:"CommNums"`
	Intro      string            `json:"Intro"`
	PostTime   util.UnixTime     `json:"PostTime"`
	UpdateTime util.UnixTime     `json:"UpdateTime"`
	// IsTop    uint32    `json:"IsTop"`
	// IsLock   uint32    `json:"IsLock"`
}

type PostArticleResponse struct {
	BasicResponse
	Data struct {
		Post Article `json:"post"`
	} `json:"data"`
}

type GetArticleResponse struct {
	BasicResponse
	Data struct {
		Post Article `json:"post"`
	} `json:"data"`
}

type ListArticleResponse struct {
	BasicResponse
	Data Data[Article] `json:"data"`
}

type PostCommentResponse struct {
	BasicResponse
}

type PageBar struct {
	AllCount     uint32 `json:"AllCount"`
	CurrentCount uint32 `json:"CurrentCount"`
	PerPageCount uint32 `json:"PerPageCount"`
	PageAll      uint32 `json:"PageAll"`
	PageNow      uint32 `json:"PageNow"`
	PageCurrent  uint32 `json:"PageCurrent"`
	PageFirst    uint32 `json:"PageFirst"`
	PageLast     uint32 `json:"PageLast"`
	PageNext     uint32 `json:"PageNext"`
	PagePrevious uint32 `json:"PagePrevious"`
}

type DeleteArticleResponse struct {
	BasicResponse
}

type Category struct {
	ID    string `json:"ID"`
	Name  string `json:"Name"`
	Count string `json:"Count"`
}

type ListCategoryResponse struct {
	BasicResponse
	Data Data[Category] `json:"data"`
}

type Data[E any] struct {
	List    []E     `json:"list"`
	PageBar PageBar `json:"pagebar"`
}

type PostMemberResponse struct {
	BasicResponse
	Data struct {
		Member Member `json:"member"`
	} `json:"data"`
}

type Tag struct {
	ID         util.NumberString `json:"ID"`
	Name       string            `json:"Name"`
	Count      string            `json:"Count"`
	CreateTime util.UnixTime     `json:"CreateTime"`
	UpdateTime util.UnixTime     `json:"UpdateTime"`
}

type ListTagResponse struct {
	BasicResponse
	Data Data[Tag] `json:"data"`
}

type PostTagResponse struct {
	BasicResponse
	Data struct {
		Tag Tag `json:"tag"`
	} `json:"data"`
}

// ----request----

type PageRequest struct {
	Page    uint32 `json:"page"`
	Perpage uint32 `json:"perpage"`
	SortBy  string `json:"sortby"`
	Order   string `json:"order"`
}

type PostArticleRequest struct {
	ID      uint32 `json:"ID"`
	Title   string `json:"Title,omitempty"`
	Content string `json:"Content,omitempty"`
	Intro   string `json:"Intro,omitempty"`
	CateID  uint32 `json:"CateID,omitempty"`
	Tag     string `json:"Tag,omitempty"`
	Type    uint32 `json:"Type,omitempty"`
}

type ListArticleRequest struct {
	CateID uint32 `json:"cate_id"`
	TagID  uint32 `json:"tag_id"`
	PageRequest
}

type PostCommentRequest struct {
	LogID   string `json:"LogID"`
	Content string `json:"Content"`
}

type PostMemberRequest struct {
	Member
	Password   string `json:"Password,omitempty"`
	PasswordRe string `json:"PasswordRe,omitempty"`
}

type ListTagRequest struct {
	PageRequest
}

type PostTagRequest struct {
	ID   uint32 `json:"ID"`
	Name string `json:"Name"`
}
