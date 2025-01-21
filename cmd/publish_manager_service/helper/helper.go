package helper

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	dbModel "github.com/ray31245/seo_cluster/pkg/db/model"
)

func ParsePagination(c *gin.Context) (page, pageSize int, err error) {
	pageQ, pageSizeQ := ParsePaginationQuery(c)

	page, pageSize, err = ParsePageAndPageSize(pageQ, pageSizeQ)
	if err != nil {
		return
	}

	return
}

func ParsePageAndPageSize(pageQ, pageSizeQ string) (int, int, error) {
	page, err := strconv.Atoi(pageQ)
	if err != nil {
		return 0, 0, errors.New("page must be a number")
	}

	pageSize, err := strconv.Atoi(pageSizeQ)
	if err != nil {
		return 0, 0, errors.New("page_size must be a number")
	}

	return page, pageSize, nil
}

func ParsePaginationQuery(c *gin.Context) (pageQ, pageSizeQ string) {
	pageQ = c.Query("page")
	pageSizeQ = c.Query("page_size")

	return
}

func ParseArticleCacheFuzzySearchQuery(c *gin.Context) (titleKeywords, contentKeywords string, operator dbModel.Operator) {
	titleKeywords = c.Query("title_keyword")
	contentKeywords = c.Query("content_keyword")
	operator = dbModel.Operator(c.Query("operator"))

	return
}
