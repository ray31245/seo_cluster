package main

import (
	"fmt"

	"github.com/ray31245/seo_cluster/pkg/db"
	"github.com/ray31245/seo_cluster/pkg/db/model"
)

func main() {
	userlist := []model.CommentUser{}

	db, err := db.NewDB("comment_bot.db")
	if err != nil {
		panic(err)
	}

	userDAO, err := db.NewCommentUserDAO()
	if err != nil {
		panic(err)
	}

	for i, user := range userlist {
		user.Name = fmt.Sprintf("comment_user_%d", i)

		if _, err := userDAO.CreateCommentUser(&user); err != nil {
			panic(err)
		}
	}
}
