package commentbot

import (
	"context"
	"fmt"
	"log"
	"time"

	aiAssistInterface "github.com/ray31245/seo_cluster/pkg/ai_assist/ai_assist_interface"
	dbInterface "github.com/ray31245/seo_cluster/pkg/db/db_interface"
	dbModel "github.com/ray31245/seo_cluster/pkg/db/model"
	zModel "github.com/ray31245/seo_cluster/pkg/z_blog_api/model"
	zInterface "github.com/ray31245/seo_cluster/pkg/z_blog_api/z_blog_Interface"
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
			time.Sleep(randomTime())

			err := c.cycleComment(ctx)
			if err != nil {
				log.Println(err)
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
		articles, err := c.listArticleForComment(ctx, site)
		if err != nil {
			return fmt.Errorf("cycleComment: %w", err)
		}

		for _, a := range articles {
			err = c.Comment(ctx, site, a)
			if err != nil {
				return fmt.Errorf("cycleComment: %w", err)
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
		// gape is the probability of comment, if gape=0, then comment
		// if value of gape is more probability of comment is less
		gape := 0
		// Decrease in probability over time
		hours := time.Since(a.PostTime.Time).Hours() + 1
		gape = int(hours) * int(a.CommNums)

		if randomNum() > int32(gape) {
			res = append(res, a)
		}
	}

	return res, nil
}

func (c CommentBot) Comment(ctx context.Context, site dbModel.Site, article zModel.Article) error {
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