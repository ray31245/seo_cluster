package rewritemanager

import (
	"fmt"

	dbModel "github.com/ray31245/seo_cluster/pkg/db/model"
)

func (r *RewriteManager) CreateRewriteTestCase(name string, source string, content string) error {
	testCase := dbModel.RewriteTestCase{
		Name:    name,
		Source:  source,
		Content: content,
	}

	_, err := r.rewriteTestCase.CreateRewriteTestCase(&testCase)
	if err != nil {
		return fmt.Errorf("RewriteManager.CreateRewriteTestCase: %w", err)
	}

	return nil
}

func (r *RewriteManager) GetRewriteTestCaseByID(id string) (*dbModel.RewriteTestCase, error) {
	testCase, err := r.rewriteTestCase.GetRewriteTestCaseByID(id)
	if err != nil {
		return nil, fmt.Errorf("RewriteManager.GetRewriteTestCaseByID: %w", err)
	}

	return testCase, nil
}

func (r *RewriteManager) ListRewriteTestCases() ([]dbModel.RewriteTestCase, error) {
	testCases, err := r.rewriteTestCase.ListRewriteTestCases()
	if err != nil {
		return nil, fmt.Errorf("RewriteManager.ListRewriteTestCases: %w", err)
	}

	return testCases, nil
}

func (r *RewriteManager) UpdateRewriteTestCase(id string, name string, source string, content string) error {
	testCase := dbModel.RewriteTestCase{
		Name:    name,
		Source:  source,
		Content: content,
	}

	err := r.rewriteTestCase.UpdateRewriteTestCase(id, &testCase)
	if err != nil {
		return fmt.Errorf("RewriteManager.UpdateRewriteTestCase: %w", err)
	}

	return nil
}

func (r *RewriteManager) DeleteRewriteTestCase(id string) error {
	err := r.rewriteTestCase.DeleteRewriteTestCase(id)
	if err != nil {
		return fmt.Errorf("RewriteManager.DeleteRewriteTestCase: %w", err)
	}

	return nil
}
