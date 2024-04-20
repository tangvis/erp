package redis

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	defaultExpireTime = 24 * time.Hour
)

type Cache interface {
	Set(ctx context.Context, key string, value any) error
	SetEx(ctx context.Context, key string, value any, expire time.Duration) error
	SetExMarshal(ctx context.Context, key string, value any, expiration time.Duration) error
	GetBytes(ctx context.Context, key string) ([]byte, error)
	GetExUnmarshal(ctx context.Context, key string, value any, expiration time.Duration) error
	Del(ctx context.Context, key string) error
}

type Client struct {
	Cli *redis.Client
}

func (c Client) Set(ctx context.Context, key string, value any) error {
	return c.SetEx(ctx, key, value, defaultExpireTime)
}

func (c Client) SetEx(ctx context.Context, key string, value any, expiration time.Duration) error {
	return c.Cli.SetEx(ctx, key, value, expiration).Err()
}

func (c Client) SetExMarshal(ctx context.Context, key string, value any, expiration time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.Cli.SetEx(ctx, key, b, expiration).Err()
}

func (c Client) Del(ctx context.Context, key string) error {
	return c.Cli.Del(ctx, key).Err()
}

func (c Client) GetExUnmarshal(ctx context.Context, key string, value any, expiration time.Duration) error {
	b, err := c.Cli.GetEx(ctx, key, expiration).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &value)
}

func (c Client) GetBytes(ctx context.Context, key string) ([]byte, error) {
	return c.Cli.Get(ctx, key).Bytes()
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
