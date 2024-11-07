package model

type User struct {
	Base
	Name         string `json:"name"`
	HashPassword string `json:"hash_password"`
	IsAdmin      bool   `json:"is_admin"`
}
