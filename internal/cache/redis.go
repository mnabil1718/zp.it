package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mnabil1718/zp.it/internal/config"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(cfg *config.Config) *RedisClient {

	opts, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(opts)

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic(fmt.Sprintf("redis connection error: %v", err))
	}

	return &RedisClient{client: rdb}
}

func (r *RedisClient) Set(ctx context.Context, k string, v any, ttl time.Duration) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, k, data, ttl).Err()
}

func (r *RedisClient) Get(ctx context.Context, k string, dest any) error {
	v, err := r.client.Get(ctx, k).Bytes()
	if err == redis.Nil {
		return ErrCacheMiss
	}

	if err != nil {
		return err
	}

	return json.Unmarshal(v, dest)
}

func (r *RedisClient) Delete(ctx context.Context, k string) error {
	return r.client.Del(ctx, k).Err()
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}
