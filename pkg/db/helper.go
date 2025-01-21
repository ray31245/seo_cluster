package db

import (
	"fmt"

	"github.com/ray31245/seo_cluster/pkg/db/model"
	"gorm.io/gorm"
)

func paginator[T any](tableModel T, query *gorm.DB, page int, limit int) ([]T, int, error) {
	var (
		result    []T
		totalPage int
		count     int64
	)

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	err := query.Model(tableModel).Count(&count).Error
	if err != nil {
		return nil, 0, fmt.Errorf("paginator: %w", err)
	}

	totalPage = int(count) / limit
	if int(count)%limit != 0 {
		totalPage++
	}

	err = query.Model(tableModel).Offset((page - 1) * limit).Limit(limit).Find(&result).Error
	if err != nil {
		return nil, 0, fmt.Errorf("paginator: %w", err)
	}

	return result, totalPage, nil
}

func articleFuzzySearch(query *gorm.DB, titleKeyword, contentKeyword string, op model.Operator) *gorm.DB {
	if op == "" {
		op = "OR"
	}

	switch {
	case titleKeyword != "" && contentKeyword != "" && op == "AND":
		query = query.Where("title LIKE ? AND content LIKE ?", "%"+titleKeyword+"%", "%"+contentKeyword+"%")
	case titleKeyword != "" && contentKeyword != "" && op == "OR":
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+titleKeyword+"%", "%"+contentKeyword+"%")
	case titleKeyword != "":
		query = query.Where("title LIKE ?", "%"+titleKeyword+"%")
	case contentKeyword != "":
		query = query.Where("content LIKE ?", "%"+contentKeyword+"%")
	}

	return query
}
