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
	lock        *sync.Mutex
	isAnonymous bool
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

func NewAnonymousClient(ctx context.Context, urlStr string) *Client {
	res := &Client{
		baseURL:     urlStr,
		lock:        &sync.Mutex{},
		isAnonymous: true,
	}

	return res
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

func (t *Client) ListMember(ctx context.Context) ([]model.Member, error) {
	res := model.ListMemberResponse{}

	var err error

	task := func() error {
		res, err = origin.ListMember(ctx, t.baseURL, t.token)
		if err != nil {
			return fmt.Errorf("ListMember: %w", err)
		}

		return nil
	}
	err = t.retry(ctx, task)

	return res.Data.List, err
}

func (t *Client) PostMember(ctx context.Context, mem model.PostMemberRequest) error {
	var err error

	if mem.ID == "" {
		mem.ID = "0"
	}

	task := func() error {
		err = origin.PostMember(ctx, t.baseURL, t.token, mem)
		if err != nil {
			return fmt.Errorf("PostMember: %w", err)
		}

		return nil
	}
	err = t.retry(ctx, task)

	return err
}

func (t *Client) PostArticle(ctx context.Context, art model.PostArticleRequest) (model.Article, error) {
	res := model.PostArticleResponse{}

	var err error

	task := func() error {
		res, err = origin.PostArticle(ctx, t.baseURL, t.token, art)
		if err != nil {
			return fmt.Errorf("PostArticle: %w", err)
		}

		return nil
	}
	err = t.retry(ctx, task)

	return res.Data.Post, err
}

func (t *Client) GetArticle(ctx context.Context, id string) (model.Article, error) {
	res := model.GetArticleResponse{}

	var err error

	task := func() error {
		res, err = origin.GetArticle(ctx, t.baseURL, t.token, id)
		if err != nil {
			return fmt.Errorf("GetArticle: %w", err)
		}

		return nil
	}
	if !t.isAnonymous {
		err = t.retry(ctx, task)
	} else {
		err = task()
	}

	return res.Data.Post, err
}

func (t *Client) ListArticle(ctx context.Context, req model.ListArticleRequest) ([]model.Article, error) {
	res := model.ListArticleResponse{}

	var err error

	task := func() error {
		res, err = origin.ListArticle(ctx, t.baseURL, t.token, req)
		if err != nil {
			return fmt.Errorf("ListArticle: %w", err)
		}

		return nil
	}
	if !t.isAnonymous {
		err = t.retry(ctx, task)
	} else {
		err = task()
	}

	return res.Data.List, err
}

func (t *Client) GetCountOfArticle(ctx context.Context, req model.ListArticleRequest) (int, error) {
	res := model.ListArticleResponse{}

	var err error

	task := func() error {
		res, err = origin.ListArticle(ctx, t.baseURL, t.token, req)
		if err != nil {
			return fmt.Errorf("ListArticle: %w", err)
		}

		return nil
	}
	err = t.retry(ctx, task)

	return int(res.Data.PageBar.AllCount), err
}

func (t *Client) DeleteArticle(ctx context.Context, id string) error {
	var err error

	task := func() error {
		err = origin.DeleteArticle(ctx, t.baseURL, t.token, id)
		if err != nil {
			return fmt.Errorf("DeleteArticle: %w", err)
		}

		return nil
	}
	err = t.retry(ctx, task)

	return err
}

func (t *Client) PostComment(ctx context.Context, comment model.PostCommentRequest) error {
	var err error

	task := func() error {
		err = origin.PostComment(ctx, t.baseURL, t.token, comment)
		if err != nil {
			return fmt.Errorf("PostComment: %w", err)
		}

		return nil
	}
	err = t.retry(ctx, task)

	return err
}

func (t *Client) ListCategory(ctx context.Context) ([]model.Category, error) {
	res := model.ListCategoryResponse{}

	var err error

	task := func() error {
		res, err = origin.ListCategory(ctx, t.baseURL, t.token)
		if err != nil {
			return fmt.Errorf("ListCategory: %w", err)
		}

		return nil
	}
	err = t.retry(ctx, task)

	return res.Data.List, err
}

func (t *Client) ListTag(ctx context.Context, req model.ListTagRequest) ([]model.Tag, error) {
	res := model.ListTagResponse{}

	var err error

	task := func() error {
		res, err = origin.ListTag(ctx, t.baseURL, t.token, req)
		if err != nil {
			return fmt.Errorf("ListTag: %w", err)
		}

		return nil
	}
	err = t.retry(ctx, task)

	return res.Data.List, err
}

func (t *Client) ListTagAll(ctx context.Context) ([]model.Tag, error) {
	res := []model.Tag{}
	taskRes := model.ListTagResponse{}

	var err error

	taskReq := model.ListTagRequest{
		PageRequest: model.PageRequest{
			Page: 1,
		},
	}
	task := func() error {
		taskRes, err = origin.ListTag(ctx, t.baseURL, t.token, taskReq)
		if err != nil {
			return fmt.Errorf("ListTagAll: %w", err)
		}

		return nil
	}

	for {
		err = t.retry(ctx, task)
		if err != nil {
			return res, err
		}

		res = append(res, taskRes.Data.List...)

		if taskRes.Data.PageBar.PageCurrent >= taskRes.Data.PageBar.PageAll {
			break
		}
		taskReq.Page = taskRes.Data.PageBar.PageCurrent + 1
	}

	return res, err
}

func (t *Client) PostTag(ctx context.Context, tag model.PostTagRequest) (model.Tag, error) {
	res := model.PostTagResponse{}

	var err error

	task := func() error {
		res, err = origin.PostTag(ctx, t.baseURL, t.token, tag)
		if err != nil {
			return fmt.Errorf("PostTag: %w", err)
		}

		return nil
	}
	err = t.retry(ctx, task)

	return res.Data.Tag, err
}
