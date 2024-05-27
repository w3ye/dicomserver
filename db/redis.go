package db

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Ctx    context.Context
	Client *redis.Client
}

func NewRedisClient(ctx context.Context, host string) (*Redis, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: "",
		DB:       0,
	})
	rdb.Set(ctx, "key", "value", 0)
	return &Redis{
		Client: rdb,
		Ctx:    ctx,
	}, nil
}
