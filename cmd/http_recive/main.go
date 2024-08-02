package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	// recive the request, any route, any method
	// then print the request method and the request URL
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
	})

	port := 8080
	fmt.Printf("Server is running on port %d\n", port)
	// listen on port specified port
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}
