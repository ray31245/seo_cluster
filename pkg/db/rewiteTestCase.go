package db

import (
	"fmt"

	"github.com/ray31245/seo_cluster/pkg/db/model"

	"gorm.io/gorm"
)

type RewriteTestCaseDAO struct {
	db *gorm.DB
}

func (d *DB) NewRewriteTestCaseDAO() (*RewriteTestCaseDAO, error) {
	err := d.db.AutoMigrate(&model.RewriteTestCase{})
	if err != nil {
		return nil, fmt.Errorf("NewRewriteTestCaseDAO: %w", err)
	}

	return &RewriteTestCaseDAO{db: d.db}, nil
}

func (d *RewriteTestCaseDAO) CreateRewriteTestCase(rewriteTestCase *model.RewriteTestCase) (model.RewriteTestCase, error) {
	err := d.db.Create(rewriteTestCase).Error
	if err != nil {
		return model.RewriteTestCase{}, err
	}

	res := model.RewriteTestCase{}

	err = d.db.Where("id = ?", rewriteTestCase.ID).First(&res).Error
	if err != nil {
		return model.RewriteTestCase{}, fmt.Errorf("CreateRewriteTestCase: %w", err)
	}

	return res, nil
}

func (d *RewriteTestCaseDAO) GetRewriteTestCaseByID(id string) (*model.RewriteTestCase, error) {
	var rewriteTestCase model.RewriteTestCase

	err := d.db.First(&rewriteTestCase, "id = ?", id).Error
	if err != nil {
		return nil, fmt.Errorf("GetRewriteTestCaseByID: %w", err)
	}

	return &rewriteTestCase, nil
}

func (d *RewriteTestCaseDAO) ListRewriteTestCases() ([]model.RewriteTestCase, error) {
	var rewriteTestCases []model.RewriteTestCase
	err := d.db.Find(&rewriteTestCases).Error

	return rewriteTestCases, err
}

func (d *RewriteTestCaseDAO) UpdateRewriteTestCase(rewriteTestCase *model.RewriteTestCase) error {
	err := d.db.Save(rewriteTestCase).Error
	if err != nil {
		return fmt.Errorf("UpdateRewriteTestCase: %w", err)
	}

	return nil
}

func (d *RewriteTestCaseDAO) DeleteRewriteTestCase(id string) error {
	err := d.db.Delete(&model.RewriteTestCase{}, "id = ?", id).Error
	if err != nil {
		return fmt.Errorf("DeleteRewriteTestCase: %w", err)
	}

	return nil
}
