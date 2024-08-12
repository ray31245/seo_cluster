package zblogapi

import (
	"goTool/pkg/z_blog_api/model"

	"github.com/google/uuid"
)

var (
	// TODO: clientPool should be sync.Map
	// key is site id
	clientPool map[uuid.UUID]*ZblogAPIClient = make(map[uuid.UUID]*ZblogAPIClient)
)

type ZblogAPI struct {
	client *ZblogAPIClient
}

func NewZblogAPI(siteID uuid.UUID, urlStr string, userName string, password string) (*ZblogAPI, error) {
	client, ok := clientPool[siteID]
	if !ok {
		client, err := NewZblogAPIClient(urlStr, userName, password)
		if err != nil {
			return nil, err
		}
		clientPool[siteID] = client
	}
	return &ZblogAPI{
		client: client,
	}, nil
}

func (t *ZblogAPI) PostArticle(article model.PostArticleRequest) error {
	return t.client.PostArticle(article)
}

func (t *ZblogAPI) ListArticle(req model.ListArticleRequest) ([]model.Article, error) {
	return t.client.ListArticle(req)
}

func (t *ZblogAPI) GetCountOfArticle(req model.ListArticleRequest) (int, error) {
	return t.client.GetCountOfArticle(req)
}

func (t *ZblogAPI) DeleteArticle(id string) error {
	return t.client.DeleteArticle(id)
}

func (t *ZblogAPI) ListCategory() ([]model.Category, error) {
	return t.client.ListCategory()
}

func (t *ZblogAPI) ListMember() (string, error) {
	return t.client.ListMember()
}
