package publishmanager

import (
	"errors"
	"fmt"
	dbInterface "goTool/pkg/db/db_interface"
	dbModel "goTool/pkg/db/model"
	zModel "goTool/pkg/z_blog_api/model"
	zInterface "goTool/pkg/z_blog_api/zblog_Interface"
	"strconv"
)

type PublishManager struct {
	zApi zInterface.ZBlogApi
	dao  dbInterface.DAOInterface
}

func NewPublishManager(zApi zInterface.ZBlogApi, dao dbInterface.DAOInterface) *PublishManager {
	return &PublishManager{
		zApi: zApi,
		dao:  dao,
	}
}

// AddSite add site to publish manager
func (p PublishManager) AddSite(urlStr string, userName string, password string) error {
	// check site is valid
	client, err := p.zApi.NewClient(urlStr, userName, password)
	if err != nil {
		return fmt.Errorf("AddSite: %w", err)
	}

	// add site
	site, err := p.dao.CreateSite(&dbModel.Site{URL: urlStr, UserName: userName, Password: password})
	if err != nil {
		return fmt.Errorf("AddSite: %w", err)
	}

	// list category of site
	categories, err := client.ListCategory()
	if err != nil {
		return fmt.Errorf("AddSite: %w", err)
	}

	// add category
	var multiErr error
	for _, cate := range categories {
		cateID, err := strconv.Atoi(cate.ID)
		if err != nil {
			multiErr = errors.Join(multiErr, err)
			continue
		}
		err = p.dao.CreateCategory(&dbModel.Category{SiteID: site.ID, ZblogID: uint32(cateID)})
		if err != nil {
			multiErr = errors.Join(multiErr, err)
		}
	}

	if multiErr != nil {
		return fmt.Errorf("AddSite: %w", multiErr)
	}

	return nil
}

// AveragePublish average publish article to all site
func (p PublishManager) AveragePublish(article zModel.PostArticleRequest) error {
	// find first publiched category
	cate, err := p.dao.FirstPublishedCategory()
	if err != nil {
		return fmt.Errorf("AveragePublish: %w", err)
	}

	// set category id
	article.CateID = cate.ZblogID

	// get zblog api client
	client, err := p.zApi.GetClient(cate.SiteID, cate.Site.URL, cate.Site.UserName, cate.Site.Password)
	if err != nil {
		return fmt.Errorf("AveragePublish: %w", err)
	}

	// post article
	err = client.PostArticle(article)
	if err != nil {
		return fmt.Errorf("AveragePublish: %w", err)
	}

	// mark last published
	err = p.dao.MarkPublished(cate.ID.String())
	if err != nil {
		return fmt.Errorf("AveragePublish: %w", err)
	}

	return nil
}