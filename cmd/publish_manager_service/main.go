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
	wordpressApi "github.com/ray31245/seo_cluster/pkg/wordpress_api"
	zBlogApi "github.com/ray31245/seo_cluster/pkg/z_blog_api"
	articleCacheManager "github.com/ray31245/seo_cluster/service/article_cache_manager"
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

	configDSN := "config.db"
	if s, ok := os.LookupEnv("CONFIG_DSN"); ok {
		configDSN = s
	}

	configBD, err := db.NewDB(configDSN)
	if err != nil {
		panic(err)
	}
	defer configBD.Close()

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

	configDAO, err := configBD.NewKVConfigDAO()
	if err != nil {
		panic(err)
	}

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
	ai, err := aiassist.NewAIAssist(mainCtx, APIKey, false)
	if err != nil {
		panic(err)
	}
	defer ai.Close()

	zAPI := zBlogApi.NewZBlogAPI()
	wordpressAPI := wordpressApi.NewWordpressApi()

	secret := util.GenerateRandomString(32)
	jwtKit := jwt_kit.NewJWTKit([]byte(secret), time.Hour, time.Hour, "loginInfo", nil, nil, nil, nil, nil)
	auth := auth.NewAuth(userDAO)

	auth.SetUpJWTKit(jwtKit)

	publisher := publishManager.NewPublishManager(zAPI, wordpressAPI, publishManager.DAO{ArticleCacheDAOInterface: articleCacheDAO, SiteDAOInterface: siteDAO, KVConfigDAOInterface: configDAO}, ai)
	siteManager := sitemanager.NewSiteManager(zAPI, wordpressAPI, siteDAO)
	userManager := usermanager.NewUserManager(userDAO, auth)
	articleCacheManager := articleCacheManager.NewArticleCacheManager(articleCacheDAO)

	err = publisher.StartRandomCyclePublishZblog(mainCtx)
	if err != nil {
		panic(err)
	}

	err = publisher.StartRandomCyclePublishWordPress(mainCtx)
	if err != nil {
		panic(err)
	}

	err = publisher.StartUpdateArticleTagSignalLoop(mainCtx, 1, 5)
	if err != nil {
		panic(err)
	}

	publisher.StartPublishByLack(mainCtx)

	commentBot := commentbot.NewCommentBot(zAPI, configDAO, siteDAO, commentUserDAO, ai)
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
	rewriteHandler := handler.NewRewriteHandler(ai, configDAO)
	articleCacheHandler := handler.NewArticleCacheHandler(articleCacheManager)

	r.Use(jwtKit.InitMiddleWare())
	r.Use(jwtKit.MiddlewareFunc())

	configRoute := r.Group("/config")
	configRoute.PUT("/set_un_cate_Name", publishHandler.SetConfigUnCategoryNameHandler)
	configRoute.GET("/get_un_cate_Name", publishHandler.GetConfigUnCateNameHandler)
	configRoute.PUT("/set_tag_blacklist", publishHandler.SetConfigTagBlackList)
	configRoute.GET("/get_tag_blacklist", publishHandler.GetConfigTagBlackList)

	articleRoute := r.Group("/article")
	articleRoute.POST("/publish", publishHandler.AveragePublishHandler)
	articleRoute.POST("/prepublish", publishHandler.PrePublishHandler)
	articleRoute.POST("/flexiblePublish", publishHandler.FlexiblePublishHandler)
	articleRoute.POST("/directPublish/:cateID", publishHandler.DirectPublishHandler)
	articleRoute.POST("/broadcastPublish", publishHandler.BroadcastPublishHandler)
	articleRoute.PUT("/stopAutoPublish", publishHandler.StopAutoPublishHandler)
	articleRoute.PUT("/startAutoPublish", publishHandler.StartAutoPublishHandler)
	articleRoute.GET("/stopAutoPublishStatus", publishHandler.GetStopAutoPublishStatusHandler)
	articleRoute.GET("/cacheCount", publishHandler.GetArticleCacheCountHandler)
	articleRoute.GET("/listPublishLaterArticleCache", articleCacheHandler.ListPublishLaterArticleCacheHandler)
	articleRoute.GET("/listEditAbleArticleCache", articleCacheHandler.ListEditAbleArticleCacheHandler)
	articleRoute.PUT("/updateArticleCacheStatus", articleCacheHandler.UpdateArticleCacheStatusHandler)
	articleRoute.PUT("/editArticleCache", articleCacheHandler.EditArticleCacheHandler)
	articleRoute.DELETE("/deleteArticleCache", articleCacheHandler.DeleteArticleCacheHandler)

	articleRewriteRoute := articleRoute.Group("/rewrite")
	articleRewriteRoute.POST("/", rewriteHandler.RewriteHandler)
	articleRewriteRoute.PUT("/set_default_system_prompt", rewriteHandler.SetDefaultSystemPromptHandler)
	articleRewriteRoute.GET("/get_default_system_prompt", rewriteHandler.GetDefaultSystemPromptHandler)
	articleRewriteRoute.PUT("/set_default_prompt", rewriteHandler.SetDefaultPromptHandler)
	articleRewriteRoute.GET("/get_default_prompt", rewriteHandler.GetDefaultPromptHandler)
	articleRewriteRoute.PUT("/set_default_extend_system_prompt", rewriteHandler.SetDefaultExtendSystemPromptHandler)
	articleRewriteRoute.GET("/get_default_extend_system_prompt", rewriteHandler.GetDefaultExtendSystemPromptHandler)
	articleRewriteRoute.PUT("/set_default_extend_prompt", rewriteHandler.SetDefaultExtendPromptHandler)
	articleRewriteRoute.GET("/get_default_extend_prompt", rewriteHandler.GetDefaultExtendPromptHandler)

	siteHandler := handler.NewSiteHandler(siteManager)

	siteRoute := r.Group("/site")
	siteRoute.POST("/", siteHandler.AddSiteHandler)
	siteRoute.GET("/", siteHandler.ListSitesHandler)
	siteRoute.PUT("/", siteHandler.UpdateSiteHandler)
	siteRoute.DELETE("/:siteID", siteHandler.DeleteSiteHandler)
	siteRoute.GET("/:siteID", siteHandler.GetSiteHandler)
	siteRoute.PUT("/syncCateFromSite/:siteID", siteHandler.SyncCategoryFromSiteHandler)
	siteRoute.PUT("/syncCateFromAllSite", siteHandler.SyncCategoryFromAllSiteHandler)
	siteRoute.POST("/increase_lack", siteHandler.IncreaseLackCountHandler)

	commentBotHandler := handler.NewCommentBotHandler(commentBot)

	commentBotRoute := r.Group("/comment_bot")
	commentBotRoute.PUT("/stopAutoComment", commentBotHandler.StopAutoCommentHandler)
	commentBotRoute.PUT("/startAutoComment", commentBotHandler.StartAutoCommentHandler)
	commentBotRoute.GET("/getStopAutoCommentStatus", commentBotHandler.GetStopAutoCommentStatusHandler)

	err = r.Run(fmt.Sprintf(":%d", *port))
	if err != nil {
		panic(err)
	}
}
