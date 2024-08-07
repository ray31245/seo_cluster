package main

import (
	"encoding/json"
	"fmt"
	zblogapi "goTool/pkg/z_blog_api"
	"io"
	"log"
	"net/http"
)

type PostArticleRequest struct {
	Title   string `json:"Title"`
	Content string `json:"Content"`
}

func main() {
	api := zblogapi.NewZblogAPI("http://www.test2.com/zb_system/api.php", "bevis", "3cc31cd246149aec68079241e71e98f6")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("request route: ", r.URL.Path)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error"))
			log.Println(err)
			return
		}
		bodyData := PostArticleRequest{}
		err = json.Unmarshal(body, &bodyData)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error"))
			log.Println(err)
			return
		}
		art := zblogapi.PostArticleRequest{
			Title:   bodyData.Title,
			Content: bodyData.Content,
		}
		err = api.PostArticle(art)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	port := 7259
	fmt.Printf("Server is running on port %d\n", port)
	// listen on port specified port
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	// members, err := api.ListMember()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(members)
}
