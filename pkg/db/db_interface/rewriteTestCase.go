package dbinterface

import (
	"github.com/ray31245/seo_cluster/pkg/db/model"
)

type RewriteTestCaseDAOInterface interface {
	CreateRewriteTestCase(rewriteTestCase *model.RewriteTestCase) (model.RewriteTestCase, error)
	GetRewriteTestCaseByID(id string) (*model.RewriteTestCase, error)
	ListRewriteTestCases() ([]model.RewriteTestCase, error)
	UpdateRewriteTestCase(id string, rewriteTestCase *model.RewriteTestCase) error
	DeleteRewriteTestCase(id string) error
}
