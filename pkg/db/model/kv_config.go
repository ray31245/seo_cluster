package model

type KVConfig struct {
	Base
	Key   string `json:"key" gorm:"unique,index"`
	Value string `json:"value"`
}
