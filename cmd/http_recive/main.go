package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	port := 8080
	server := http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		ReadHeaderTimeout: 5 * time.Second, //nolint:mnd
	}

	// receive the request, any route, any method
	// then print the request method and the request URL
	server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { //nolint:revive
		log.Println("--------------------")
		// print the request method and the request URL
		log.Println("Request Route: ", r.URL.Path)
		// the request method is the method
		log.Println("Request Method: ", r.Method)
		// the request parameter is the parameter
		log.Println("Request Parameter: ", r.URL.Query())
		// the request body is the body
		log.Println("Request Body: ", r.Body)

		s, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("Error: ", err)
		} else {
			// log.Println("Request Body: ", string(s))
			body := make(map[string]interface{})

			err = json.Unmarshal(s, &body)
			if err != nil {
				log.Println("Request Body: ", string(s))
			}

			for key, value := range body {
				log.Println("Key: ", key, " Value: ", value)
			}
		}
		// the request header is the header
		log.Println("Request Header: ", r.Header)
	})

	log.Printf("Server is running on port %d\n", port)
	// listen on port specified port
	err := server.ListenAndServe()
	if err != nil {
		log.Println("Error: ", err)
	}
}
