package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTConfig JWT配置
type JWTConfig struct {
	SecretKey     string        // 密钥
	TokenExpiry   time.Duration // 令牌过期时间
	RefreshExpiry time.Duration // 刷新令牌过期时间
	Issuer        string        // 签发者
}

// DefaultJWTConfig 默认JWT配置
var DefaultJWTConfig = JWTConfig{
	SecretKey:     "your-secret-key-change-in-production", // 生产环境中应该更改为复杂的随机字符串
	TokenExpiry:   time.Hour * 24,                        // 令牌有效期24小时
	RefreshExpiry: time.Hour * 24 * 7,                    // 刷新令牌有效期7天
	Issuer:        "ssop-api",                            // 签发者
}

// CustomClaims 自定义JWT声明
type CustomClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID uint, username string, config JWTConfig) (string, string, error) {
	// 生成访问令牌
	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.TokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    config.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return "", "", err
	}

	// 生成刷新令牌
	refreshClaims := CustomClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.RefreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    config.Issuer,
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(config.SecretKey))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshTokenString, nil
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString string, config JWTConfig) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// 验证签名方法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid token signing method")
			}
			return []byte(config.SecretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// RefreshToken 刷新令牌
func RefreshToken(refreshToken string, config JWTConfig) (string, string, error) {
	// 解析刷新令牌
	claims, err := ParseToken(refreshToken, config)
	if err != nil {
		return "", "", err
	}

	// 生成新令牌
	return GenerateToken(claims.UserID, claims.Username, config)
} 