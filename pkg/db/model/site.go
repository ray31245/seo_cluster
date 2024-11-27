package model

type CMSType string

const (
	CMSTypeWordPress CMSType = "wordpress"
	CMSTypeZBlog     CMSType = "zblog"
)

type Site struct {
	Base
	URL        string `json:"url" gorm:"unique"`
	UserName   string `json:"username"`
	Password   string `json:"password"`
	LackCount  int    `json:"lack_count" gorm:"default:0"`
	Categories []Category
	CmsType    CMSType `json:"cms_type"`
}
