package cache

import (
	"biometrics/internal/config"
	"biometrics/pkg/utils"
	"context"
	"time"
)

type Redis struct {
	Client interface{}
	log    *utils.Logger
}

func NewRedis(cfg config.RedisConfig) (*Redis, error) {
	log := utils.NewLogger("info", "development")

	return &Redis{
		Client: nil,
		log:    log,
	}, nil
}

func (r *Redis) Close() error {
	return nil
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	return "", nil
}

func (r *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return nil
}

func (r *Redis) Del(ctx context.Context, keys ...string) error {
	return nil
}

func (r *Redis) Exists(ctx context.Context, keys ...string) (int64, error) {
	return 0, nil
}

func (r *Redis) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return nil
}

func (r *Redis) Incr(ctx context.Context, key string) (int64, error) {
	return 0, nil
}

func (r *Redis) Decr(ctx context.Context, key string) (int64, error) {
	return 0, nil
}
