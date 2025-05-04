package infrastructure

import (
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

type RedisClient struct {
	client *redis.Client
}

func InitRedis() *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	return &RedisClient{client: client}
}

func (r *RedisClient) CacheMessage(messageID string, sentAt time.Time) {
	ctx := context.Background()
	r.client.Set(ctx, messageID, sentAt.Format(time.RFC3339), 0)
}

func (r *RedisClient) GetKeys(pattern string) ([]string, error) {
	ctx := context.Background()
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}
	return keys, nil
}
