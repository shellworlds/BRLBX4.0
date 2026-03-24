package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// PredictCache stores serialized prediction payloads (optional).
type PredictCache interface {
	Get(ctx context.Context, key string, dest any) bool
	Set(ctx context.Context, key string, v any, ttl time.Duration) error
}

type Noop struct{}

func (Noop) Get(_ context.Context, _ string, _ any) bool { return false }

func (Noop) Set(_ context.Context, _ string, _ any, _ time.Duration) error { return nil }

type Redis struct {
	c        *redis.Client
	prefix   string
	disabled bool
}

func NewRedis(addr, password string, db int, prefix string) *Redis {
	if addr == "" {
		return &Redis{disabled: true}
	}
	return &Redis{
		c:      redis.NewClient(&redis.Options{Addr: addr, Password: password, DB: db}),
		prefix: prefix,
	}
}

func (r *Redis) key(k string) string {
	return r.prefix + k
}

func (r *Redis) Get(ctx context.Context, key string, dest any) bool {
	if r == nil || r.disabled || r.c == nil {
		return false
	}
	b, err := r.c.Get(ctx, r.key(key)).Bytes()
	if err != nil {
		return false
	}
	return json.Unmarshal(b, dest) == nil
}

func (r *Redis) Set(ctx context.Context, key string, v any, ttl time.Duration) error {
	if r == nil || r.disabled || r.c == nil {
		return nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return r.c.Set(ctx, r.key(key), b, ttl).Err()
}
