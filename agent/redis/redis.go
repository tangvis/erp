package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Set(ctx context.Context, key string, value any, expire int) error
}

type Client struct {
	Cli *redis.Client
}

func (c Client) Set(ctx context.Context, key string, value any, expire int) error {
	//TODO implement me
	panic("implement me")
}

func NewCache(config Config) Cache {
	return &Client{
		Cli: redis.NewClient(&redis.Options{
			Addr:     config.Addr,
			Password: config.Passwd,
			DB:       config.DB,
		}),
	}
}
