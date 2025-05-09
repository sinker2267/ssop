package utils

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	// 字符集
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// GenerateID 生成唯一ID
// prefix: ID前缀，如user, ds, task等
// 返回格式: prefix + 当前时间戳 + 随机字符串
func GenerateID(prefix string) string {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	randomStr := RandomString(6)
	return fmt.Sprintf("%s_%d_%s", prefix, timestamp, randomStr)
}

// GenerateToken 生成随机令牌
func GenerateToken(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
} 