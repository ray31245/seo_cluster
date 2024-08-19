package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	zBlogApi "github.com/ray31245/seo_cluster/pkg/z_blog_api"
	zModel "github.com/ray31245/seo_cluster/pkg/z_blog_api/model"
)

type PostArticleRequest struct {
	Title   string `json:"Title"`
	Content string `json:"Content"`
}

//nolint:gochecknoglobals
var (
	URL      string
	UserName string
	Password string
)

const (
	requestTimeout = 5 * time.Second
)

func main() {
	if v, ok := os.LookupEnv("URL"); ok {
		URL = v
	}

	if v, ok := os.LookupEnv("LOGIN_USERNAME"); ok {
		UserName = v
	}

	if v, ok := os.LookupEnv("LOGIN_PASSWORD"); ok {
		Password = v
	}

	mainCtx := context.TODO()

	api, err := zBlogApi.NewClient(mainCtx, URL, UserName, Password)
	if err != nil {
		log.Fatalln(err)
	}

	port := 7259
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		ReadHeaderTimeout: requestTimeout,
	}

	server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("request route: ", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			http.Error(w, "error", http.StatusInternalServerError)

			return
		}

		bodyData := PostArticleRequest{}

		err = json.Unmarshal(body, &bodyData)
		if err != nil {
			log.Println(err)
			http.Error(w, "error", http.StatusBadRequest)

			return
		}

		art := zModel.PostArticleRequest{
			Title:   bodyData.Title,
			Content: bodyData.Content,
		}

		err = api.PostArticle(mainCtx, art)
		if err != nil {
			log.Println(err)
			http.Error(w, "error", http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)

		_, err = w.Write([]byte("ok"))
		if err != nil {
			log.Println(err)
		}
	})

	log.Printf("Server is running on port %d\n", port)
	// listen on port specified port
	err = server.ListenAndServe()
	if err != nil {
		log.Println("Error: ", err)
	}
}
