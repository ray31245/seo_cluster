package authinterface

import (
	dbModel "github.com/ray31245/seo_cluster/pkg/db/model"
)

type AuthInterface interface {
	GenerateUser(userName, password string) (*dbModel.User, error)
}
