package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/ray31245/seo_cluster/cmd/publish_manager_service/handler"
	aiassist "github.com/ray31245/seo_cluster/pkg/ai_assist"
	"github.com/ray31245/seo_cluster/pkg/db"
	zBlogApi "github.com/ray31245/seo_cluster/pkg/z_blog_api"
	commentbot "github.com/ray31245/seo_cluster/service/comment_bot"
	publishManager "github.com/ray31245/seo_cluster/service/publish_manager"
	sitemanager "github.com/ray31245/seo_cluster/service/site_manager"

	"github.com/gin-gonic/gin"
)

var APIKey string //nolint:gochecknoglobals // APIKey can input from ldflags

func main() {
	port := flag.Int("port", 7259, "port")
	flag.Parse()

	mainCtx := context.TODO()

	dsn := "publish_manager.db"
	if s, ok := os.LookupEnv("DSN"); ok {
		dsn = s
	}

	publishDB, err := db.NewDB(dsn)
	if err != nil {
		panic(err)
	}
	defer publishDB.Close()

	commentBotDSN := "comment_bot.db"
	if s, ok := os.LookupEnv("COMMENT_BOT_DSN"); ok {
		commentBotDSN = s
	}

	commentBotDB, err := db.NewDB(commentBotDSN)
	if err != nil {
		panic(err)
	}
	defer commentBotDB.Close()

	siteDAO, err := publishDB.NewSiteDAO()
	if err != nil {
		panic(err)
	}

	articleCacheDAO, err := publishDB.NewArticleCacheDAO()
	if err != nil {
		panic(err)
	}

	commentUserDAO, err := commentBotDB.NewCommentUserDAO()
	if err != nil {
		panic(err)
	}

	if e, ok := os.LookupEnv("API_KEY"); ok {
		APIKey = e
	}

	if APIKey == "" {
		panic("api key is not set")
	}
	// Access your API key as an environment variable (see "Set up your API key" above)
	ai, err := aiassist.NewAIAssist(mainCtx, APIKey)
	if err != nil {
		panic(err)
	}
	defer ai.Close()

	zAPI := zBlogApi.NewZBlogAPI()
	publisher := publishManager.NewPublishManager(zAPI, publishManager.DAO{ArticleCacheDAOInterface: articleCacheDAO, SiteDAOInterface: siteDAO})
	siteManager := sitemanager.NewSiteManager(zAPI, siteDAO)

	err = publisher.StartRandomCyclePublish(mainCtx)
	if err != nil {
		panic(err)
	}

	publisher.StartPublishByLack(mainCtx)

	commentBot := commentbot.NewCommentBot(zAPI, siteDAO, commentUserDAO, ai)
	commentBot.StartCycleComment(mainCtx)

	publishHandler := handler.NewPublishHandler(publisher)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/publish", publishHandler.AveragePublishHandler)
	r.POST("/prepublish", publishHandler.PrePublishHandler)
	r.POST("/flexiblePublish", publishHandler.FlexiblePublishHandler)

	siteHandler := handler.NewSiteHandler(siteManager)

	r.POST("/site", siteHandler.AddSiteHandler)
	r.DELETE("/site/:siteID", siteHandler.DeleteSiteHandler)
	r.GET("/site", siteHandler.ListSitesHandler)
	r.GET("/site/:siteID", siteHandler.GetSiteHandler)
	r.PUT("/site", siteHandler.UpdateSiteHandler)
	r.POST("/site/increase_lack", siteHandler.IncreaseLackCountHandler)

	rewriteHandler := handler.NewRewriteHandler(ai)

	r.POST("/rewrite", rewriteHandler.RewriteHandler)

	err = r.Run(fmt.Sprintf(":%d", *port))
	if err != nil {
		panic(err)
	}
}
