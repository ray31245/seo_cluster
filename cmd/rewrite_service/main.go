package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	aiassist "github.com/ray31245/seo_cluster/pkg/ai_assist"
	"github.com/ray31245/seo_cluster/pkg/util"
)

var APIKey string //nolint:gochecknoglobals // APIKey can input from ldflags

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	ctx := context.Background()

	if e, ok := os.LookupEnv("API_KEY"); ok {
		APIKey = e
	}

	if APIKey == "" {
		log.Fatal("api key is not set")
	}
	// Access your API key as an environment variable (see "Set up your API key" above)
	ai, err := aiassist.NewAIAssist(ctx, APIKey)
	if err != nil {
		log.Fatal(err)
	}
	defer ai.Close()

	port := 9527
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		ReadHeaderTimeout: 5 * time.Second, //nolint:mnd
	}
	// recive the request, any route, any method
	// then print the request method and the request URL
	server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// debugPrint(r)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)

			return
		}

		art, err := ai.Rewrite(ctx, string(body))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)

			return
		}

		art.Content = string(util.MdToHTML([]byte(art.Content)))
		log.Printf("%s", art.Content)

		result, err := util.EscapeHTMLMarshal(art)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)

			return
		}
		// log.Printf("%s", result)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", result)
	})

	log.Printf("Server is running on port %d\n", port)
	// listen on port specified port
	err = server.ListenAndServe()
	if err != nil {
		log.Println("Error: ", err)
	}
}
