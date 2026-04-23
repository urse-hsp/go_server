package client

import (
	"context"
	"go-server/pkg/server/httpclient"
	"time"
)

type UserClient interface {
	GetRandomUser(ctx context.Context) (any, error)
}

type userClient struct {
	client *httpclient.Client
}

func NewUserClient() UserClient {
	return &userClient{
		client: httpclient.New("https://randomuser.me/api", 30*time.Second),
	}
}

func (c *userClient) GetRandomUser(ctx context.Context) (any, error) {
	var result any
	err := c.client.Get(ctx, "", &result)
	return result, err
}
