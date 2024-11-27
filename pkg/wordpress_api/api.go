package wordpressapi

import (
	"context"

	"github.com/google/uuid"

	"github.com/ray31245/seo_cluster/pkg/wordpress_api/model"
	wordpressinterface "github.com/ray31245/seo_cluster/pkg/wordpress_api/wordpress_interface"
)

// WordpressApi is a struct to implement WordpressApi interface
var _ wordpressinterface.WordpressAPI = &WordpressApi{}

type WordpressApi struct {
	// TODO: clientPool should be sync.Map
	// key is site id or user id
	clientPool map[uuid.UUID]*Client
}

func NewWordpressApi() *WordpressApi {
	return &WordpressApi{
		clientPool: make(map[uuid.UUID]*Client),
	}
}

func (t *WordpressApi) GetClient(ctx context.Context, ID uuid.UUID, urlStr string, userName string, password string) (wordpressinterface.WordpressClient, error) {
	client, ok := t.clientPool[ID]

	var err error
	if !ok {
		client, err = NewClient(ctx, urlStr, model.BasicAuthentication{
			Username: userName,
			Password: password,
		})
		if err != nil {
			return nil, err
		}

		t.clientPool[ID] = client
	}

	return client, nil
}

func (t *WordpressApi) UpdateClient(ctx context.Context, ID uuid.UUID, urlStr string, userName string, password string) (wordpressinterface.WordpressClient, error) {
	client, err := NewClient(ctx, urlStr, model.BasicAuthentication{
		Username: userName,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	t.clientPool[ID] = client

	return client, nil
}

func (t *WordpressApi) DeleteClient(ID uuid.UUID) {
	delete(t.clientPool, ID)
}

func (t *WordpressApi) NewClient(ctx context.Context, urlStr string, userName string, password string) (wordpressinterface.WordpressClient, error) {
	return NewClient(ctx, urlStr, model.BasicAuthentication{
		Username: userName,
		Password: password,
	})
}

func (t *WordpressApi) NewAnonymousClient(ctx context.Context, urlStr string) wordpressinterface.WordpressClient {
	return NewAnonymousClient(ctx, urlStr)
}
