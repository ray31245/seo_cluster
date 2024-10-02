package commentbot

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	aiAssistInterface "github.com/ray31245/seo_cluster/pkg/ai_assist/ai_assist_interface"
	dbInterface "github.com/ray31245/seo_cluster/pkg/db/db_interface"
	dbModel "github.com/ray31245/seo_cluster/pkg/db/model"
	zModel "github.com/ray31245/seo_cluster/pkg/z_blog_api/model"
	zInterface "github.com/ray31245/seo_cluster/pkg/z_blog_api/z_blog_Interface"
)

const (
	// coefficientOfGape is the coefficient of gape
	// if want increase probability of comment, decrease this value
	coefficientOfGape = 10
	// rateLimitDelay is the delay between each comment
	// to avoid rate limit of Gemini( googleapi: Error 429: Resource has been exhausted (e.g. check quota) )
	rateLimitDelay = time.Millisecond * 500
)

type CommentBot struct {
	zBlogAPI       zInterface.ZBlogAPI
	siteDAO        dbInterface.SiteDAOInterface
	commentUserDAO dbInterface.CommentUserDAOInterface
	aiAssist       aiAssistInterface.AIAssistInterface
}

func NewCommentBot(zBlogAPI zInterface.ZBlogAPI, siteDAO dbInterface.SiteDAOInterface, commentUserDAO dbInterface.CommentUserDAOInterface, aiAssist aiAssistInterface.AIAssistInterface) *CommentBot {
	return &CommentBot{
		zBlogAPI:       zBlogAPI,
		siteDAO:        siteDAO,
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
	log.Println("cycleComment running...")

	sites, err := c.siteDAO.ListSites()
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

		for _, a := range articles {
			err = c.Comment(ctx, site, a)
			if err != nil {
				log.Printf("site url %s, article id %s, Comment error: %v", site.URL, a.ID, err)

				continue
			}

			for {
				select {
				case <-ctx.Done():
					return nil
				case <-time.After(rateLimitDelay):
					break
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
		Sortby: "PostTime",
		Order:  "desc",
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

		log.Printf("site url %s, article id %s, gape: %d, randomN: %d", site.URL, a.ID, gap, randomN)
	}

	return res, nil
}

func computeGap(article zModel.Article) int {
	// observe the time of post and the number of comments
	log.Printf("article id %s, post time: %s, commNums: %d", article.ID, article.PostTime, article.CommNums)

	hours := time.Since(article.PostTime.Time).Hours() + 1

	return int(math.Sqrt(hours)*coefficientOfGape*float64(article.CommNums+1)) - int(hours)*5
}

func (c CommentBot) Comment(ctx context.Context, site dbModel.Site, article zModel.Article) error {
	log.Printf("site url %s, article id %s, start comment", site.URL, article.ID)

	commentUser, err := c.commentUserDAO.GetRandomCommentUser()
	if err != nil {
		return fmt.Errorf("Comment: %w", err)
	}

	client, err := c.zBlogAPI.GetClient(ctx, commentUser.ID, site.URL, commentUser.Name, commentUser.Password)
	if err != nil {
		return fmt.Errorf("Comment: %w", err)
	}

	article, err = client.GetArticle(ctx, article.ID)
	if err != nil {
		return fmt.Errorf("Comment: %w", err)
	}

	// comment
	comment, err := c.aiAssist.Comment(ctx, []byte(article.Content))
	if err != nil {
		return fmt.Errorf("Comment: %w", err)
	}

	err = client.PostComment(ctx, zModel.PostCommentRequest{LogID: article.ID, Content: comment.Comment})
	if err != nil {
		return fmt.Errorf("Comment: %w", err)
	}

	log.Printf("site url %s, article id %s,score: %d, comment success", site.URL, article.ID, comment.Score)

	return nil
}
