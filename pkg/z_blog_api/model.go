package zblogapi

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
}

type PostArticleResponse struct {
	BasicResponse
}

// ----request----

type ArticleRequest struct {
	ID      int    `json:"ID"`
	Title   string `json:"Title"`
	Content string `json:"Content"`
	Type    int    `json:"Type"`
}
