package zblogapi

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/ray31245/seo_cluster/pkg/z_blog_api/model"
	"github.com/ray31245/seo_cluster/pkg/z_blog_api/origin"
)

type Client struct {
	baseURL  string
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

	res := &Client{
		baseURL:  urlStr,
		lock:     &sync.Mutex{},
		userName: userName,
		password: password,
		token:    "YmV2aXN8fHw3ZWYxMmJkNTQ1ZmU5MTRhNTMwYTFlYjMyODUxYTA5YTg4YjE0OGRmYjExN2Y2ODRkZmZmNzM1ZjM2YTcwMmI4MTcyMjQxODY0OA==",
	}

	err := res.Login(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *Client) Login(ctx context.Context) error {
	resData, err := origin.Login(ctx, t.baseURL, t.token, t.userName, t.password)
	if err != nil {
		return fmt.Errorf("login error: %w", err)
	}

	t.token = resData.Data.Token

	log.Println("login success")

	return nil
}

func (t *Client) ListMember(ctx context.Context) (model.ListMemberResponse, error) {
	res := model.ListMemberResponse{}

	var err error

	task := func() error {
		res, err = origin.ListMember(ctx, t.baseURL, t.token)

		return err
	}
	err = t.retry(ctx, task)

	return res, err
}

func (t *Client) PostArticle(ctx context.Context, art model.PostArticleRequest) error {
	var err error

	task := func() error {
		err = origin.PostArticle(ctx, t.baseURL, t.token, art)
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
		res, err = origin.ListArticle(ctx, t.baseURL, t.token, req)

		return err
	}
	err = t.retry(ctx, task)

	return res.Data.List, err
}

func (t *Client) GetCountOfArticle(ctx context.Context, req model.ListArticleRequest) (int, error) {
	res := model.ListArticleResponse{}

	var err error

	task := func() error {
		res, err = origin.ListArticle(ctx, t.baseURL, t.token, req)

		return err
	}
	err = t.retry(ctx, task)

	return int(res.Data.PageBar.AllCount), err
}

func (t *Client) DeleteArticle(ctx context.Context, id string) error {
	var err error

	task := func() error {
		err = origin.DeleteArticle(ctx, t.baseURL, t.token, id)

		return err
	}
	err = t.retry(ctx, task)

	return err
}

func (t *Client) ListCategory(ctx context.Context) ([]model.Category, error) {
	res := model.ListCategoryResponse{}

	var err error

	task := func() error {
		res, err = origin.ListCategory(ctx, t.baseURL, t.token)

		return err
	}
	err = t.retry(ctx, task)

	return res.Data.List, err
}
