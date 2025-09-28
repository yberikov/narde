package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"narde/internal/config"
	"time"
)

func setupCacheConnection(ctx context.Context, cfg config.DragonflyCreds) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:    cfg.Password,
		DB:          cfg.DB,
		PoolSize:    cfg.PoolSize,
		PoolTimeout: time.Second,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("can't ping cache: %w", err)
	}

	return client, nil
}
