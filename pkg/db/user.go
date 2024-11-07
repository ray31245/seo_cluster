package db

import (
	"fmt"

	dbErr "github.com/ray31245/seo_cluster/pkg/db/error"
	"github.com/ray31245/seo_cluster/pkg/db/model"

	"gorm.io/gorm"
)

type UserDAO struct {
	db *gorm.DB
}

func (d *DB) NewUserDAO() (*UserDAO, error) {
	err := d.db.AutoMigrate(&model.User{})
	if err != nil {
		return nil, fmt.Errorf("NewUserDAO: %w", err)
	}

	return &UserDAO{db: d.db}, nil
}

func (d *UserDAO) CreateFirstAdminUser(user *model.User) error {
	var count int64

	err := d.db.Model(&model.User{}).Count(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		return dbErr.ErrUserAlreadyExist
	}

	user.IsAdmin = true

	err = d.db.Create(user).Error
	if err != nil {
		return err
	}

	return nil
}
