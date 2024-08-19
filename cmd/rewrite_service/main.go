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

	"goTool/pkg/util"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Article struct {
	Title   string `json:"Title"`
	Content string `json:"Content"`
}

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
	client, err := genai.NewClient(ctx, option.WithAPIKey(APIKey))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// The Gemini 1.5 models are versatile and work with both text-only and multimodal prompts
	model := client.GenerativeModel("gemini-1.5-flash")
	model.GenerationConfig = genai.GenerationConfig{
		ResponseMIMEType: "application/json",
	}

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

		//nolint:gosmopolitan // prompt is a string
		prompt := "你是一位收悉區塊鏈的專欄作家，請你將以下內容用你的話重新闡述文章中的內容，並訂一個標題。請使用json格式輸出：{Title: string,Content: string}"

		resp, err := model.GenerateContent(ctx, genai.Text(fmt.Sprintf("%s\n%s", prompt, body)))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)

			return
		}

		if len(resp.Candidates) < 1 {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("no candidate")

			return
		}
		// log.Println(len(resp.Candidates))
		art := Article{}

		err = json.Unmarshal([]byte(fmt.Sprintf("%s", resp.Candidates[0].Content.Parts[0])), &art)
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
