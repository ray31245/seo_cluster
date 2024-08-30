package main

import (
	"context"
	"net/http"
	"os"

	"github.com/ray31245/seo_cluster/cmd/publish_manager_service/handler"
	"github.com/ray31245/seo_cluster/pkg/db"
	zBlogApi "github.com/ray31245/seo_cluster/pkg/z_blog_api"
	publishManager "github.com/ray31245/seo_cluster/service/publish_manager"

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

	mainCtx := context.TODO()

	err = publisher.StartRandomCyclePublish(mainCtx)
	if err != nil {
		panic(err)
	}

	handler := handler.NewHandler(publisher)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/publish", handler.AveragePublishHandler)
	r.POST("/site", handler.AddSiteHandler)
	r.GET("/site", handler.ListSitesHandler)
	r.POST("/prepublish", handler.PrePublishHandler)
	r.POST("/flexiblePublish", handler.FlexiblePublishHandler)

	err = r.Run(":7259")
	if err != nil {
		panic(err)
	}
}
