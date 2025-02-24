package publishmanager

import (
	"context"
	"errors"
	"fmt"
	"log"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	aiAssistInterface "github.com/ray31245/seo_cluster/pkg/ai_assist/ai_assist_interface"
	aiAssistModel "github.com/ray31245/seo_cluster/pkg/ai_assist/model"
	dbInterface "github.com/ray31245/seo_cluster/pkg/db/db_interface"
	dbErr "github.com/ray31245/seo_cluster/pkg/db/error"
	dbModel "github.com/ray31245/seo_cluster/pkg/db/model"
	"github.com/ray31245/seo_cluster/pkg/util"
	wordpressModel "github.com/ray31245/seo_cluster/pkg/wordpress_api/model"
	wordpressInterface "github.com/ray31245/seo_cluster/pkg/wordpress_api/wordpress_interface"
	zModel "github.com/ray31245/seo_cluster/pkg/z_blog_api/model"
	zInterface "github.com/ray31245/seo_cluster/pkg/z_blog_api/z_blog_Interface"
	"github.com/ray31245/seo_cluster/service/publish_manager/model"
)

const (
	ConfigUnCateName  = "un_cate_name"
	TagsBlockList     = "tags_block_list"
	IsStopAutoPublish = "is_stop_auto_publish"

	maxKeyWords                  = 5
	updateArticleTagSignalBuffer = 100000
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
	dbInterface.KVConfigDAOInterface
}

type PublishManager struct {
	zAPI                    zInterface.ZBlogAPI
	wordpressAPI            wordpressInterface.WordpressAPI
	aiAssist                aiAssistInterface.AIAssistInterface
	dao                     DAO
	publishLock             sync.Mutex
	updateTagSignal         chan updateArticleTagSignal
	maxUpdateTagThreads     int
	updateArticleTagThreads atomic.Int32
}

var ErrStopAutoPublish = errors.New("system is set to stop auto publish, break the cycle")

func NewPublishManager(zAPI zInterface.ZBlogAPI, wordpressAPI wordpressInterface.WordpressAPI, dao DAO, aiAssist aiAssistInterface.AIAssistInterface) *PublishManager {
	updateArticleTagSignal := make(chan updateArticleTagSignal, updateArticleTagSignalBuffer)

	return &PublishManager{
		zAPI:            zAPI,
		wordpressAPI:    wordpressAPI,
		aiAssist:        aiAssist,
		dao:             dao,
		updateTagSignal: updateArticleTagSignal,
	}
}

