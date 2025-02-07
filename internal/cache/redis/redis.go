package redis

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"url-shortener/internal/cache"
	"url-shortener/internal/config"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

// New creates redis connection.
func New(ctx context.Context, cfg config.Redis) (*Cache, error) {
	const op = "cache.Redis.New"

	addr := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	connResp := client.Ping(ctx)
	if connResp.Val() != "PONG" {
		return nil, fmt.Errorf("%s: %w", op, cache.ErrNoPing)
	}

	return &Cache{
		client: client,
	}, nil
}
