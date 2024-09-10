package model

type CommentUser struct {
	Base
	Name         string `json:"name" gorm:"unique"`
	Alias        string `json:"alias"`
	Password     string `json:"password"`
	HashPassword string `json:"hash_password"`
}
