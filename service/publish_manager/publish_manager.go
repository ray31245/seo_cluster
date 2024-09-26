package publishmanager

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	dbInterface "github.com/ray31245/seo_cluster/pkg/db/db_interface"
	dbErr "github.com/ray31245/seo_cluster/pkg/db/error"
	dbModel "github.com/ray31245/seo_cluster/pkg/db/model"
	zModel "github.com/ray31245/seo_cluster/pkg/z_blog_api/model"
	zInterface "github.com/ray31245/seo_cluster/pkg/z_blog_api/z_blog_Interface"
)

var ErrNoCategoryNeedToBePublished = errors.New("no category need to be published")

type PublishErr struct {
	SiteID uuid.UUID
	CateID uuid.UUID
}

func (s PublishErr) Error() string {
	return fmt.Sprintf("site id %s, cate id %s", s.SiteID, s.CateID)
}

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

	log.Printf("category id %s, site id %s in AveragePublish", cate.ID, cate.SiteID)

	// set category id
	article.CateID = cate.ZBlogID

	// do publish
	err = p.doPublish(ctx, article, cate.Site)
	if err != nil {
		return errors.Join(PublishErr{SiteID: cate.SiteID, CateID: cate.ID}, err)
	}

	// mark last published
	err = p.dao.MarkPublished(cate.ID.String())
	if err != nil {
		return fmt.Errorf("AveragePublish: %w", err)
	}

	return nil
}

func (p PublishManager) doPublish(ctx context.Context, article zModel.PostArticleRequest, site dbModel.Site) error {
	// get zblog api client
	client, err := p.zAPI.GetClient(ctx, site.ID, site.URL, site.UserName, site.Password)
	if err != nil {
		return fmt.Errorf("doPublish: %w", err)
	}

	// post article
	err = client.PostArticle(ctx, article)
	if err != nil {
		return fmt.Errorf("doPublish: %w", err)
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
		log.Printf("last Publish time %s in StartRandomCyclePublish", lastCategory.LastPublished)
		log.Printf("time now %s in StartRandomCyclePublish", time.Now())

		if time.Since(lastCategory.LastPublished).Minutes() > maxCycleTime {
			log.Println("Duration is more than maxCycleTime in StartRandomCyclePublish, cyclePublish forced to run")

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
			nextTime := randomTime()
			log.Printf("next time run cyclePublish is %s in StartRandomCyclePublish", time.Now().Add(nextTime))
			select {
			case <-ctx.Done():
				// Exit the loop if the context is cancelled
				return
			case <-time.After(nextTime):
				// Proceed with the publishing cycle after a random duration
				if err := p.cyclePublish(ctx); err != nil {
					log.Println("Error during cyclePublish:", err)
				}
			}
		}
	}()

	return nil
}

func (p PublishManager) cyclePublish(ctx context.Context) error {
	log.Println("cyclePublish running...")

	sites, err := p.dao.ListSites()
	if err != nil {
		return fmt.Errorf("cyclePublish: %w", err)
	}

	for _, site := range sites {
		if site.LackCount != 0 {
			continue
		}

		lackCount := randomNum()
		if lackCount > 0 {
			log.Printf("site id %s, lack count %d in cyclePublish", site.ID, lackCount)

			err := p.dao.IncreaseLackCount(site.ID.String(), int(lackCount))
			if err != nil {
				return fmt.Errorf("cyclePublish: %w", err)
			}
		}
	}

	// get total lack count
	totalLackCount, err := p.dao.SumLackCount()
	if err != nil {
		return fmt.Errorf("cyclePublish: %w", err)
	}

	articles, err := p.dao.ListArticleCacheByLimit(totalLackCount)
	if err != nil {
		return fmt.Errorf("cyclePublish: %w", err)
	}

	for _, article := range articles {
		err := p.AveragePublish(ctx, zModel.PostArticleRequest{Title: article.Title, Content: article.Content})
		if err != nil {
			log.Printf("Error in AveragePublish: %v", err)

			var pErr PublishErr
			if errors.As(err, &pErr) {
				// mark published, if error is PublishErr
				// avoid publish to the same category
				// usually caused by the site is down or the domain is expired
				err = p.dao.MarkPublished(pErr.CateID.String())
				if err != nil {
					return fmt.Errorf("cyclePublish: %w", err)
				}
			} else {
				return fmt.Errorf("cyclePublish: %w", err)
			}

			continue
		}

		err = p.dao.DeleteArticleCache(article.ID.String())
		if err != nil {
			return fmt.Errorf("cyclePublish: %w", err)
		}
	}

	return nil
}
