package zblogapi

import (
	zInterface "goTool/pkg/z_blog_api/zblog_Interface"

	"github.com/google/uuid"
)

type ZblogAPI struct {
	// TODO: clientPool should be sync.Map
	// key is site id
	clientPool map[uuid.UUID]*ZblogAPIClient
}

func NewZblogAPI() *ZblogAPI {
	return &ZblogAPI{
		clientPool: make(map[uuid.UUID]*ZblogAPIClient),
	}
}

func (t *ZblogAPI) GetClient(siteID uuid.UUID, urlStr string, userName string, password string) (zInterface.ZBlogApiClient, error) {
	client, ok := t.clientPool[siteID]
	var err error
	if !ok {
		client, err = NewZblogAPIClient(urlStr, userName, password)
		if err != nil {
			return nil, err
		}
		t.clientPool[siteID] = client
	}
	return client, nil
}

func (t *ZblogAPI) NewClient(urlStr string, userName string, password string) (zInterface.ZBlogApiClient, error) {
	return NewZblogAPIClient(urlStr, userName, password)
}
