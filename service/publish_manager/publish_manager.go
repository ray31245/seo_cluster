package publishmanager

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	dbInterface "github.com/ray31245/seo_cluster/pkg/db/db_interface"
	dbErr "github.com/ray31245/seo_cluster/pkg/db/error"
	dbModel "github.com/ray31245/seo_cluster/pkg/db/model"
	wordpressModel "github.com/ray31245/seo_cluster/pkg/wordpress_api/model"
	wordpressInterface "github.com/ray31245/seo_cluster/pkg/wordpress_api/wordpress_interface"
	zModel "github.com/ray31245/seo_cluster/pkg/z_blog_api/model"
	zInterface "github.com/ray31245/seo_cluster/pkg/z_blog_api/z_blog_Interface"
	"github.com/ray31245/seo_cluster/service/publish_manager/model"
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
	zAPI         zInterface.ZBlogAPI
	wordpressAPI wordpressInterface.WordpressAPI
	dao          DAO
	publishLock  sync.Mutex
}

func NewPublishManager(zAPI zInterface.ZBlogAPI, wordpressAPI wordpressInterface.WordpressAPI, dao DAO) *PublishManager {
	return &PublishManager{
		zAPI:         zAPI,
		wordpressAPI: wordpressAPI,
		dao:          dao,
	}
}

// AveragePublish average publish article to all site and category
func (p *PublishManager) AveragePublish(ctx context.Context, article model.Article) error {
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
	if cate.Site.CmsType == dbModel.CMSTypeWordPress {
		article.CateID = cate.WordpressID
	} else if cate.Site.CmsType == dbModel.CMSTypeZBlog {
		article.CateID = cate.ZBlogID
	} else {
		return fmt.Errorf("AveragePublish: %w", errors.New("cms type not support"))
	}

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

func (p *PublishManager) doPublish(ctx context.Context, article model.Article, site dbModel.Site) error {
	var err error
	if site.CmsType == dbModel.CMSTypeWordPress {
		err = p.doPublishWordPress(ctx, article, site)
	} else if site.CmsType == dbModel.CMSTypeZBlog {
		err = p.doPublishZblog(ctx, article, site)
	} else {
		err = errors.New("cms type not support")
	}

	if err != nil {
		return fmt.Errorf("doPublish: %w", err)
	}

	return nil
}

func (p *PublishManager) doPublishWordPress(ctx context.Context, article model.Article, site dbModel.Site) error {
	// set post article request
	postArticle := wordpressModel.CreateArticleArgs{
		Title:      article.Title,
		Content:    article.Content,
		Categories: []uint32{article.CateID},
		Status:     wordpressModel.StatusPublish,
	}

	// get wordpress api client
	client, err := p.wordpressAPI.GetClient(ctx, site.ID, site.URL, site.UserName, site.Password)
	if err != nil {
		return fmt.Errorf("doPublish: %w", err)
	}

	// post article
	_, err = client.CreateArticle(ctx, postArticle)
	if err != nil {
		return fmt.Errorf("doPublish: %w", err)
	}

	return nil
}

func (p *PublishManager) doPublishZblog(ctx context.Context, article model.Article, site dbModel.Site) error {
	// set post article request
	postArticle := zModel.PostArticleRequest{
		Title:   article.Title,
		Content: article.Content,
		CateID:  article.CateID,
	}

	// get zblog api client
	client, err := p.zAPI.GetClient(ctx, site.ID, site.URL, site.UserName, site.Password)
	if err != nil {
		return fmt.Errorf("doPublishZblog: %w", err)
	}

	// post article
	err = client.PostArticle(ctx, postArticle)
	if err != nil {
		return fmt.Errorf("doPublishZblog: %w", err)
	}

	return nil
}

func (p *PublishManager) PrePublish(article model.Article) error {
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

func (p *PublishManager) StartRandomCyclePublishZblog(ctx context.Context) error {
	lastCategory, err := p.dao.LastPublishedCategoryByCMSType(dbModel.CMSTypeZBlog)
	if err == nil {
		log.Printf("last Publish time %s in StartRandomCyclePublish", lastCategory.LastPublished)
		log.Printf("time now %s in StartRandomCyclePublish", time.Now())

		if time.Since(lastCategory.LastPublished).Minutes() > maxCycleTime {
			log.Println("Duration is more than maxCycleTime in StartRandomCyclePublish, cyclePublish forced to run")

			err = p.CyclePublishZblog(ctx)
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
				if err := p.CyclePublishZblog(ctx); err != nil {
					log.Println("Error during cyclePublish:", err)
				}
			}
		}
	}()

	return nil
}

func (p *PublishManager) CyclePublishZblog(ctx context.Context) error {
	p.publishLock.Lock()
	defer p.publishLock.Unlock()

	return p.cyclePublishZblog(ctx)
}

func (p *PublishManager) cyclePublishZblog(ctx context.Context) error {
	log.Println("cyclePublish running...")

	sites, err := p.dao.ListSitesByCMSType(dbModel.CMSTypeZBlog)
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

	err = p.publishByLack(ctx)
	if err != nil {
		return fmt.Errorf("cyclePublish: %w", err)
	}

	return nil
}

func (p *PublishManager) StartRandomCyclePublishWordPress(ctx context.Context) error {
	lastCategory, err := p.dao.LastPublishedCategoryByCMSType(dbModel.CMSTypeWordPress)
	if err == nil {
		log.Printf("last Publish time %s in StartRandomCyclePublish", lastCategory.LastPublished)
		log.Printf("time now %s in StartRandomCyclePublish", time.Now())

		if time.Since(lastCategory.LastPublished).Minutes() > maxCycleTime {
			log.Println("Duration is more than maxCycleTime in StartRandomCyclePublish, cyclePublish forced to run")

			err = p.CyclePublishWordPress(ctx)
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
				if err := p.CyclePublishWordPress(ctx); err != nil {
					log.Println("Error during cyclePublish:", err)
				}
			}
		}
	}()

	return nil
}

