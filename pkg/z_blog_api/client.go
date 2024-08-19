package zblogapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/ray31245/seo_cluster/pkg/z_blog_api/model"
)

type Client struct {
	baseURL  url.URL
	token    string
	userName string
	password string
	// avoid duplicate login
	lock *sync.Mutex
}

func NewClient(ctx context.Context, urlStr string, userName string, password string) (*Client, error) {
	log.Println("following url is used to login")
	log.Println(urlStr)
	log.Println("following username is used to login")
	log.Println(userName)
	log.Println("following password is used to login")
	log.Println(password)

	baseURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("parse url error: %w", err)
	}
	// add default path of z-blog api
	baseURL = baseURL.JoinPath("zb_system/api.php")
	res := &Client{
		baseURL:  *baseURL,
		lock:     &sync.Mutex{},
		userName: userName,
		password: password,
		token:    "YmV2aXN8fHw3ZWYxMmJkNTQ1ZmU5MTRhNTMwYTFlYjMyODUxYTA5YTg4YjE0OGRmYjExN2Y2ODRkZmZmNzM1ZjM2YTcwMmI4MTcyMjQxODY0OA==",
	}

	err = res.Login(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *Client) Login(ctx context.Context) error {
	data := map[string]interface{}{}
	data["username"] = t.userName
	data["password"] = t.password

	bytesData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	resBody, err := t.requestWithBlock(ctx, http.MethodPost, map[string]interface{}{"mod": "member", "act": "login"}, bytesData)
	if err != nil {
		return fmt.Errorf("login error: %w", err)
	}

	var loginRes model.LoginResponse
	if err := json.Unmarshal(resBody, &loginRes); err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}

	t.token = loginRes.Data.Token

	log.Println("login success")

	return nil
}

func (t *Client) ListMember(ctx context.Context) (model.ListMemberResponse, error) {
	res := model.ListMemberResponse{}

	var err error

	task := func() error {
		res, err = t.listMember(ctx)

		return err
	}
	err = t.retry(ctx, task)

	return res, err
}

func (t *Client) PostArticle(ctx context.Context, art model.PostArticleRequest) error {
	var err error

	task := func() error {
		err = t.postArticle(ctx, art)
		if err != nil {
			return fmt.Errorf("PostArticle: %w", err)
		}

		return nil
	}
	err = t.retry(ctx, task)

	return err
}

func (t *Client) ListArticle(ctx context.Context, req model.ListArticleRequest) ([]model.Article, error) {
	res := model.ListArticleResponse{}

	var err error

	task := func() error {
		res, err = t.listArticle(ctx, req)

		return err
	}
	err = t.retry(ctx, task)

	return res.Data.List, err
}

func (t *Client) GetCountOfArticle(ctx context.Context, req model.ListArticleRequest) (int, error) {
	res := model.ListArticleResponse{}

	var err error

	task := func() error {
		res, err = t.listArticle(ctx, req)

		return err
	}
	err = t.retry(ctx, task)

	return int(res.Data.PageBar.AllCount), err
}

func (t *Client) DeleteArticle(ctx context.Context, id string) error {
	var err error

	task := func() error {
		err = t.deleteArticle(ctx, id)

		return err
	}
	err = t.retry(ctx, task)

	return err
}

func (t *Client) ListCategory(ctx context.Context) ([]model.Category, error) {
	res := model.ListCategoryResponse{}

	var err error

	task := func() error {
		res, err = t.listCategory(ctx)

		return err
	}
	err = t.retry(ctx, task)

	return res.Data.List, err
}
