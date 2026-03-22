package service

import (
	"context"
	"fmt"
	"time"
)

// Cache is the interface AccountService and TransactionService depend on.
// RedisClient satisfies this interface; in tests you can use noopCacheImpl.
type Cache interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Get(ctx context.Context, key string, dest any) error
	Delete(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, key string) (bool, error)
}

// noopCacheImpl is a no-op Cache used in unit tests.
type noopCacheImpl struct{}

func (n *noopCacheImpl) Set(_ context.Context, _ string, _ any, _ time.Duration) error {
	return nil
}
func (n *noopCacheImpl) Get(_ context.Context, _ string, _ any) error {
	return fmt.Errorf("cache miss")
}
func (n *noopCacheImpl) Delete(_ context.Context, _ ...string) error      { return nil }
func (n *noopCacheImpl) Exists(_ context.Context, _ string) (bool, error) { return false, nil }
