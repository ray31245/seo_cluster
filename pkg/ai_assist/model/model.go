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
