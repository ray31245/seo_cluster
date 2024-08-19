package main

import (
	"context"
	"net/http"
	"os"

	"goTool/cmd/publish_manager_service/handler"
	"goTool/pkg/db"
	zBlogApi "goTool/pkg/z_blog_api"
	publishManager "goTool/service/publish_manager"

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
	r.POST("/prepublish", handler.PrePublishHandler)
	r.POST("/flexiblePublish", handler.FlexiblePublishHandler)

	err = r.Run(":7259")
	if err != nil {
		panic(err)
	}
}
