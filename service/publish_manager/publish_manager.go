package publishManager

import (
	"errors"
	"fmt"
	dbInterface "goTool/pkg/db/db_interface"
	dbErr "goTool/pkg/db/error"
	dbModel "goTool/pkg/db/model"
	zModel "goTool/pkg/z_blog_api/model"
	zInterface "goTool/pkg/z_blog_api/z_blog_Interface"
	"log"
	"strconv"
	"time"
)

var (
	ErrNoCategoryNeedToBePublished = dbErr.ErrNoCategoryNeedToBePublished
)

type DAO struct {
	dbInterface.ArticleCacheDAOInterface
	dbInterface.SiteDAOInterface
}

type PublishManager struct {
	zApi zInterface.ZBlogApi
	dao  DAO
}

func NewPublishManager(zApi zInterface.ZBlogApi, dao DAO) *PublishManager {
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
		err = p.dao.CreateCategory(&dbModel.Category{SiteID: site.ID, ZBlogID: uint32(cateID)})
		if err != nil {
			multiErr = errors.Join(multiErr, err)
		}
	}

	if multiErr != nil {
		return fmt.Errorf("AddSite: %w", multiErr)
	}

	return nil
}

// AveragePublish average publish article to all site and category
func (p PublishManager) AveragePublish(article zModel.PostArticleRequest) error {
	// find first published category
	cate, err := p.dao.FirstPublishedCategory()
	if err != nil {
		if errors.Is(err, dbErr.ErrNoCategoryNeedToBePublished) {
			err = ErrNoCategoryNeedToBePublished
		}
		return fmt.Errorf("AveragePublish: %w", err)
	}

	// set category id
	article.CateID = cate.ZBlogID

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

func (p PublishManager) PrePublish(article zModel.PostArticleRequest) error {
	cache := dbModel.ArticleCache{
		Title:   article.Title,
		Content: article.Content,
	}
	err := p.dao.AddArticleToCache(cache)
	if err != nil {
		return fmt.Errorf("PrePublish: %w", err)
	}
	return nil
}

func (p PublishManager) StartRandomCyclePublish() error {
	lastCategory, err := p.dao.LastPublishedCategory()
	if err != nil {
		return fmt.Errorf("StartRandomCyclePublish: %w", err)
	}
	if time.Since(lastCategory.LastPublished).Minutes() > maxCycleTime {
		err = p.cyclePublish()
		if err != nil {
			return fmt.Errorf("StartRandomCyclePublish: %w", err)
		}
	}
	go func() {
		for {
			time.Sleep(randomTime())
			err = p.cyclePublish()
			if err != nil {
				log.Println(err)
			}
		}

	}()
	return nil
}

func (p PublishManager) cyclePublish() error {
	sites, err := p.dao.ListSites()
	if err != nil {
		return fmt.Errorf("cyclePublish: %w", err)
	}
	totalLackCount := 0
	for _, site := range sites {
		if site.LackCount != 0 {
			continue
		}
		lackCount := randomNum()
		if lackCount > 0 {
			err := p.dao.IncreaseLackCount(site.ID.String(), int(lackCount))
			if err != nil {
				return fmt.Errorf("cyclePublish: %w", err)
			}
			totalLackCount += int(lackCount)
		}
	}
	articles, err := p.dao.ListArticleCacheByLimit(totalLackCount)
	if err != nil {
		return fmt.Errorf("cyclePublish: %w", err)
	}
	for _, article := range articles {
		err := p.AveragePublish(zModel.PostArticleRequest{Title: article.Title, Content: article.Content})
		if err != nil {
			return fmt.Errorf("cyclePublish: %w", err)
		}
		err = p.dao.DeleteArticleCache(article.ID.String())
		if err != nil {
			return fmt.Errorf("cyclePublish: %w", err)
		}
	}
	return nil
}
