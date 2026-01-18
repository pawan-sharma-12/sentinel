package cache

import (
	"context"
	"fmt"
	"time"
	"github.com/go-redis/redis/v8"
)
type RedisClient struct{
	Client *redis.Client
}
func NewRedisClient(addr string) *RedisClient{
	return &RedisClient{
		Client : redis.NewClient(&redis.Options{
			Addr : addr,
		}),
	}
}
func (r *RedisClient) IsRateLimited(ctx context.Context, userID string)(bool, error){
	key := fmt.Sprintf("rate_limit: %s", userID)
	count, err := r.Client.Incr(ctx, key).Result() 
	if err != nil {
		return false, fmt.Errorf("failed to increment rate limit: %w", err)
	}
	if count == 1{
		err := r.Client.Expire(ctx, key, 1 * time.Minute).Err()
		if err != nil{
			return false, fmt.Errorf("failed to set expiration : %w", err)
		}
	}
	const threshold = 5
	if count > threshold{
		return true, nil
	}
	return false, nil 

}