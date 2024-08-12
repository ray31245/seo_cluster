package main

import (
	"goTool/cmd/publish_manager_service/handler"
	"goTool/pkg/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := db.NewDB("publish_manager.db")
	if err != nil {
		panic(err)
	}
	err = db.Migrate()
	if err != nil {
		panic(err)
	}

	dao := db.NewDAO()

	handler := handler.NewHandler(dao)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/publish", handler.AveragePublishHandler)
	r.POST("/site", handler.AddSiteHandler)

	err = r.Run(":7259")
	if err != nil {
		panic(err)
	}
}
