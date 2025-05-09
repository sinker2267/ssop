package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/sinker/ssop/pkg/logger"
	"github.com/sinker/ssop/pkg/redis"
)

const (
	// TokenBlacklistPrefix Redis中Token黑名单的前缀
	TokenBlacklistPrefix = "token:blacklist:"
)

// TokenService 令牌服务接口
type TokenService interface {
	AddToBlacklist(tokenString string, claims *Claims) error
	IsBlacklisted(tokenString string) (bool, error)
}

// tokenService 令牌服务实现
type tokenService struct{}

// NewTokenService 创建令牌服务
func NewTokenService() TokenService {
	return &tokenService{}
}

// AddToBlacklist 将令牌添加到黑名单
func (s *tokenService) AddToBlacklist(tokenString string, claims *Claims) error {
	if tokenString == "" {
		return errors.New("令牌为空")
	}

	// 计算剩余过期时间
	expirationTime := claims.ExpiresAt.Time
	now := time.Now()

	// 如果已过期，无需加入黑名单
	if expirationTime.Before(now) {
		return nil
	}

	// 计算剩余过期时间（秒）
	remainingTime := expirationTime.Sub(now)

	// 将令牌添加到黑名单，过期时间设置为JWT的过期时间
	key := fmt.Sprintf("%s%s", TokenBlacklistPrefix, tokenString)
	if err := redis.Set(key, "revoked", remainingTime); err != nil {
		logger.Error("将令牌添加到黑名单失败", "error", err)
		return err
	}

	return nil
}

// IsBlacklisted 检查令牌是否在黑名单中
func (s *tokenService) IsBlacklisted(tokenString string) (bool, error) {
	if tokenString == "" {
		return false, errors.New("令牌为空")
	}

	key := fmt.Sprintf("%s%s", TokenBlacklistPrefix, tokenString)
	exists, err := redis.Exists(key)
	if err != nil {
		logger.Error("检查令牌是否在黑名单中失败", "error", err)
		return false, err
	}

	return exists, nil
}
