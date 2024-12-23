package cache

import (
	redis "github.com/go-redis/redis/v8"
	"time"
)

// Redis is a Redis cache client
type Redis struct {
	client *redis.Client
}

// NewRedisClient creates a new Redis client
func NewRedisClient() *Redis {
	return &Redis{
		client: redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		}),
	}
}

func (r *Redis) Set(key string, value string, expiration time.Duration) error {
	return r.client.Set(r.client.Context(), key, value, expiration).Err()
}

func (r *Redis) Get(key string) (string, error) {
	return r.client.Get(r.client.Context(), key).Result()
}

func (r *Redis) Del(key string) error {
	return r.client.Del(r.client.Context(), key).Err()
}

func (r *Redis) SExists(key string, member string) (bool, error) {
	cmd := r.client.SIsMember(r.client.Context(), key, member)
	return cmd.Val(), cmd.Err()
}

func (r *Redis) SAdd(key string, members ...string) error {
	return r.client.SAdd(r.client.Context(), key, members).Err()
}

func (r *Redis) Close() error {
	return r.client.Close()
}
