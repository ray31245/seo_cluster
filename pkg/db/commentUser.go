package db

import (
	"fmt"

	"github.com/ray31245/seo_cluster/pkg/db/model"
	"gorm.io/gorm"
)

type CommentUserDAO struct {
	db *gorm.DB
}

func (d *DB) NewCommentUserDAO() (*CommentUserDAO, error) {
	err := d.db.AutoMigrate(&model.CommentUser{})
	if err != nil {
		return nil, fmt.Errorf("NewUserDAO: %w", err)
	}

	return &CommentUserDAO{db: d.db}, nil
}

func (d *CommentUserDAO) CreateCommentUser(user *model.CommentUser) (model.CommentUser, error) {
	err := d.db.Create(user).Error
	if err != nil {
		return model.CommentUser{}, err
	}

	res := model.CommentUser{}

	err = d.db.Where("name = ?", user.Name).First(&res).Error
	if err != nil {
		return model.CommentUser{}, err
	}

	return res, nil
}

func (d *CommentUserDAO) GetRandomCommentUser() (model.CommentUser, error) {
	var user model.CommentUser
	err := d.db.Order("RANDOM()").First(&user).Error

	return user, err
}

func (d *CommentUserDAO) ListCommentUsers() ([]model.CommentUser, error) {
	users := []model.CommentUser{}
	err := d.db.Find(&users).Error

	return users, err
}
