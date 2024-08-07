package zblogapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
)

type ZblogAPI struct {
	baseURL  url.URL
	token    string
	userName string
	password string
	// avoid duplicate login
	lock *sync.Mutex
}

func NewZblogAPI(urlStr string, userName string, password string) *ZblogAPI {
	baseURL, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}
	return &ZblogAPI{
		baseURL:  *baseURL,
		lock:     &sync.Mutex{},
		userName: userName,
		password: password,
		token:    "YmV2aXN8fHw3ZWYxMmJkNTQ1ZmU5MTRhNTMwYTFlYjMyODUxYTA5YTg4YjE0OGRmYjExN2Y2ODRkZmZmNzM1ZjM2YTcwMmI4MTcyMjQxODY0OA==",
	}
}

func (t *ZblogAPI) Login() error {
	t.lock.Lock()
	defer t.lock.Unlock()
	values := t.baseURL.Query()
	values.Add("mod", "member")
	values.Add("act", "login")
	t.baseURL.RawQuery = values.Encode()
	data := map[string]interface{}{}
	data["username"] = t.userName
	data["password"] = t.password
	bytesData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, t.baseURL.String(), bytes.NewReader(bytesData))
	if err != nil {
		return err
	}
	// req.Header.Add("Authorization", "Bearer "+t.token)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("status code error: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read body error: %w", err)
	}
	var loginRes LoginResponse
	if err := json.Unmarshal(body, &loginRes); err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}
	if loginRes.Code != 200 {
		return fmt.Errorf("login error: %s", loginRes.Message)
	}
	t.token = loginRes.Data.Token
	return nil
}

func (t *ZblogAPI) ListMember() (string, error) {
	res := ""
	var err error
	task := func() error {
		res, err = t.listMember()
		return err
	}
	err = t.retry(task)
	return res, err
}

func (t *ZblogAPI) PostArticle(art PostArticleRequest) error {
	var err error
	task := func() error {
		err = t.postArticle(art)
		return err
	}
	err = t.retry(task)
	return err
}

func (t *ZblogAPI) ListArticle(req ListArticleRequest) ([]Article, error) {
	res := ListArticleResponse{}
	var err error
	task := func() error {
		res, err = t.listArticle(req)
		return err
	}
	err = t.retry(task)
	return res.Data.List, err
}
