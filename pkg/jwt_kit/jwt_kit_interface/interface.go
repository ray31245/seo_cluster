package jwtkitinterface

import (
	"github.com/gin-gonic/gin"
)

type SetJWTKit interface {
	SetPayloadFunc(payloadFunc func(data interface{}) map[string]interface{})
	SetIdentityHandler(identityHandler func(c *gin.Context) interface{})
	SetAuthenticator(authenticator func(c *gin.Context) (interface{}, error))
	SetAuthorizator(authorizator func(data interface{}, c *gin.Context) bool)
	SetUnauthorized(unauthorized func(c *gin.Context, code int, message string))
}
