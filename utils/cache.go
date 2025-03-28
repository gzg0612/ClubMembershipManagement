package utils

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

var (
	redisClient *redis.Client
	ctx         = context.Background()
)

// InitCache 初始化Redis缓存
func InitCache(addr string, password string, db int) error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// 测试连接
	_, err := redisClient.Ping(ctx).Result()
	return err
}

// Cache 设置缓存
func Cache(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return redisClient.Set(ctx, key, data, expiration).Err()
}

// GetCache 获取缓存
func GetCache(key string, value interface{}) error {
	data, err := redisClient.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, value)
}

// DeleteCache 删除缓存
func DeleteCache(key string) error {
	return redisClient.Del(ctx, key).Err()
}

// ClearCache 清空所有缓存
func ClearCache() error {
	return redisClient.FlushAll(ctx).Err()
}
