package dbinterface

import (
	"github.com/ray31245/seo_cluster/pkg/db/model"
)

type UserDAOInterface interface {
	CreateFirstAdminUser(user *model.User) error
}
