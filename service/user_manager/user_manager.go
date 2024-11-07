package usermanager

import (
	"context"
	"fmt"

	dbInterface "github.com/ray31245/seo_cluster/pkg/db/db_interface"
	dbModel "github.com/ray31245/seo_cluster/pkg/db/model"
	"golang.org/x/crypto/bcrypt"
)

type UserManager struct {
	userDAO dbInterface.UserDAOInterface
}

func NewUserManager(userDAO dbInterface.UserDAOInterface) *UserManager {
	return &UserManager{
		userDAO: userDAO,
	}
}

func (u *UserManager) CreateFirstAdminUser(ctx context.Context, userName string, password string) error {
	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("CreateFirstAdminUser: %w", err)
	}

	user := &dbModel.User{
		Name:         userName,
		HashPassword: string(hashedPassword),
		IsAdmin:      true,
	}

	if err := u.userDAO.CreateFirstAdminUser(user); err != nil {
		return fmt.Errorf("CreateFirstAdminUser: %w", err)
	}

	return nil
}
