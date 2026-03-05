package cache

import (
	"context"
	"errors"
	"time"
)

type ICache interface {
	// Set v value to cache on key k. v is anything serializable by json.Marshal()
	Set(ctx context.Context, k string, v any, ttl time.Duration) error
	// Get dest from key k. dest is a pointer to struct. dest shape must match the data unmarshaled from the cache.
	Get(ctx context.Context, k string, dest any) error
	// Delete cache value from key k
	Delete(ctx context.Context, k string) error
	// Close cache connection
	Close() error
}

var (
	ErrCacheMiss = errors.New("cache miss")
)
