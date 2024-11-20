package main

import (
	"context"
	"log"
	"runtime"
	"sync"
	"time"

	zBlogApi "github.com/ray31245/seo_cluster/pkg/z_blog_api"
	zModel "github.com/ray31245/seo_cluster/pkg/z_blog_api/model"

	"github.com/google/uuid"
)

func main() {
	watchNumbersOfCategory()
}

type site struct {
	id       uuid.UUID
	url      string
	userName string
	password string
}

func watchNumbersOfCategory() {
	ctx := context.Background()

	a := zBlogApi.NewZBlogAPI()
	sites := []site{
		{id: uuid.New(), url: "https://www.example.com", userName: "usr", password: "pwd"},
	}

	for {
		var wg *sync.WaitGroup = new(sync.WaitGroup)
		counts := make([]int, len(sites))

		nCPU := runtime.NumCPU()
		reqCh := make(chan recordRequest, nCPU)
		wg.Add(nCPU)
		for range nCPU {
			go recordArticleCount(ctx, wg, a, counts, reqCh)
		}

		for i, v := range sites {
			reqCh <- recordRequest{index: i, siteInfo: v}
		}

		wg.Wait()
		for i, s := range sites {
			log.Printf("site: %s  count: %d\n", s.url, counts[i])
		}

		time.Sleep(600 * time.Second) //nolint: mnd
	}
}

type recordRequest struct {
	index    int
	siteInfo site
}

func recordArticleCount(ctx context.Context, wg *sync.WaitGroup, api *zBlogApi.ZBlogAPI, countSlice []int, reqCh <-chan recordRequest) {
	for {
		select {
		case req := <-reqCh:
			test := api.NewAnonymousClient(ctx, req.siteInfo.url)

			count, err := test.GetCountOfArticle(ctx, zModel.ListArticleRequest{})
			if err != nil {
				log.Println(err)
			}

			countSlice[req.index] = count
		default:
			wg.Done()
			return
		}
	}
}
