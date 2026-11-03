package utils

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

func GetRedisStr(ctx context.Context, rdb *redis.Client, key string) (string, error) {
	return rdb.Get(ctx, key).Result()
}

func GetRedisJson[T any](ctx context.Context, rdb *redis.Client, key string) (T, error) {
	var data T
	str, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return data, err
	}

	if err := json.Unmarshal([]byte(str), &data); err != nil {
		return data, err
	}

	return data, nil
}

func SetRedisJson[T any](ctx context.Context, rdb *redis.Client, key string, data T, expiration ...time.Duration) error {
	var ttl time.Duration
	if len(expiration) > 0 {
		ttl = expiration[0]
	}

	str, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return rdb.Set(ctx, key, str, ttl).Err()
}

func DelRedis(ctx context.Context, rdb *redis.Client, key string) error {
	_, err := rdb.Del(ctx, key).Result()
	if err != nil {
		return err
	}

	return nil
}
