package auth

import (
	"fmt"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ray31245/seo_cluster/pkg/auth/model"
	dbInterface "github.com/ray31245/seo_cluster/pkg/db/db_interface"
	dbModel "github.com/ray31245/seo_cluster/pkg/db/model"
	jwtkitInterface "github.com/ray31245/seo_cluster/pkg/jwt_kit/jwt_kit_interface"
	"golang.org/x/crypto/bcrypt"
)

const (
	// claim key for user id
	ClaimKeyUserID = "userID"
)

type Auth struct {
	userDAO dbInterface.UserDAOInterface
}

func NewAuth(userDAO dbInterface.UserDAOInterface) *Auth {
	return &Auth{
		userDAO: userDAO,
	}
}

func (a *Auth) SetUpJWTKit(setJWT jwtkitInterface.SetJWTKit) {
	payloadFunc := func(data interface{}) map[string]interface{} {
		user, ok := data.(*dbModel.User)
		if !ok {
			return map[string]interface{}{}
		}

		return map[string]interface{}{
			ClaimKeyUserID: user.ID,
		}
	}
	setJWT.SetPayloadFunc(payloadFunc)

	identityHandler := func(c *gin.Context) interface{} {
		claims := jwt.ExtractClaims(c)

		userID, ok := claims[ClaimKeyUserID]
		if !ok {
			return nil
		}

		userIDStr, ok := userID.(string)
		if !ok {
			return nil
		}

		userUUID := uuid.MustParse(userIDStr)
		if userUUID == uuid.Nil {
			return nil
		}

		return &dbModel.User{
			Base: dbModel.Base{
				ID: userUUID,
			},
		}
	}
	setJWT.SetIdentityHandler(identityHandler)

	authenticator := func(c *gin.Context) (interface{}, error) {
		var login model.LoginRequest

		if err := c.ShouldBind(&login); err != nil {
			return nil, jwt.ErrMissingLoginValues
		}

		user, err := a.userDAO.GetUserByName(login.UserName)
		if err != nil {
			return nil, jwt.ErrFailedAuthentication
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(login.Password)); err != nil {
			return nil, jwt.ErrFailedAuthentication
		}

		return user, nil
	}
	setJWT.SetAuthenticator(authenticator)

	authorizator := func(data interface{}, c *gin.Context) bool {
		if data == nil {
			return false
		}

		if user, ok := data.(*dbModel.User); !ok || user.ID == uuid.Nil {
			return false
		}

		return true
	}
	setJWT.SetAuthorizator(authorizator)
}

func (a *Auth) GenerateUser(userName, password string) (*dbModel.User, error) {
	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("GenerateUser: %w", err)
	}

	user := &dbModel.User{
		Name:         userName,
		HashPassword: string(hashedPassword),
	}

	return user, nil
}
