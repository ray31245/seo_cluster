package main

import (
	// "goTools/pkg/util"
	"context"
	"encoding/json"
	"fmt"
	"goTool/pkg/util"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type Article struct {
	Title   string `json:"Title"`
	Content string `json:"Content"`
}

var API_KEY string

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	ctx := context.Background()
	if e, ok := os.LookupEnv("API_KEY"); ok {
		API_KEY = e
	}

	if API_KEY == "" {
		log.Fatal("api key is not set")
	}
	// Access your API key as an environment variable (see "Set up your API key" above)
	client, err := genai.NewClient(ctx, option.WithAPIKey(API_KEY))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// The Gemini 1.5 models are versatile and work with both text-only and multimodal prompts
	model := client.GenerativeModel("gemini-1.5-flash")
	model.GenerationConfig = genai.GenerationConfig{
		ResponseMIMEType: "application/json",
	}

	// recive the request, any route, any method
	// then print the request method and the request URL
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// debugPrint(r)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
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
		art.Content = string(mdToHTML([]byte(art.Content)))
		log.Printf("%s", art.Content)
		result, err := util.EscapeHTMLMarshual(art)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
		// log.Printf("%s", result)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", result)
	})

	port := 9527
	fmt.Printf("Server is running on port %d\n", port)
	// listen on port specified port
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extenstions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extenstions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)
	return markdown.Render(doc, renderer)
}

func debugPrint(r *http.Request) {
	fmt.Println("--------------------")
	// print the request method and the request URL
	fmt.Println("Request Route: ", r.URL.Path)
	// the request method is the method
	fmt.Println("Request Method: ", r.Method)
	// the request parameter is the parameter
	fmt.Println("Request Parameter: ", r.URL.Query())
	// the request body is the body
	fmt.Println("Request Body: ", r.Body)
	s, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		// fmt.Println("Request Body: ", string(s))
		body := make(map[string]interface{})
		err = json.Unmarshal(s, &body)
		if err != nil {
			fmt.Println("Request Body: ", string(s))
		}
		for key, value := range body {
			fmt.Println("Key: ", key, " Value: ", value)
		}
	}
	// the request header is the header
	fmt.Println("Request Header: ", r.Header)
}
