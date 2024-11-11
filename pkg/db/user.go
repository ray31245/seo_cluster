package db

import (
	"fmt"

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

func (d *UserDAO) CreateUser(user *model.User) error {
	err := d.db.Create(user).Error
	if err != nil {
		return err
	}

	return nil
}

func (d *UserDAO) Count() (int64, error) {
	var count int64

	err := d.db.Model(&model.User{}).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (d *UserDAO) GetUserByID(id string) (*model.User, error) {
	var user model.User

	err := d.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (d *UserDAO) GetUserByName(name string) (*model.User, error) {
	var user model.User

	err := d.db.First(&user, "name = ?", name).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
