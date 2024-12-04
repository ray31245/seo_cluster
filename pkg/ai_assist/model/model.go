package model

type RewriteResponse struct {
	Title   string `json:"Title"`
	Content string `json:"Content"`
}

type CommentResponse struct {
	Comment string `json:"Comment"`
	Score   int    `json:"Score"`
}

type EvaluateResponse struct {
	Score int `json:"Score"`
}

type FindKeyWordsResponse struct {
	KeyWords []string `json:"KeyWords"`
}

type CategoryOption struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SelectCategoryRequest struct {
	Text             []byte
	CategoriesOption []CategoryOption
}

type SelectCategoryResponse struct {
	ID     string `json:"id"`
	IsFind bool   `json:"isFind"`
}
