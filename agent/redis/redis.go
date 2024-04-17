package redis

import (
	"context"
	"github.com/gin-contrib/sessions"
	ginRedis "github.com/gin-contrib/sessions/redis"
	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Set(ctx context.Context, key string, value any, expire int) error
	GetSessionStore() ginRedis.Store
}

type Client struct {
	Cli   *redis.Client
	Store ginRedis.Store
}

func (c Client) GetSessionStore() ginRedis.Store {
	return c.Store
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
		Store: initCacheStore(config),
	}
}

func initCacheStore(config Config) ginRedis.Store {
	store, err := ginRedis.NewStore(
		10,
		"tcp",
		config.Addr,
		config.Passwd,
		[]byte("secret"),
	)
	if err != nil {
		panic(err)
	}
	store.Options(sessions.Options{
		MaxAge:   86400,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})

	return store
}
