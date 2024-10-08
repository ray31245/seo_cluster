package zblogapi

import (
	"context"

	zInterface "github.com/ray31245/seo_cluster/pkg/z_blog_api/z_blog_Interface"

	"github.com/google/uuid"
)

var _ zInterface.ZBlogAPI = (*ZBlogAPI)(nil)

type ZBlogAPI struct {
	// TODO: clientPool should be sync.Map
	// key is site id or user id
	clientPool map[uuid.UUID]*Client
}

func NewZBlogAPI() *ZBlogAPI {
	return &ZBlogAPI{
		clientPool: make(map[uuid.UUID]*Client),
	}
}

func (t *ZBlogAPI) GetClient(ctx context.Context, ID uuid.UUID, urlStr string, userName string, password string) (zInterface.ZBlogAPIClient, error) {
	client, ok := t.clientPool[ID]

	var err error
	if !ok {
		client, err = NewClient(ctx, urlStr, userName, password)
		if err != nil {
			return nil, err
		}

		t.clientPool[ID] = client
	}

	return client, nil
}

func (t *ZBlogAPI) UpdateClient(ctx context.Context, ID uuid.UUID, urlStr string, userName string, password string) (zInterface.ZBlogAPIClient, error) {
	client, err := NewClient(ctx, urlStr, userName, password)
	if err != nil {
		return nil, err
	}

	t.clientPool[ID] = client

	return client, nil
}

func (t *ZBlogAPI) DeleteClient(ID uuid.UUID) {
	delete(t.clientPool, ID)
}

func (t *ZBlogAPI) NewClient(ctx context.Context, urlStr string, userName string, password string) (zInterface.ZBlogAPIClient, error) {
	return NewClient(ctx, urlStr, userName, password)
}

func (t *ZBlogAPI) NewAnonymousClient(ctx context.Context, urlStr string) zInterface.ZBlogAPIClient {
	return NewAnonymousClient(ctx, urlStr)
}