// AveragePublish average publish article to all site and category
func (p *PublishManager) AveragePublish(ctx context.Context, article model.Article) error {
	isStopAutoPublish, err := p.dao.GetBoolByKeyWithDefault(IsStopAutoPublish, false)
	if err != nil {
		return fmt.Errorf("AveragePublish: %w", err)
	}

	if isStopAutoPublish {
		return fmt.Errorf("AveragePublish: %w", ErrStopAutoPublish)
	}

	cate, err := p.findFirstMatchCategory(ctx, article)
	if err != nil {
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

func (p *PublishManager) DirectPublish(ctx context.Context, cateID string, article model.Article) error {
	cate, err := p.dao.GetCategory(cateID)
	if err != nil {
		return fmt.Errorf("DirectPublish: %w", err)
	}

	site, err := p.dao.GetSite(cate.SiteID.String())
	if err != nil {
		return fmt.Errorf("DirectPublish: %w", err)
	}

	// set category id
	if site.CmsType == dbModel.CMSTypeWordPress {
		article.CateID = cate.WordpressID
	} else if site.CmsType == dbModel.CMSTypeZBlog {
		article.CateID = cate.ZBlogID
	} else {
		return fmt.Errorf("DirectPublish: %w", errors.New("cms type not support"))
	}

	err = p.doPublish(ctx, article, cate.Site)
	if err != nil {
		return errors.Join(PublishErr{SiteID: cate.SiteID, CateID: cate.ID}, err)
	}

	return nil
}

type PublishSignal struct {
	Site    dbModel.Site
	Article model.Article
}

func (p *PublishManager) BroadcastPublish(ctx context.Context, article model.Article) error {
	sites, err := p.dao.ListSites()
	if err != nil {
		return fmt.Errorf("BroadcastPublish: %w", err)
	}

	sitesCount := len(sites)

	if sitesCount == 0 {
		return fmt.Errorf("BroadcastPublish: %w", errors.New("no site found"))
	}

	threads := runtime.NumCPU() * 10
	threads = min(threads, sitesCount)

	signal := make(chan PublishSignal, sitesCount)
	defer close(signal)

	errCh := make(chan error, sitesCount)
	defer close(errCh)

	go func() {
		for _, site := range sites {
			signal <- PublishSignal{Site: site, Article: article}
		}
	}()

	var wg *sync.WaitGroup = new(sync.WaitGroup)

	wg.Add(threads)

	for range threads {
		go p.startPublishWorker(ctx, wg, signal, errCh)
	}

	var errs error

FanInLoop:
	for {
		select {
		case err := <-errCh:
			errs = errors.Join(errs, err)
		case <-ctx.Done():
			return nil
		case <-util.WaitGroupChan(wg):
			break FanInLoop
		}
	}

	var resErr error

	if errs != nil {
		resErr = fmt.Errorf("BroadcastPublish: %w", errs)
	}

	return resErr
}

func (p *PublishManager) startPublishWorker(ctx context.Context, wg *sync.WaitGroup, signal chan PublishSignal, errCh chan error) {
	defer wg.Done()

	for {
		select {
		case s := <-signal:
			site, err := p.dao.GetSite(s.Site.ID.String())
			if err != nil {
				errCh <- err

				continue
			}

			cate, err := p.MatchCategory(ctx, site.Categories, s.Article)
			if err != nil {
				errCh <- err

				continue
			}

			// set category id
			if cate.Site.CmsType == dbModel.CMSTypeWordPress {
				s.Article.CateID = cate.WordpressID
			} else if cate.Site.CmsType == dbModel.CMSTypeZBlog {
				s.Article.CateID = cate.ZBlogID
			} else {
				errCh <- errors.New("cms type not support")

				continue
			}

			err = p.doPublish(ctx, s.Article, cate.Site)
			if err != nil {
				errCh <- errors.Join(PublishErr{SiteID: cate.SiteID, CateID: cate.ID}, err)

				continue
			}
		case <-ctx.Done():
			return

		default:
			return
		}
	}
}

func (p *PublishManager) findFirstMatchCategory(ctx context.Context, article model.Article) (*dbModel.Category, error) {
	cates, err := p.dao.ListPublishedCategories()
	if err != nil {
		return nil, fmt.Errorf("FindFirstMatchCategory: %w", err)
	}

	if len(cates) == 0 {
		return nil, fmt.Errorf("FindFirstMatchCategory: %w", ErrNoCategoryNeedToBePublished)
	}

	cate, err := p.MatchCategory(ctx, cates, article)
	if err != nil {
		return nil, fmt.Errorf("FindFirstMatchCategory: %w", err)
	}

	return cate, nil
}

func (p *PublishManager) MatchCategory(ctx context.Context, cates []dbModel.Category, article model.Article) (*dbModel.Category, error) {
	notMatchCate := dbModel.Category{}
	cateOpts := []aiAssistModel.CategoryOption{}

	configUnCateName, err := p.dao.GetByKeyWithDefault(ConfigUnCateName, "ThisIsUnCate")
	if err != nil {
		return nil, fmt.Errorf("MatchCategory: %w", err)
	}

	for _, cate := range cates {
		if strings.Trim(cate.Name, " ") == configUnCateName.Value {
			if notMatchCate.ID.String() == "" {
				{
					notMatchCate = cate
				}

				continue
			}
		}

		cateOpts = append(cateOpts, aiAssistModel.CategoryOption{ID: cate.ID.String(), Name: cate.Name})
	}

	selectResp, err := p.aiAssist.SelectCategory(ctx,
		aiAssistModel.SelectCategoryRequest{
			Text:             []byte(article.Content),
			CategoriesOption: cateOpts,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("MatchCategory: %w", err)
	}

	var cate *dbModel.Category

	switch {
	case selectResp.ID != "" && selectResp.IsFind:
		cate, err = p.dao.GetCategory(selectResp.ID)
		if err != nil && !dbErr.IsNotfoundErr(err) {
			return nil, fmt.Errorf("MatchCategory: %w", err)
		} else if err == nil {
			break
		}
		// if not found, fall back to notMatchCate or first published category
		fallthrough
	case notMatchCate.Name != "":
		cate = &notMatchCate
	default:
		// find first published category
		cate = &cates[0]
	}

	return cate, nil
}

func (p *PublishManager) doPublish(ctx context.Context, article model.Article, site dbModel.Site) error {
	var err error
	if site.CmsType == dbModel.CMSTypeWordPress {
		_, err = p.doPublishWordPress(ctx, article, site)
	} else if site.CmsType == dbModel.CMSTypeZBlog {
		_, err = p.doPublishZblog(ctx, article, site)
	} else {
		err = errors.New("cms type not support")
	}

	if err != nil {
		return fmt.Errorf("doPublish: %w", err)
	}

	return nil
}

func (p *PublishManager) doPublishWordPress(ctx context.Context, article model.Article, site dbModel.Site) (wordpressModel.CreateArticleResponse, error) {
	// set post article request
	postArticle := article.ToWordpressCreateArgs(wordpressModel.StatusPublish)

	// get wordpress api client
	client, err := p.wordpressAPI.GetClient(ctx, site.ID, site.URL, site.UserName, site.Password)
	if err != nil {
		return wordpressModel.CreateArticleResponse{}, fmt.Errorf("doPublishWordPress: %w", err)
	}

	// post article
	postArt, err := client.CreateArticle(ctx, postArticle)
	if err != nil {
		return wordpressModel.CreateArticleResponse{}, fmt.Errorf("doPublishWordPress: %w", err)
	}

	// update article tag
	err = p.pushUpdateTagSignal(updateArticleTagSignal{ArtContent: article.Content, ArtID: postArt.ID, Site: site})
	if err != nil {
		return wordpressModel.CreateArticleResponse{}, fmt.Errorf("doPublishWordPress: %w", err)
	}

	return postArt, nil
}

func (p *PublishManager) doPublishZblog(ctx context.Context, article model.Article, site dbModel.Site) (zModel.Article, error) {
	// set post article request
	postArticle := article.ToZBlogCreateRequest()

	// get zblog api client
	client, err := p.zAPI.GetClient(ctx, site.ID, site.URL, site.UserName, site.Password)
	if err != nil {
		return zModel.Article{}, fmt.Errorf("doPublishZblog: %w", err)
	}

	// post article
	postArt, err := client.PostArticle(ctx, postArticle)
	if err != nil {
		return zModel.Article{}, fmt.Errorf("doPublishZblog: %w", err)
	}

	artID, err := strconv.Atoi(string(postArt.ID))
	if err != nil {
		return zModel.Article{}, fmt.Errorf("doPublishZblog: %w", err)
	}

	// update article tag
	err = p.pushUpdateTagSignal(updateArticleTagSignal{ArtContent: article.Content, ArtID: artID, Site: site})
	if err != nil {
		return zModel.Article{}, fmt.Errorf("doPublishZblog: %w", err)
	}

	return postArt, nil
}

type updateArticleTagSignal struct {
	ArtContent string
	ArtID      int
	Site       dbModel.Site
}

func (p *PublishManager) StartUpdateArticleTagSignalLoop(ctx context.Context, threads int, maxThreads int) error {
	p.maxUpdateTagThreads = maxThreads

	var err error
	for i := 0; i < threads; i++ {
		err = p.newUpdateArticleTagSignalLoopThread(ctx, true)
		if err != nil {
			return fmt.Errorf("StartUpdateArticleTagSignalLoop: %w", err)
		}
	}

	return nil
}

func (p *PublishManager) pushUpdateTagSignal(signal updateArticleTagSignal) (err error) {
	select {
	case p.updateTagSignal <- signal:
	default:
		log.Println("updateTagSignal is full, open new goroutine to handle")

		err = p.newUpdateArticleTagSignalLoopThread(context.Background(), false)
		if err != nil {
			return fmt.Errorf("pushUpdateTagSignal: %w", err)
		}
		p.updateTagSignal <- signal
	}

	return
}

func (p *PublishManager) newUpdateArticleTagSignalLoopThread(ctx context.Context, isPersistent bool) error {
	if p.updateArticleTagThreads.Load() >= int32(p.maxUpdateTagThreads) {
		return fmt.Errorf("newUpdateArticleTagSignalLoopThread: %w", errors.New("max threads reached"))
	}

	go p.updateArticleTagSignalLoop(ctx, isPersistent)

	return nil
}

func (p *PublishManager) updateArticleTagSignalLoop(ctx context.Context, isPersistent bool) {
	p.updateArticleTagThreads.Add(1)
	defer p.updateArticleTagThreads.Add(-1)

	isIdle := false
	idleCheck := func() <-chan time.Time {
		if !isPersistent {
			return time.Tick(5 * time.Minute)
		}

		return make(<-chan time.Time)
	}

	for {
		select {
		case signal := <-p.updateTagSignal:
			isIdle = false

			err := p.updateArticleTag(ctx, signal.ArtContent, signal.ArtID, signal.Site)
			if err != nil {
				log.Printf("Error in updateArticleTagSignalLoop: %v", err)
			}
		case <-idleCheck():
			if isIdle {
				log.Println("updateArticleTagSignalLoop is idle, exit")

				return
			}

			isIdle = true
		case <-ctx.Done():
			return
		}
	}
}

func (p *PublishManager) updateArticleTag(ctx context.Context, artContent string, artID int, site dbModel.Site) error {
	var err error
	if site.CmsType == dbModel.CMSTypeWordPress {
		err = p.updateArticleTagWordpress(ctx, artContent, artID, site)
	} else if site.CmsType == dbModel.CMSTypeZBlog {
		err = p.updateArticleTagZblog(ctx, artContent, artID, site)
	} else {
		err = errors.New("cms type not support")
	}

	if err != nil {
		return fmt.Errorf("updateArticleTag: %w", err)
	}

	return nil
}

func (p *PublishManager) updateArticleTagZblog(ctx context.Context, artContent string, artID int, site dbModel.Site) error {
	client, err := p.zAPI.GetClient(ctx, site.ID, site.URL, site.UserName, site.Password)
	if err != nil {
		return fmt.Errorf("updateArticleTagZblog: %w", err)
	}

	tagBlackList, err := p.GetTagsBlockList()
	if err != nil && !dbErr.IsNotfoundErr(err) {
		return fmt.Errorf("updateArticleTagZblog: %w", err)
	}

	siteTags, err := client.ListTagAll(ctx)
	if err != nil {
		return fmt.Errorf("updateArticleTagZblog: %w", err)
	}

	tagMatcher := NewTagMatcher(tagBlackList, siteTags)

	matchedTags := []string{}

	// find matched tags
	keywords, err := p.aiAssist.FindKeyWords(ctx, []byte(artContent))
	if err != nil {
		return fmt.Errorf("updateArticleTagZblog: %w", err)
	}

	for _, keyword := range keywords.KeyWords {
		if tagMatcher.IsTagBlackList(keyword) {
			continue
		}

		isMatch, matchTag := tagMatcher.IsTagInSite(keyword)
		if isMatch {
			keyword = matchTag.GetName()
		} else {
			newTag, err := client.PostTag(ctx, zModel.PostTagRequest{Name: keyword})
			if err != nil {
				log.Printf("Error in PostTag: %v, Err msg: %v", keyword, err)

				continue
			}

			keyword = newTag.Name
		}

		matchedTags = append(matchedTags, keyword)

		if len(matchedTags) >= maxKeyWords {
			break
		}
	}

	_, err = client.PostArticle(ctx, zModel.PostArticleRequest{
		ID:  uint32(artID),
		Tag: strings.Join(matchedTags, ","),
	})
	if err != nil {
		return fmt.Errorf("updateArticleTagZblog: %w", err)
	}

	return nil
}

func (p *PublishManager) updateArticleTagWordpress(ctx context.Context, artContent string, artID int, site dbModel.Site) error {
	client, err := p.wordpressAPI.GetClient(ctx, site.ID, site.URL, site.UserName, site.Password)
	if err != nil {
		return fmt.Errorf("updateArticleWordpress: %w", err)
	}

	tagBlackList, err := p.GetTagsBlockList()
	if err != nil && !dbErr.IsNotfoundErr(err) {
		return fmt.Errorf("updateArticleWordpress: %w", err)
	}

	siteTags, err := client.ListTagAll(ctx)
	if err != nil {
		return fmt.Errorf("updateArticleWordpress: %w", err)
	}

	tagMatcher := NewTagMatcher(tagBlackList, siteTags)

	matchedTags := []int{}

	// find matched tags
	keywords, err := p.aiAssist.FindKeyWords(ctx, []byte(artContent))
	if err != nil {
		return fmt.Errorf("updateArticleWordpress: %w", err)
	}

	for _, keyword := range keywords.KeyWords {
		if tagMatcher.IsTagBlackList(keyword) {
			continue
		}

		if isMatch, matchTag := tagMatcher.IsTagInSite(keyword); isMatch {
			matchedTags = append(matchedTags, matchTag.GetID())
		} else {
			newTag, err := client.CreateTag(ctx, wordpressModel.CreateTagArgs{Name: keyword})
			if err != nil {
				log.Printf("Error in PostTag: %v, Err msg: %v", keyword, err)

				continue
			}

			matchedTags = append(matchedTags, newTag.ID)
		}

		if len(matchedTags) >= maxKeyWords {
			break
		}
	}

	_, err = client.UpdateArticle(ctx, wordpressModel.UpdateArticleArgs{
		ID:   artID,
		Tags: matchedTags,
	})
	if err != nil {
		return fmt.Errorf("updateArticleWordpress: %w", err)
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
		log.Printf("last Publish time %s in StartRandomCyclePublishZblog", lastCategory.LastPublished)
		log.Printf("time now %s in StartRandomCyclePublishZblog", time.Now())

		if time.Since(lastCategory.LastPublished).Minutes() > maxCycleTime {
			log.Println("Duration is more than maxCycleTime in StartRandomCyclePublishZblog, cyclePublish forced to run")

			err = p.CyclePublishZblog(ctx)
			if err != nil {
				return fmt.Errorf("StartRandomCyclePublishZblog: %w", err)
			}
		}
	} else if !dbErr.IsNotfoundErr(err) {
		return fmt.Errorf("StartRandomCyclePublishZblog: %w", err)
	}

	go func() {
		for {
			nextTime := randomTime()

			multi, err := p.multiOfArticleCount()
			if err != nil {
				log.Println("Error in multiOfArticleCount:", err)
			} else if multi > 0 {
				nextTime = nextTime / time.Duration(multi)
			}

			log.Printf("next time run cyclePublish is %s in StartRandomCyclePublishZblog", time.Now().Add(nextTime))
			select {
			case <-ctx.Done():
				// Exit the loop if the context is cancelled
				return
			case <-time.After(nextTime):
				// Proceed with the publishing cycle after a random duration
				if err := p.CyclePublishZblog(ctx); err != nil {
					log.Println("Error during CyclePublishZblog:", err)
				}
			}
		}
	}()

	return nil
}

func (p *PublishManager) multiOfArticleCount() (int, error) {
	count, err := p.CountArticleCache()
	if err != nil {
		return 0, fmt.Errorf("multiOfArticleCount: %w", err)
	}

	sites, err := p.dao.ListSites()
	if err != nil {
		return 0, fmt.Errorf("multiOfArticleCount: %w", err)
	}

	// expected count, preserver 200 articles
	expectedCount := 200
	for _, site := range sites {
		expectedCount += site.LackCount
		if site.CmsType == dbModel.CMSTypeZBlog {
			expectedCount += 1
		} else if site.CmsType == dbModel.CMSTypeWordPress {
			expectedCount += 5
		} else {
			return 0, fmt.Errorf("multiOfArticleCount: %w", errors.New("cms type not support"))
		}
	}

	// if count is less than expectedCount, return 0
	if int(count) <= expectedCount {
		return 0, nil
	}

	// if count is more than expectedCount, return multi
	multi := int(count) / expectedCount

	return multi, nil
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
			log.Printf("site id %s, lack count %d in cyclePublishZblog", site.ID, lackCount)

			err := p.dao.IncreaseLackCount(site.ID.String(), int(lackCount))
			if err != nil {
				return fmt.Errorf("cyclePublishZblog: %w", err)
			}
		}
	}

	err = p.publishByLack(ctx)
	if err != nil {
		return fmt.Errorf("cyclePublishZblog: %w", err)
	}

	return nil
}

// StartRandomCyclePublishWordPress start random cycle publish wordpress
// do 4 or 5 times publish in a day
// 2 times in 12 hours
// 2 times in next 12 hours
// may extra 1 times in next 24 hours
func (p *PublishManager) StartRandomCyclePublishWordPress(ctx context.Context) error {
	go func() {
		for {
			multi, err := p.multiOfArticleCount()
			if err != nil {
				log.Println("Error in multiOfArticleCount:", err)
			}

			if multi <= 0 {
				multi = 1
			}

			timeArr := []time.Time{}

			// TODO: optimize wordpress, and delete this
			multi = 1

			for range multi {
				timeArr = append(timeArr, computeTimePointArray()...)
			}

			sort.Slice(timeArr, func(i, j int) bool {
				return timeArr[i].Before(timeArr[j])
			})

			log.Printf("timeArr is %v in StartRandomCyclePublishWordPress", timeArr)
			select {
			case <-ctx.Done():
				return
			default:
				timeArrSchedule(ctx, timeArr, func() {
					err := p.CyclePublishWordPress(ctx)
					if err != nil {
						log.Println("Error during CyclePublishWordPress:", err)
					}
				})
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
	log.Println("cyclePublishWordPress running...")

	sites, err := p.dao.ListSitesByCMSType(dbModel.CMSTypeWordPress)
	if err != nil {
		return fmt.Errorf("cyclePublishWordPress: %w", err)
	}

	for _, site := range sites {
		if site.LackCount != 0 {
			continue
		}

		log.Printf("site id %s, lack count %d in cyclePublishWordPress", site.ID, 1)

		err := p.dao.IncreaseLackCount(site.ID.String(), 1)
		if err != nil {
			return fmt.Errorf("cyclePublishWordPress: %w", err)
		}
	}

	err = p.publishByLack(ctx)
	if err != nil {
		return fmt.Errorf("cyclePublishWordPress: %w", err)
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

	articles, err := p.dao.ListReadyToPublishArticleCacheByLimit(totalLackCount)
	if err != nil {
		return fmt.Errorf("publishByLack: %w", err)
	}

	articleIDs := []string{}
	for _, article := range articles {
		articleIDs = append(articleIDs, article.ID.String())
	}

	err = p.dao.UpdateArticleCacheStatusByIDs(articleIDs, dbModel.ArticleCacheStatusInBuffer)
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

		err = p.dao.DeleteArticleCacheByIDs([]string{article.ID.String()})
		if err != nil {
			return fmt.Errorf("publishByLack: %w", err)
		}
	}

	return nil
}

func (p *PublishManager) CountArticleCache() (int64, error) {
	return p.dao.CountArticleCache()
}

func (p *PublishManager) SetConfigUnCateName(ctx context.Context, name string) error {
	return p.dao.UpsertByKey(ConfigUnCateName, name)
}

func (p *PublishManager) GetConfigUnCateName() (string, error) {
	res, err := p.dao.GetByKey(ConfigUnCateName)
	if err != nil {
		return "", fmt.Errorf("GetConfigUnCateName: %w", err)
	}

	return res.Value, nil
}

func (p *PublishManager) SetTagsBlockList(ctx context.Context, tags []string) error {
	return p.dao.UpsertByKey(TagsBlockList, strings.Join(tags, ","))
}

func (p *PublishManager) GetTagsBlockList() ([]string, error) {
	res, err := p.dao.GetByKey(TagsBlockList)
	if err != nil {
		return nil, fmt.Errorf("GetTagsBlockList: %w", err)
	}

	return strings.Split(res.Value, ","), nil
}

func (p *PublishManager) StopAutoPublish() error {
	return p.dao.UpsertByKeyBool(IsStopAutoPublish, true)
}

func (p *PublishManager) StartAutoPublish() error {
	return p.dao.UpsertByKeyBool(IsStopAutoPublish, false)
}

func (p *PublishManager) StopAutoPublishStatus() (bool, error) {
	res, err := p.dao.GetBoolByKeyWithDefault(IsStopAutoPublish, false)
	if err != nil {
		return false, fmt.Errorf("AutoPublishStatus: %w", err)
	}

	return res, nil
}
