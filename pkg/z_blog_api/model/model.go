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
	Level  string            `json:"Level"`
	Status string            `json:"Status"`
	Name   string            `json:"Name"`
	Alias  string            `json:"Alias"`
	Email  string            `json:"Email"`
}

type Article struct {
	ID         string            `json:"ID"`
	CateID     string            `json:"CateID"`
	AuthorID   string            `json:"AuthorID"`
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

// ----request----

type PostArticleRequest struct {
	ID      uint32 `json:"ID"`
	Title   string `json:"Title"`
	Content string `json:"Content"`
	Intro   string `json:"Intro"`
	CateID  uint32 `json:"CateID"`
	Type    uint32 `json:"Type"`
}

type ListArticleRequest struct {
	Page    uint32 `json:"page"`
	CateID  uint32 `json:"cate_id"`
	TagID   uint32 `json:"tag_id"`
	Perpage uint32 `json:"perpage"`
	Sortby  string `json:"sortby"`
	Order   string `json:"order"`
}

type PostCommentRequest struct {
	LogID   string `json:"LogID"`
	Content string `json:"Content"`
}

type PostMemberRequest struct {
	Member
	Password   string `json:"Password"`
	PasswordRe string `json:"PasswordRe"`
}
