package model

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
		Token string `json:"token"`
	} `json:"data"`
}

type ListMemberResponse struct {
	BasicResponse
	Data Data[Member] `json:"data"`
}

type Member struct {
	ID     string `json:"ID"`
	Level  string `json:"Level"`
	Status string `json:"Status"`
	Name   string `json:"Name"`
	Email  string `json:"Email"`
}

type Article struct {
	ID         string `json:"ID"`
	CateID     string `json:"CateID"`
	AuthorID   string `json:"AuthorID"`
	Title      string `json:"Title"`
	Content    string `json:"Content"`
	Intro      string `json:"Intro"`
	PostTime   string `json:"PostTime"`
	UpdateTime string `json:"UpdateTime"`
	// IsTop    uint32    `json:"IsTop"`
	// IsLock   uint32    `json:"IsLock"`
}

type PostArticleResponse struct {
	BasicResponse
}

type ListArticleResponse struct {
	BasicResponse
	Data Data[Article] `json:"data"`
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
	Page   uint32 `json:"page"`
	CateID uint32 `json:"cate_id"`
	TagID  uint32 `json:"tag_id"`
}
