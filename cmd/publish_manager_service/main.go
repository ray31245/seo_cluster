package main

import (
	"goTool/cmd/publish_manager_service/handler"
	"goTool/pkg/db"
	publishmanager "goTool/pkg/publish_manager"
	zblogapi "goTool/pkg/z_blog_api"
	"net/http"
	"os"

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
	err = db.Migrate()
	if err != nil {
		panic(err)
	}

	dao := db.NewDAO()
	zApi := zblogapi.NewZblogAPI()
	publisher := publishmanager.NewPublishManager(zApi, dao)

	err = publisher.StartRandomCyclePublish()
	if err != nil {
		panic(err)
	}

	handler := handler.NewHandler(dao, zApi, publisher)

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
