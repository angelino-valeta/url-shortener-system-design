package repository

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

// RedisRepository for counter and cache
type RedisRepository struct {
	client *redis.Client
	logger *zap.Logger
}

// NewRedisRepository
func NewRedisRepository(addr, password string, logger *zap.Logger) *RedisRepository {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	return &RedisRepository{client: client, logger: logger}
}

// IncrementCounter atomic ID
func (r *RedisRepository) IncrementCounter(ctx context.Context) (int64, error) {
	return r.client.Incr(ctx, "url_counter").Result()
}

// GetCachedURL
func (r *RedisRepository) GetCachedURL(ctx context.Context, shortcode string) (string, error) {
	return r.client.Get(ctx, "url:"+shortcode).Result()
}

// SetCachedURL with TTL
func (r *RedisRepository) SetCachedURL(ctx context.Context, shortcode, longURL string, ttl int) error {
	return r.client.Set(ctx, "url:"+shortcode, longURL, time.Duration(ttl)*time.Second).Err()
}