func (p *PublishManager) CyclePublishWordPress(ctx context.Context) error {
	p.publishLock.Lock()
	defer p.publishLock.Unlock()

	return p.cyclePublishWordPress(ctx)
}

func (p *PublishManager) cyclePublishWordPress(ctx context.Context) error {
	log.Println("cyclePublish running...")

	sites, err := p.dao.ListSitesByCMSType(dbModel.CMSTypeWordPress)
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

	err = p.publishByLack(ctx)
	if err != nil {
		return fmt.Errorf("cyclePublish: %w", err)
	}

	return nil
}

func (p *PublishManager) StartPublishByLack(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				// Exit the loop if the context is cancelled
				return
			case <-time.After(time.Minute * 5):
				// Proceed with the publishing cycle after a random duration
				if err := p.PublishByLack(ctx); err != nil {
					log.Println("Error during publishByLack:", err)
				}
			}
		}
	}()

	return
}

func (p *PublishManager) PublishByLack(ctx context.Context) error {
	if ok := p.publishLock.TryLock(); !ok {
		return nil
	}
	defer p.publishLock.Unlock()

	return p.publishByLack(ctx)
}

func (p *PublishManager) publishByLack(ctx context.Context) error {
	// get total lack count
	totalLackCount, err := p.dao.SumLackCount()
	if err != nil {
		return fmt.Errorf("publishByLack: %w", err)
	}

	articles, err := p.dao.ListArticleCacheByLimit(totalLackCount)
	if err != nil {
		return fmt.Errorf("publishByLack: %w", err)
	}

	for _, article := range articles {
		err := p.AveragePublish(ctx, model.Article{Title: article.Title, Content: article.Content})
		if err != nil {
			log.Printf("Error in AveragePublish: %v", err)

			var pErr PublishErr
			if errors.As(err, &pErr) {
				// mark published, if error is PublishErr
				// avoid publish to the same category
				// usually caused by the site is down or the domain is expired
				err = p.dao.MarkPublished(pErr.CateID.String())
				if err != nil {
					return fmt.Errorf("publishByLack: %w", err)
				}
			} else {
				return fmt.Errorf("publishByLack: %w", err)
			}

			continue
		}

		err = p.dao.DeleteArticleCache(article.ID.String())
		if err != nil {
			return fmt.Errorf("publishByLack: %w", err)
		}
	}

	return nil
}

func (p *PublishManager) CountArticleCache() (int64, error) {
	return p.dao.CountArticleCache()
}
