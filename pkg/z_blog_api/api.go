package zblogapi

import (
	"context"

	zInterface "github.com/ray31245/seo_cluster/pkg/z_blog_api/z_blog_Interface"

	"github.com/google/uuid"
)

type ZBlogAPI struct {
	// TODO: clientPool should be sync.Map
	// key is site id
	clientPool map[uuid.UUID]*Client
}

func NewZBlogAPI() *ZBlogAPI {
	return &ZBlogAPI{
		clientPool: make(map[uuid.UUID]*Client),
	}
}

func (t *ZBlogAPI) GetClient(ctx context.Context, siteID uuid.UUID, urlStr string, userName string, password string) (zInterface.ZBlogAPIClient, error) {
	client, ok := t.clientPool[siteID]

	var err error
	if !ok {
		client, err = NewClient(ctx, urlStr, userName, password)
		if err != nil {
			return nil, err
		}

		t.clientPool[siteID] = client
	}

	return client, nil
}

func (t *ZBlogAPI) NewClient(ctx context.Context, urlStr string, userName string, password string) (zInterface.ZBlogAPIClient, error) {
	return NewClient(ctx, urlStr, userName, password)
}
