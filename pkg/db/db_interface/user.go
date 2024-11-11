package dbinterface

import (
	"github.com/ray31245/seo_cluster/pkg/db/model"
)

type UserDAOInterface interface {
	CreateUser(user *model.User) error
	Count() (int64, error)
	GetUserByID(id string) (*model.User, error)
	GetUserByName(name string) (*model.User, error)
}
