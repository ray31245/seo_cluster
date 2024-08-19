package main

import (
	"context"
	"log"
	"time"

	zBlogApi "github.com/ray31245/seo_cluster/pkg/z_blog_api"
	zModel "github.com/ray31245/seo_cluster/pkg/z_blog_api/model"

	"github.com/google/uuid"
)

func main() {
	ctx := context.Background()
	code := 100
	log.Println(code/200 == 1)

	api, err := zBlogApi.NewClient(ctx, "http://www.test.com", "bevis", "3cc31cd246149aec68079241e71e98f6")
	if err != nil {
		log.Fatalln(err)
	}

	// list category
	// list, err := api.ListCategory()
	// if err != nil {
	// 	log.Println(err)
	// }
	// for _, v := range list {
	// 	log.Printf("category id: %s\nname: %s\n", v.ID, v.Name)
	// }

	// // list article
	// list, err := api.ListArticle(zBlogApi.ListArticleRequest{})
	// if err != nil {
	// 	log.Println(err)
	// }
	// for _, v := range list {
	// 	log.Printf("article id: %s\nname: %s\n", v.ID, v.Title)
	// }

	// get count of article
	count, err := api.GetCountOfArticle(ctx, zModel.ListArticleRequest{})
	if err != nil {
		log.Println(err)
	}

	log.Printf("count of article: %d\n", count)

	// // empty article
	// for {
	// 	list, err := api.ListArticle(zModel.ListArticleRequest{})
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// 	if len(list) == 0 {
	// 		break
	// 	}
	// 	for _, v := range list {
	// 		log.Printf("delete article id: %s\nname: %s\n", v.ID, v.Title)
	// 		err := api.DeleteArticle(v.ID)
	// 		if err != nil {
	// 			log.Println(err)
	// 		}
	// 	}
	// }
	watchNumbersOfCategory(ctx)
}

func watchNumbersOfCategory(ctx context.Context) {
	type site struct {
		id       uuid.UUID
		url      string
		userName string
		password string
	}

	a := zBlogApi.NewZBlogAPI()
	sites := []site{
		{id: uuid.New(), url: "http://www.test.com", userName: "bevis", password: "3cc31cd246149aec68079241e71e98f6"},
		{id: uuid.New(), url: "http://www.test2.com", userName: "bevis", password: "3cc31cd246149aec68079241e71e98f6"},
		{id: uuid.New(), url: "http://www.test3.com", userName: "bevis", password: "3cc31cd246149aec68079241e71e98f6"},
		{id: uuid.New(), url: "http://www.test4.com", userName: "bevis", password: "3cc31cd246149aec68079241e71e98f6"},
		{id: uuid.New(), url: "http://www.test5.com", userName: "bevis", password: "3cc31cd246149aec68079241e71e98f6"},
	}

	for {
		counts := make([]int, len(sites))

		for i, v := range sites {
			test, err := a.GetClient(ctx, v.id, v.url, v.userName, v.password)
			if err != nil {
				log.Fatalln(err)
			}

			count, err := test.GetCountOfArticle(ctx, zModel.ListArticleRequest{})
			if err != nil {
				log.Println(err)
			}

			counts[i] = count
		}

		log.Printf("count of 1: %d  2: %d  3: %d  4: %d  5: %d ", counts[0], counts[1], counts[2], counts[3], counts[4])

		time.Sleep(10 * time.Second) //nolint: mnd
	}
}
