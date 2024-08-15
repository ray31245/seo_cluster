package model

type Site struct {
	Base
	URL        string `json:"url" gorm:"unique"`
	UserName   string `json:"username"`
	Password   string `json:"password"`
	LackCount  int    `json:"lack_count" gorm:"default:0"`
	Categories []Category
}
