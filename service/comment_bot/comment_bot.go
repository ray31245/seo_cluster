package commentbot

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	aiAssistInterface "github.com/ray31245/seo_cluster/pkg/ai_assist/ai_assist_interface"
	"github.com/ray31245/seo_cluster/pkg/ai_assist/model"
	dbInterface "github.com/ray31245/seo_cluster/pkg/db/db_interface"
	dbModel "github.com/ray31245/seo_cluster/pkg/db/model"
	zModel "github.com/ray31245/seo_cluster/pkg/z_blog_api/model"
	zInterface "github.com/ray31245/seo_cluster/pkg/z_blog_api/z_blog_Interface"
)

const (
	IsStopAutoCommentConfigKey = "IsStopAutoComment"
	// coefficientOfGape is the coefficient of gape
	// if want increase probability of comment, decrease this value
	coefficientOfGape = 30
	// rateLimitDelay is the delay between each comment
	// to avoid rate limit of Gemini( googleapi: Error 429: Resource has been exhausted (e.g. check quota) )
	rateLimitDelay        = time.Millisecond * 500
	maxContinueErrorCount = 3
)

type CommentBot struct {
	zBlogAPI       zInterface.ZBlogAPI
	siteDAO        dbInterface.SiteDAOInterface
	configDAO      dbInterface.KVConfigDAOInterface
	commentUserDAO dbInterface.CommentUserDAOInterface
	aiAssist       aiAssistInterface.AIAssistInterface
}

func NewCommentBot(zBlogAPI zInterface.ZBlogAPI, configDAO dbInterface.KVConfigDAOInterface, siteDAO dbInterface.SiteDAOInterface, commentUserDAO dbInterface.CommentUserDAOInterface, aiAssist aiAssistInterface.AIAssistInterface) *CommentBot {
	return &CommentBot{
		zBlogAPI:       zBlogAPI,
		siteDAO:        siteDAO,
		configDAO:      configDAO,
		commentUserDAO: commentUserDAO,
		aiAssist:       aiAssist,
	}
}

func (c CommentBot) StartCycleComment(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				// Exit the loop if the context is cancelled
				return
			case <-time.After(randomTime()):
				// Proceed with the publishing cycle after a random duration
				if err := c.cycleComment(ctx); err != nil {
					log.Printf("Error in cycleComment: %v", err)
				}
			}
		}
	}()
}

func (c CommentBot) cycleComment(ctx context.Context) error {
	// check if auto comment is stopped
	isStopAutoComment, err := c.configDAO.GetBoolByKeyWithDefault(IsStopAutoCommentConfigKey, false)
	if err != nil {
		return fmt.Errorf("cycleComment: %w", err)
	}

	if isStopAutoComment {
		log.Println("cycleComment is stopped")

		return nil
	}

	log.Println("cycleComment running...")

	sites, err := c.siteDAO.ListSitesRandom()
	if err != nil {
		return fmt.Errorf("cycleComment: %w", err)
	}

	for _, site := range sites {
		log.Printf("site url %s, start cycleComment", site.URL)

		articles, err := c.listArticleForComment(ctx, site)
		if err != nil {
			log.Printf("site url %s, listArticleForComment error: %v", site.URL, err)

			continue
		}

		continueErrorCount := 0

		for _, a := range articles {
			err = c.Comment(ctx, site, a)
			if err != nil {
				continueErrorCount++

				log.Printf("site url %s, article id %s, Comment error: %v", site.URL, a.ID, err)

				if continueErrorCount > maxContinueErrorCount {
					log.Printf("site url %s, article id %s, Comment error count > 3, break", site.URL, a.ID)

					break
				}

				continue
			}

			continueErrorCount = 0

		rateLimitDelay:
			for {
				select {
				case <-ctx.Done():
					return nil
				case <-time.After(rateLimitDelay):
					break rateLimitDelay
				}
			}
		}
	}

	return nil
}

