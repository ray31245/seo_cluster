package main

import (
	"context"
	"net/http"
	"os"

	"github.com/ray31245/seo_cluster/cmd/publish_manager_service/handler"
	"github.com/ray31245/seo_cluster/pkg/db"
	zBlogApi "github.com/ray31245/seo_cluster/pkg/z_blog_api"
	publishManager "github.com/ray31245/seo_cluster/service/publish_manager"
	sitemanager "github.com/ray31245/seo_cluster/service/site_manager"

	"github.com/gin-gonic/gin"
)

func main() {
	dsn := "publish_manager.db"
	if s, ok := os.LookupEnv("DSN"); ok {
		dsn = s
	}

	db, err := db.NewDB(dsn)
	if err != nil {
		panic(err)
	}

	siteDAO, err := db.NewSiteDAO()
	if err != nil {
		panic(err)
	}

	articleCacheDAO, err := db.NewArticleCacheDAO()
	if err != nil {
		panic(err)
	}

	zAPI := zBlogApi.NewZBlogAPI()
	publisher := publishManager.NewPublishManager(zAPI, publishManager.DAO{ArticleCacheDAOInterface: articleCacheDAO, SiteDAOInterface: siteDAO})
	siteManager := sitemanager.NewSiteManager(zAPI, siteDAO)

	mainCtx := context.TODO()

	err = publisher.StartRandomCyclePublish(mainCtx)
	if err != nil {
		panic(err)
	}

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
	r.GET("/site", siteHandler.ListSitesHandler)

	err = r.Run(":7259")
	if err != nil {
		panic(err)
	}
}
