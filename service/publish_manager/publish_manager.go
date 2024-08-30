package publishmanager

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	dbInterface "github.com/ray31245/seo_cluster/pkg/db/db_interface"
	dbErr "github.com/ray31245/seo_cluster/pkg/db/error"
	dbModel "github.com/ray31245/seo_cluster/pkg/db/model"
	zModel "github.com/ray31245/seo_cluster/pkg/z_blog_api/model"
	zInterface "github.com/ray31245/seo_cluster/pkg/z_blog_api/z_blog_Interface"
)

var ErrNoCategoryNeedToBePublished = errors.New("no category need to be published")

type DAO struct {
	dbInterface.ArticleCacheDAOInterface
	dbInterface.SiteDAOInterface
}

type PublishManager struct {
	zAPI zInterface.ZBlogAPI
	dao  DAO
}

func NewPublishManager(zAPI zInterface.ZBlogAPI, dao DAO) *PublishManager {
	return &PublishManager{
		zAPI: zAPI,
		dao:  dao,
	}
}

// AddSite add site to publish manager
func (p PublishManager) AddSite(ctx context.Context, urlStr string, userName string, password string) error {
	// check site is valid
	client, err := p.zAPI.NewClient(ctx, urlStr, userName, password)
	if err != nil {
		return fmt.Errorf("AddSite: %w", err)
	}

	// add site
	site, err := p.dao.CreateSite(&dbModel.Site{URL: urlStr, UserName: userName, Password: password})
	if err != nil {
		return fmt.Errorf("AddSite: %w", err)
	}

	// list category of site
	categories, err := client.ListCategory(ctx)
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

// List sites of publish manager
func (p PublishManager) ListSites() ([]dbModel.Site, error) {
	res, err := p.dao.ListSites()
	if err != nil {
		return nil, fmt.Errorf("ListSites: %w", err)
	}

	return res, nil
}

// AveragePublish average publish article to all site and category
func (p PublishManager) AveragePublish(ctx context.Context, article zModel.PostArticleRequest) error {
	// find first published category
	cate, err := p.dao.FirstPublishedCategory()
	if err != nil {
		if dbErr.IsNotfoundErr(err) {
			err = errors.Join(ErrNoCategoryNeedToBePublished, err)
		}

		return fmt.Errorf("AveragePublish: %w", err)
	}

	// set category id
	article.CateID = cate.ZBlogID

	// get zblog api client
	client, err := p.zAPI.GetClient(ctx, cate.SiteID, cate.Site.URL, cate.Site.UserName, cate.Site.Password)
	if err != nil {
		return fmt.Errorf("AveragePublish: %w", err)
	}

	// post article
	err = client.PostArticle(ctx, article)
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

func (p PublishManager) StartRandomCyclePublish(ctx context.Context) error {
	lastCategory, err := p.dao.LastPublishedCategory()
	if err == nil {
		if time.Since(lastCategory.LastPublished).Minutes() > maxCycleTime {
			err = p.cyclePublish(ctx)
			if err != nil {
				return fmt.Errorf("StartRandomCyclePublish: %w", err)
			}
		}
	} else if !dbErr.IsNotfoundErr(err) {
		return fmt.Errorf("StartRandomCyclePublish: %w", err)
	}

	go func() {
		for {
			time.Sleep(randomTime())

			err = p.cyclePublish(ctx)
			if err != nil {
				log.Println(err)
			}
		}
	}()

	return nil
}

func (p PublishManager) cyclePublish(ctx context.Context) error {
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
		err := p.AveragePublish(ctx, zModel.PostArticleRequest{Title: article.Title, Content: article.Content})
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