func (c CommentBot) listArticleForComment(ctx context.Context, site dbModel.Site) ([]zModel.Article, error) {
	res := []zModel.Article{}
	client := c.zBlogAPI.NewAnonymousClient(ctx, site.URL)
	// get all articles
	articles, err := client.ListArticle(ctx, zModel.ListArticleRequest{
		PageRequest: zModel.PageRequest{
			SortBy: "PostTime",
			Order:  "desc",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("listArticleForComment: %w", err)
	}

	for _, a := range articles {
		gap := computeGap(a)

		randomN := randomNum()
		if randomN > int32(gap) {
			res = append(res, a)
		}

		// log.Printf("site url %s, article id %s, gape: %d, randomN: %d", site.URL, a.ID, gap, randomN)
	}

	return res, nil
}

func computeGap(article zModel.Article) int {
	hours := time.Since(article.PostTime.Time).Hours() + 1
	// hours should be at least 1
	if hours < 1 {
		hours = 1
	}

	// observe the time of post and the number of comments
	// log.Printf("article id %s, post time: %s,hours: %f, commNums: %d", article.ID, article.PostTime, hours, article.CommNums)

	return int(math.Sqrt(hours)*coefficientOfGape*float64(article.CommNums+1)) - int(hours)*2
}

func (c CommentBot) Comment(ctx context.Context, site dbModel.Site, article zModel.Article) error {
	log.Printf("site url %s, article id %s, start comment", site.URL, article.ID)

	var err error
	if site.CmsType == dbModel.CMSTypeWordPress {
		err = c.commentWordPress(ctx, site, article)
	} else if site.CmsType == dbModel.CMSTypeZBlog {
		err = c.commentZBlog(ctx, site, article)
	} else {
		err = errors.New("unsupported cms type")
	}

	if err != nil {
		return fmt.Errorf("Comment: %w", err)
	}

	return nil
}

func (c CommentBot) commentZBlog(ctx context.Context, site dbModel.Site, article zModel.Article) error {
	commentUser, err := c.commentUserDAO.GetRandomCommentUser()
	if err != nil {
		return fmt.Errorf("commentZBlog: %w", err)
	}

	client, err := c.zBlogAPI.GetClient(ctx, commentUser.ID, site.URL, commentUser.Name, commentUser.Password)
	if err != nil {
		return fmt.Errorf("commentZBlog: %w", err)
	}

	article, err = client.GetArticle(ctx, string(article.ID))
	if err != nil {
		return fmt.Errorf("commentZBlog: %w", err)
	}

	// comment
	comment, err := c.comment(ctx, article)
	if err != nil {
		return fmt.Errorf("commentZBlog: %w", err)
	}

	err = client.PostComment(ctx, zModel.PostCommentRequest{LogID: string(article.ID), Content: comment.Comment})
	if err != nil {
		return fmt.Errorf("commentZBlog: %w", err)
	}

	log.Printf("site url %s, article id %s,score: %d, comment success", site.URL, article.ID, comment.Score)

	return nil
}

func (c CommentBot) commentWordPress(ctx context.Context, site dbModel.Site, article zModel.Article) error {
	// TODO: implement commentWordPress
	return nil
}

func (c CommentBot) comment(ctx context.Context, article zModel.Article) (model.CommentResponse, error) {
	if ok := c.aiAssist.TryLock(); !ok {
		return model.CommentResponse{}, fmt.Errorf("Comment: %w", errors.New("AIAssist is locked"))
	}
	defer c.aiAssist.Unlock()

	comment, err := c.aiAssist.Comment(ctx, []byte(article.Content))
	if err != nil {
		return model.CommentResponse{}, fmt.Errorf("Comment: %w", err)
	}

	return comment, nil
}

func (c CommentBot) StopAutoComment(ctx context.Context) error {
	err := c.configDAO.UpsertByKeyBool(IsStopAutoCommentConfigKey, true)
	if err != nil {
		return fmt.Errorf("StopAutoComment: %w", err)
	}

	return nil
}

func (c CommentBot) StartAutoComment(ctx context.Context) error {
	err := c.configDAO.UpsertByKeyBool(IsStopAutoCommentConfigKey, false)
	if err != nil {
		return fmt.Errorf("StartAutoComment: %w", err)
	}

	return nil
}

func (c CommentBot) IsAutoCommentStopped() (bool, error) {
	isStopAutoComment, err := c.configDAO.GetBoolByKeyWithDefault(IsStopAutoCommentConfigKey, false)
	if err != nil {
		return false, fmt.Errorf("IsAutoCommentStopped: %w", err)
	}

	return isStopAutoComment, nil
}
