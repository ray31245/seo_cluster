package main

import (
	"context"

	"github.com/ray31245/seo_cluster/pkg/db"
	zblogapi "github.com/ray31245/seo_cluster/pkg/z_blog_api"
	zModel "github.com/ray31245/seo_cluster/pkg/z_blog_api/model"
)

func main() {
	db, err := db.NewDB("comment_bot.db")
	if err != nil {
		panic(err)
	}

	userDAO, err := db.NewCommentUserDAO()
	if err != nil {
		panic(err)
	}

	users, err := userDAO.ListCommentUsers()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	client, err := zblogapi.NewClient(ctx, "http://www.test.com/", "admin_bevis_toor", "3cc31cd246149aec68079241e71e98f6")
	if err != nil {
		panic(err)
	}

	for _, user := range users {
		err := client.PostMember(ctx, zModel.PostMemberRequest{
			Member: zModel.Member{
				Level: "5",
				Name:  user.Name,
				Alias: user.Alias,
			},
			Password:   user.Password,
			PasswordRe: user.Password,
		})
		if err != nil {
			panic(err)
		}
	}
}
