package model

type Article struct {
	Title   string `json:"Title"`
	Content string `json:"Content"`
	CateID  uint32 `json:"CateID"`
}
