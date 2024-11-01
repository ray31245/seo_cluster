package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"unicode/utf8"

	aiAssist "github.com/ray31245/seo_cluster/pkg/ai_assist"
	model "github.com/ray31245/seo_cluster/pkg/ai_assist/model"
	zBlogApi "github.com/ray31245/seo_cluster/pkg/z_blog_api"
	zModel "github.com/ray31245/seo_cluster/pkg/z_blog_api/model"
)

func main() {
	ctx := context.Background()

	file, err := os.Create("keywords")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	token := os.Getenv("AI_ASSIST_TOKEN")

	aiAssistClient, err := aiAssist.NewAIAssist(ctx, token)
	if err != nil {
		log.Fatal(err)
	}

	client := zBlogApi.NewZBlogAPI().NewAnonymousClient(ctx, "https://www.example.com/")

	count, err := client.GetCountOfArticle(ctx, zModel.ListArticleRequest{})
	if err != nil {
		log.Fatal(err)
	}

	repeatMap := make(map[string]bool)

	for count > 0 {
		arts, err := client.ListArticle(ctx, zModel.ListArticleRequest{})
		if err != nil {
			log.Println(err)
		}

		count -= len(arts)

		for _, art := range arts {
			var keywords model.FindKeyWordsResponse

			retryMax := 300

			for {
				keywords, err = aiAssistClient.FindKeyWords(ctx, []byte(art.Content))
				if err != nil && retryMax > 0 {
					log.Println(err)

					retryMax--

					time.Sleep(1 * time.Second)
				} else {
					break
				}
			}

			log.Println("Keywords: ", keywords)

			for _, w := range keywords.KeyWords {
				if utf8.RuneCountInString(w) > 4 && utf8.RuneCountInString(w) < 8 && !repeatMap[w] {
					file.WriteString(fmt.Sprintf("%s\n", w))
					repeatMap[w] = true
				}
			}
		}
	}
}
