package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sinker/ssop/internal/config"
	"github.com/sinker/ssop/pkg/logger"
)

var (
	// Client Redis客户端实例
	Client *redis.Client
	ctx    = context.Background()
)

// InitRedis 初始化Redis连接
func InitRedis(config config.RedisConfig) error {
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	// 测试连接
	if err := Client.Ping(ctx).Err(); err != nil {
		return err
	}

	logger.Info("Redis连接成功", "host", config.Host, "port", config.Port)
	return nil
}

// Get 获取键值
func Get(key string) (string, error) {
	return Client.Get(ctx, key).Result()
}

// Set 设置键值
func Set(key string, value interface{}, expiration time.Duration) error {
	return Client.Set(ctx, key, value, expiration).Err()
}

// SetEx 设置键值并设置过期时间（秒）
func SetEx(key string, value interface{}, seconds int) error {
	return Client.Set(ctx, key, value, time.Duration(seconds)*time.Second).Err()
}

// Del 删除键
func Del(key string) error {
	return Client.Del(ctx, key).Err()
}

// Exists 检查键是否存在
func Exists(key string) (bool, error) {
	val, err := Client.Exists(ctx, key).Result()
	return val > 0, err
}

// Keys 获取匹配的键
func Keys(pattern string) ([]string, error) {
	return Client.Keys(ctx, pattern).Result()
}

// Expire 设置过期时间（秒）
func Expire(key string, seconds int) error {
	return Client.Expire(ctx, key, time.Duration(seconds)*time.Second).Err()
}

// Close 关闭连接
func Close() error {
	if Client != nil {
		return Client.Close()
	}
	return nil
} 