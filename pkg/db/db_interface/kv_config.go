package dbinterface

import (
	"github.com/ray31245/seo_cluster/pkg/db/model"
)

type KVConfigDAOInterface interface {
	UpsertByKey(key string, value string) error
	UpsertByKeyInt(key string, value int) error
	UpsertByKeyBool(key string, value bool) error
	GetByKey(key string) (model.KVConfig, error)
	GetByKeyWithDefault(key string, defaultValue string) (model.KVConfig, error)
	GetIntByKeyWithDefault(key string, defaultValue int) (int, error)
	GetBoolByKeyWithDefault(key string, defaultValue bool) (bool, error)
	GetByKeyString(key string) (string, error)
	GetByKeyInt(key string) (int, error)
	GetByKeyBool(key string) (bool, error)
}
