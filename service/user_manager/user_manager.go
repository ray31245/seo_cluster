package usermanager

import (
	"context"
	"errors"
	"fmt"

	authInterface "github.com/ray31245/seo_cluster/pkg/auth/auth_interface"
	dbInterface "github.com/ray31245/seo_cluster/pkg/db/db_interface"
)

var ErrUserAlreadyExist = errors.New("user already exist")

type UserManager struct {
	userDAO dbInterface.UserDAOInterface
	auth    authInterface.AuthInterface
}

func NewUserManager(userDAO dbInterface.UserDAOInterface, auth authInterface.AuthInterface) *UserManager {
	return &UserManager{
		userDAO: userDAO,
		auth:    auth,
	}
}

func (u *UserManager) CreateFirstAdminUser(ctx context.Context, userName string, password string) error {
	// check if there is any user
	count, err := u.userDAO.Count()
	if err != nil {
		return fmt.Errorf("CreateFirstAdminUser: %w", err)
	}

	if count > 0 {
		return ErrUserAlreadyExist
	}

	// create admin user
	user, err := u.auth.GenerateUser(userName, password)
	if err != nil {
		return fmt.Errorf("CreateFirstAdminUser: %w", err)
	}

	user.IsAdmin = true

	if err := u.userDAO.CreateUser(user); err != nil {
		return fmt.Errorf("CreateFirstAdminUser: %w", err)
	}

	return nil
}
