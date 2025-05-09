package utils

import (
	"math/rand"
	"time"
)

var (
	// 确保随机种子在包初始化时被设置
	_ = func() bool {
		rand.Seed(time.Now().UnixNano())
		return true
	}()
)

// RandomFloat64 生成指定范围的随机浮点数
// min: 最小值
// max: 最大值
func RandomFloat64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// RandomInt 生成指定范围的随机整数
// min: 最小值(包含)
// max: 最大值(包含)
func RandomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

// RandomString 生成指定长度的随机字符串
// length: 字符串长度
// chars: 可选的字符集，默认使用字母和数字
func RandomString(length int, chars ...string) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	if len(chars) > 0 {
		charset = chars[0]
	}
	
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
} 