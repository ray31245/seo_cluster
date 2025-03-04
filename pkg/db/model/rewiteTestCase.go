package model

type RewriteTestCase struct {
	Base
	Name    string `json:"name"`
	Source  string `json:"source"`
	Content string `json:"content"`
}
