package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	authModel "github.com/ray31245/seo_cluster/pkg/auth/model"
)

const (
	// DefaultPort is the default port for the server
	DefaultPort = 8080
	// RefreshTokenInterval is the interval to refresh token
	RefreshTokenInterval = time.Minute * 30
)

func main() {
	target := flag.String("target", "", "target url")
	port := flag.Int("port", DefaultPort, "port")
	userName := flag.String("username", "", "username for login")
	password := flag.String("password", "", "password for login")
	flag.Parse()

	if *target == "" {
		panic("target is required")
	}

	handler, err := handleForwarder(
		*target,
		authModel.LoginRequest{
			UserName: *userName,
			Password: *password,
		},
	)
	if err != nil {
		panic("handler not available: " + err.Error())
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}

func handleForwarder(target string, loginInfo authModel.LoginRequest) (func(http.ResponseWriter, *http.Request), error) {
	client, err := NewStatusClient(target, loginInfo)
	if err != nil {
		return nil, fmt.Errorf("error when create status client: %w", err)
	}

	err = client.Login()
	if err != nil {
		return nil, fmt.Errorf("error when login: %w", err)
	}

	go func() {
		for {
			select {
			case <-time.After(RefreshTokenInterval):
				err := client.RefreshToken()
				if err != nil {
					log.Printf("error when refresh token: %s", err.Error())
				}
			}
		}
	}()

	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("forward request to %s", r.URL.String())

		client.ForwardRequest(w, r)
	}, nil
}

type StatusClient struct {
	c         *http.Client
	token     string
	targetUrl *url.URL
	loginBody []byte
}

func NewStatusClient(target string, login authModel.LoginRequest) (*StatusClient, error) {
	targetUrl, err := url.Parse(target)
	if err != nil {
		return nil, fmt.Errorf("error when parse target: %s", err.Error())
	}

	if targetUrl.Scheme == "" {
		targetUrl.Scheme = "https"
	}

	if targetUrl.Path != "/" {
		log.Printf("path %s force to /", targetUrl.Path)
		targetUrl.Path = "/"
	}

	testClient := http.Client{Timeout: time.Second * 3}

	_, err = testClient.Get(targetUrl.String())
	if err != nil {
		return nil, fmt.Errorf("error when test target: %s", err.Error())
	}

	log.Println("target is reachable")

	loginBody, err := json.Marshal(login)
	if err != nil {
		return nil, fmt.Errorf("error when marshal login: %s", err.Error())
	}

	return &StatusClient{
		c:         &http.Client{},
		targetUrl: targetUrl,
		loginBody: loginBody,
	}, nil
}

func (s *StatusClient) Login() error {
	loginUrl := s.targetUrl
	loginUrl.Path = "/login"

	req, err := http.NewRequest(http.MethodPost, loginUrl.String(), bytes.NewReader(s.loginBody))
	if err != nil {
		return fmt.Errorf("error when create login request: %s", err.Error())
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.c.Do(req)
	if err != nil {
		return fmt.Errorf("error when login: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error when read login body: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error when login: status code %d, msg: %s", resp.StatusCode, string(body))
	}

	err = s.ParseAndSetToken(body)
	if err != nil {
		return fmt.Errorf("error when parse login token: %s", err.Error())
	}

	log.Println("login success")

	return nil
}

func (s *StatusClient) RefreshToken() error {
	req, err := http.NewRequest(http.MethodPost, "/refresh_token", nil)
	if err != nil {
		return fmt.Errorf("error when create refresh token request: %s", err.Error())
	}

	res, err := s.retryIfUnauthorized(req)
	if err != nil {
		return fmt.Errorf("error when refresh token: %s", err.Error())
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error when refresh token: status code %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error when read refresh token body: %s", err.Error())
	}

	err = s.ParseAndSetToken(body)
	if err != nil {
		return fmt.Errorf("error when parse refresh token: %s", err.Error())
	}

	return nil
}

func (s *StatusClient) ParseAndSetToken(body []byte) error {
	var token struct {
		Token string `json:"token"`
	}

	err := json.Unmarshal(body, &token)
	if err != nil {
		return fmt.Errorf("error when parse token: %s", err.Error())
	}

	s.token = token.Token

	return nil
}

func (s *StatusClient) ForwardRequest(w http.ResponseWriter, r *http.Request) {
	// copy request
	proxyReq, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)

		return
	}

	// copy header
	for k, vs := range r.Header {
		for _, v := range vs {
			proxyReq.Header.Add(k, v)
		}
	}

	resp, err := s.retryIfUnauthorized(proxyReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("error forward request: %s", err.Error()), http.StatusBadRequest)

		return
	}
	defer resp.Body.Close()

	for k, vs := range resp.Header {
		for _, v := range vs {
			w.Header().Add(k, v)
		}
	}

	w.WriteHeader(resp.StatusCode)

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Failed to write response: %v", err)
	}

	return
}

func (s *StatusClient) DoRequest(req *http.Request) (*http.Response, error) {
	// over write the url
	req.URL.Scheme = s.targetUrl.Scheme
	req.URL.Host = s.targetUrl.Host

	// set token
	req.Header.Set("Authorization", "Bearer "+s.token)

	return s.c.Do(req)
}

func (s *StatusClient) retryIfUnauthorized(req *http.Request) (*http.Response, error) {
	resp, err := s.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error when do request: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		resp.Body.Close()

		err = s.Login()
		if err != nil {
			return nil, fmt.Errorf("error when login: %w", err)
		}

		resp, err = s.DoRequest(req)
		if err != nil {
			return nil, fmt.Errorf("error when do request: %w", err)
		}
	}

	return resp, nil
}
