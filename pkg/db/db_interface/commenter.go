package dbinterface

import (
	"github.com/ray31245/seo_cluster/pkg/db/model"
)

type CommentUserDAOInterface interface {
	GetRandomCommentUser() (model.CommentUser, error)
}
