package dbinterface

import (
	"github.com/ray31245/seo_cluster/pkg/db/model"
)

type KVConfigDAOInterface interface {
	UpsertByKey(key string, value string) error
	GetByKey(key string) (model.KVConfig, error)
	GetByKeyWithDefault(key string, defaultValue string) (model.KVConfig, error)
}
