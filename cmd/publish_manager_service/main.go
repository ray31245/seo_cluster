package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ray31245/seo_cluster/cmd/publish_manager_service/handler"
	aiassist "github.com/ray31245/seo_cluster/pkg/ai_assist"
	"github.com/ray31245/seo_cluster/pkg/auth"
	"github.com/ray31245/seo_cluster/pkg/db"
	jwt_kit "github.com/ray31245/seo_cluster/pkg/jwt_kit"
	util "github.com/ray31245/seo_cluster/pkg/util"
	zBlogApi "github.com/ray31245/seo_cluster/pkg/z_blog_api"
	commentbot "github.com/ray31245/seo_cluster/service/comment_bot"
	publishManager "github.com/ray31245/seo_cluster/service/publish_manager"
	sitemanager "github.com/ray31245/seo_cluster/service/site_manager"
	usermanager "github.com/ray31245/seo_cluster/service/user_manager"

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

	userDSN := "user.db"
	if s, ok := os.LookupEnv("USER_DSN"); ok {
		userDSN = s
	}

	userDB, err := db.NewDB(userDSN)
	if err != nil {
		panic(err)
	}
	defer userDB.Close()

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

	userDAO, err := userDB.NewUserDAO()
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

	secret := util.GenerateRandomString(32)
	jwtKit := jwt_kit.NewJWTKit([]byte(secret), time.Hour, time.Hour, "loginInfo", nil, nil, nil, nil, nil)
	auth := auth.NewAuth(userDAO)

	auth.SetUpJWTKit(jwtKit)

	publisher := publishManager.NewPublishManager(zAPI, publishManager.DAO{ArticleCacheDAOInterface: articleCacheDAO, SiteDAOInterface: siteDAO})
	siteManager := sitemanager.NewSiteManager(zAPI, siteDAO)
	userManager := usermanager.NewUserManager(userDAO, auth)

	err = publisher.StartRandomCyclePublish(mainCtx)
	if err != nil {
		panic(err)
	}

	publisher.StartPublishByLack(mainCtx)

	commentBot := commentbot.NewCommentBot(zAPI, siteDAO, commentUserDAO, ai)
	commentBot.StartCycleComment(mainCtx)

	r := gin.Default()
	r.Use(gin.Recovery())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/login", jwtKit.LoginHandler)
	r.POST("/refresh_token", jwtKit.RefreshHandler)

	userHandler := handler.NewUserHandler(userManager)

	r.POST("/first_user", userHandler.AddFirstAdminUser)

	publishHandler := handler.NewPublishHandler(publisher)
	rewriteHandler := handler.NewRewriteHandler(ai)

	r.Use(jwtKit.InitMiddleWare())
	r.Use(jwtKit.MiddlewareFunc())

	articleRoute := r.Group("/article")
	articleRoute.POST("/publish", publishHandler.AveragePublishHandler)
	articleRoute.POST("/prepublish", publishHandler.PrePublishHandler)
	articleRoute.POST("/flexiblePublish", publishHandler.FlexiblePublishHandler)
	articleRoute.POST("/rewrite", rewriteHandler.RewriteHandler)
	articleRoute.GET("/cacheCount", publishHandler.GetArticleCacheCountHandler)

	siteHandler := handler.NewSiteHandler(siteManager)

	siteRoute := r.Group("/site")
	siteRoute.POST("/", siteHandler.AddSiteHandler)
	siteRoute.GET("/", siteHandler.ListSitesHandler)
	siteRoute.PUT("/", siteHandler.UpdateSiteHandler)
	siteRoute.DELETE("/:siteID", siteHandler.DeleteSiteHandler)
	siteRoute.GET("/:siteID", siteHandler.GetSiteHandler)
	siteRoute.POST("/increase_lack", siteHandler.IncreaseLackCountHandler)

	err = r.Run(fmt.Sprintf(":%d", *port))
	if err != nil {
		panic(err)
	}
}
