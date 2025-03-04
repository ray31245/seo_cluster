package model

type RewriteTestCase struct {
	Base
	Name    string `json:"name" gorm:"unique"`
	Source  string `json:"source"`
	Content string `json:"content"`
}
