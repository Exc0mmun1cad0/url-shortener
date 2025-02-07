package redis

import (
	"context"
	"errors"
	"fmt"
	"url-shortener/internal/cache"

	"github.com/redis/go-redis/v9"
)

func (c *Cache) Insert(ctx context.Context, key string, value string) error {
	const op = "cache.Redis.Insert"

	err := c.client.Set(ctx, key, value, 0).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	const op = "cache.Redis.Get"

	resp := c.client.Get(ctx, key)

	result, err := resp.Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", fmt.Errorf("%s: %w", op, cache.ErrNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	const op = "cache.Redis.Delete"

	resp := c.client.Del(ctx, key)

	result, err := resp.Result()
	if result == 0 {
		return fmt.Errorf("%s: %w", op, cache.ErrNotFound)
	}
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
