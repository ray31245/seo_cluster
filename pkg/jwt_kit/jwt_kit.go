package jwtkit

import (
	"log"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	jwt_kit_interface "github.com/ray31245/seo_cluster/pkg/jwt_kit/jwt_kit_interface"
)

// assert the JWTKit struct implements the SetJWTKit interface
var _ jwt_kit_interface.SetJWTKit = &JWTKit{}

type JWTKit struct {
	secretKey   []byte
	expireTime  time.Duration
	refreshTime time.Duration
	identityKey string
	middleware  *jwt.GinJWTMiddleware
}

func NewJWTKit(
	secretKey []byte,
	expireTime time.Duration,
	refreshTime time.Duration,
	identityKey string,
	payloadFunc func(data interface{}) map[string]interface{},
	identityHandler func(c *gin.Context) interface{},
	authenticator func(c *gin.Context) (interface{}, error),
	authorizator func(data interface{}, c *gin.Context) bool,
	unauthorized func(c *gin.Context, code int, message string),
) *JWTKit {
	res := &JWTKit{
		secretKey:   secretKey,
		expireTime:  expireTime,
		refreshTime: refreshTime,
		identityKey: identityKey,
		middleware: &jwt.GinJWTMiddleware{
			Key:         secretKey,
			Timeout:     expireTime,
			MaxRefresh:  refreshTime,
			IdentityKey: identityKey,

			IdentityHandler: identityHandler,
			Authenticator:   authenticator,
			Authorizator:    authorizator,
			Unauthorized:    unauthorized,
			TokenLookup:     "header: Authorization, query: token, cookie: jwt",
			TokenHeadName:   "Bearer",
			TimeFunc:        time.Now,
		},
	}

	res.SetPayloadFunc(payloadFunc)

	return res
}

func (j *JWTKit) SetPayloadFunc(payloadFunc func(data interface{}) map[string]interface{}) {
	j.middleware.PayloadFunc = func(data interface{}) jwt.MapClaims {
		return jwt.MapClaims(payloadFunc(data))
	}
}

func (j *JWTKit) SetIdentityHandler(identityHandler func(c *gin.Context) interface{}) {
	j.middleware.IdentityHandler = identityHandler
}

func (j *JWTKit) SetAuthenticator(authenticator func(c *gin.Context) (interface{}, error)) {
	j.middleware.Authenticator = authenticator
}

func (j *JWTKit) SetAuthorizator(authorizator func(data interface{}, c *gin.Context) bool) {
	j.middleware.Authorizator = authorizator
}

func (j *JWTKit) SetUnauthorized(unauthorized func(c *gin.Context, code int, message string)) {
	j.middleware.Unauthorized = unauthorized
}

func (j *JWTKit) InitMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errInit := j.middleware.MiddlewareInit()
		if errInit != nil {
			log.Printf("JWTKit InitMiddleWare: %v", errInit)
			ctx.Abort()

			return
		}
	}
}

func (j *JWTKit) MiddlewareFunc() gin.HandlerFunc {
	return j.middleware.MiddlewareFunc()
}

func (j *JWTKit) LoginHandler(c *gin.Context) {
	j.middleware.LoginHandler(c)
}

func (j *JWTKit) RefreshHandler(c *gin.Context) {
	j.middleware.RefreshHandler(c)
}

func (j *JWTKit) LogoutHandler(c *gin.Context) {
	j.middleware.LogoutHandler(c)
}
