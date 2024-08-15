package zBlogApi

import (
	zInterface "goTool/pkg/z_blog_api/z_blog_Interface"

	"github.com/google/uuid"
)

type ZBlogAPI struct {
	// TODO: clientPool should be sync.Map
	// key is site id
	clientPool map[uuid.UUID]*ZBlogAPIClient
}

func NewZBlogAPI() *ZBlogAPI {
	return &ZBlogAPI{
		clientPool: make(map[uuid.UUID]*ZBlogAPIClient),
	}
}

func (t *ZBlogAPI) GetClient(siteID uuid.UUID, urlStr string, userName string, password string) (zInterface.ZBlogApiClient, error) {
	client, ok := t.clientPool[siteID]
	var err error
	if !ok {
		client, err = NewZBlogAPIClient(urlStr, userName, password)
		if err != nil {
			return nil, err
		}
		t.clientPool[siteID] = client
	}
	return client, nil
}

func (t *ZBlogAPI) NewClient(urlStr string, userName string, password string) (zInterface.ZBlogApiClient, error) {
	return NewZBlogAPIClient(urlStr, userName, password)
}
