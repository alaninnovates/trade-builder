package database

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	context context.Context
	client  *redis.Client
}

func NewRedis() *Redis {
	return &Redis{
		context: context.Background(),
	}
}

func (r *Redis) Context() context.Context {
	return r.context
}

func (r *Redis) Client() *redis.Client {
	return r.client
}

func (r *Redis) Connect(uri string) (*redis.Client, error) {
	opts, err := redis.ParseURL(uri)
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(opts)
	r.client = client
	return client, nil
}
