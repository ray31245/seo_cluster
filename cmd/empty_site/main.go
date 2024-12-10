package main

import (
	"context"
	"fmt"
	"log"

	zBlogApi "github.com/ray31245/seo_cluster/pkg/z_blog_api"
	zModel "github.com/ray31245/seo_cluster/pkg/z_blog_api/model"
)

func main() {
	// print hint to request user to input URL and user name, password
	fmt.Print("Please input the URL: ")

	var url string

	fmt.Scanln(&url)
	fmt.Print("Please input the user name: ")

	var userName string

	fmt.Scanln(&userName)
	fmt.Print("Please input the password: ")

	var password string

	fmt.Scanln(&password)

	// print the input URL and user name, password
	fmt.Println("following url is used to login")
	fmt.Println(url)
	fmt.Println("following username is used to login")
	fmt.Println(userName)
	fmt.Println("following password is used to login")
	fmt.Println(password)

	// confirm that the user really want to do this
	var confirm string

	for {
		fmt.Print("Do you want to login? (y/n):")
		fmt.Scanln(&confirm)

		if confirm == "n" {
			fmt.Println("You have canceled")

			return
		} else if confirm == "y" {
			fmt.Println("You have confirmed")

			break
		} else {
			fmt.Println("Invalid input")
		}
	}

	fmt.Println("process running...")

	ctx := context.Background()

	api, err := zBlogApi.NewClient(ctx, url, userName, password)
	if err != nil {
		log.Fatalln(err)
	}

	total := 0

	for {
		list, err := api.ListArticle(context.Background(), zModel.ListArticleRequest{})
		if err != nil {
			log.Println(err)
		}

		if len(list) == 0 {
			break
		}

		for _, v := range list {
			log.Printf("delete article id: %s\nname: %s\n", v.ID, v.Title)

			if err := api.DeleteArticle(context.Background(), string(v.ID)); err != nil {
				log.Println(err)
			}

			total++
		}
	}

	log.Printf("total deleted article: %d\n", total)
}
