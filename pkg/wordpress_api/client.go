package wordpressapi

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/ray31245/seo_cluster/pkg/wordpress_api/model"
	"github.com/ray31245/seo_cluster/pkg/wordpress_api/origin"
)

type Client struct {
	baseURL   string
	basicAuth model.BasicAuthentication
	// avoid duplicate login
	lock *sync.Mutex
}

func NewClient(ctx context.Context, urlStr string, basicAuth model.BasicAuthentication) (*Client, error) {
	log.Println("following url is used to login")
	log.Println(urlStr)

	res := &Client{
		baseURL:   urlStr,
		basicAuth: basicAuth,
		lock:      &sync.Mutex{},
	}

	_, err := res.RetrieveUserMe(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func NewAnonymousClient(ctx context.Context, urlStr string) *Client {
	res := &Client{
		baseURL: urlStr,
		basicAuth: model.BasicAuthentication{
			IsAnonymous: true,
		},
		lock: &sync.Mutex{},
	}

	return res
}

func (c *Client) RetrieveUserMe(ctx context.Context) (model.RetrieveUserMeResponse, error) {
	res, err := origin.RetrieveUserMe(ctx, c.baseURL, c.basicAuth, model.ContextView)
	if err != nil {
		return model.RetrieveUserMeResponse{}, fmt.Errorf("retrieve user me error: %w", err)
	}

	return res, nil
}

func (c *Client) ListTag(ctx context.Context, args model.ListTagArgs) (model.ListTagResponse, error) {
	res, _, err := origin.ListTag(ctx, c.baseURL, c.basicAuth, args)
	if err != nil {
		return model.ListTagResponse{}, fmt.Errorf("list tag error: %w", err)
	}

	return res, nil
}

func (c *Client) ListTagAll(ctx context.Context) (model.ListTagResponse, error) {
	res := model.ListTagResponse{}

	listArgs := model.ListTagArgs{
		Page: 1,
	}

	for {
		listRes, page, err := origin.ListTag(ctx, c.baseURL, c.basicAuth, listArgs)
		if err != nil {
			return model.ListTagResponse{}, fmt.Errorf("list tag all error: %w", err)
		}

		res = append(res, listRes...)

		if len(res) >= page.Total {
			break
		}

		listArgs.Page += 1
	}

	return res, nil
}

func (c *Client) CreateTag(ctx context.Context, args model.CreateTagArgs) (model.CreateTagResponse, error) {
	res, err := origin.CreateTag(ctx, c.baseURL, c.basicAuth, args)
	if err != nil {
		return model.CreateTagResponse{}, fmt.Errorf("create tag error: %w", err)
	}

	return res, nil
}

func (c *Client) ListCategory(ctx context.Context, args model.ListCategoryArgs) (model.ListCategoryResponse, error) {
	res, err := origin.ListCategory(ctx, c.baseURL, c.basicAuth, args)
	if err != nil {
		return model.ListCategoryResponse{}, fmt.Errorf("list category error: %w", err)
	}

	return res, nil
}

func (c *Client) GetCountOfArticle(ctx context.Context, req model.ListArticleArgs) (int, error) {
	_, page, err := origin.ListArticle(ctx, c.baseURL, c.basicAuth, req)
	if err != nil {
		return 0, fmt.Errorf("list article error: %w", err)
	}

	return page.Total, nil
}

func (c *Client) ListArticle(ctx context.Context, args model.ListArticleArgs) (model.ListArticleResponse, error) {
	res, _, err := origin.ListArticle(ctx, c.baseURL, c.basicAuth, args)
	if err != nil {
		return model.ListArticleResponse{}, fmt.Errorf("list article error: %w", err)
	}

	return res, nil
}

func (c *Client) CreateArticle(ctx context.Context, args model.CreateArticleArgs) (model.CreateArticleResponse, error) {
	res, err := origin.CreateArticle(ctx, c.baseURL, c.basicAuth, args)
	if err != nil {
		return model.CreateArticleResponse{}, fmt.Errorf("create article error: %w", err)
	}

	return res, nil
}

func (c *Client) UpdateArticle(ctx context.Context, args model.UpdateArticleArgs) (model.UpdateArticleResponse, error) {
	res, err := origin.UpdateArticle(ctx, c.baseURL, c.basicAuth, args)
	if err != nil {
		return model.UpdateArticleResponse{}, fmt.Errorf("update article error: %w", err)
	}

	return res, nil
}

func (c *Client) RetrieveArticle(ctx context.Context, args model.RetrieveArticleArgs) (model.RetrieveArticleResponse, error) {
	res, err := origin.RetrieveArticle(ctx, c.baseURL, c.basicAuth, args)
	if err != nil {
		return model.RetrieveArticleResponse{}, fmt.Errorf("retrieve article error: %w", err)
	}

	return res, nil
}

func (c *Client) CreateComment(ctx context.Context, args model.CreateCommentArgs) (model.CreateCommentResponse, error) {
	res, err := origin.CreateComment(ctx, c.baseURL, c.basicAuth, args)
	if err != nil {
		return model.CreateCommentResponse{}, fmt.Errorf("create comment error: %w", err)
	}

	return res, nil
}
