package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sinker/ssop/internal/services"
	"github.com/sinker/ssop/pkg/logger"
	"github.com/sinker/ssop/pkg/response"
)

// AuthMiddleware 认证中间件
func AuthMiddleware(authService services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "缺少认证信息")
			c.Abort()
			return
		}

		// 检查格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Unauthorized(c, "认证格式错误")
			c.Abort()
			return
		}

		// 获取token
		tokenString := parts[1]

		// 验证token
		claims, err := authService.VerifyToken(tokenString)
		if err != nil {
			if err == services.ErrTokenBlacklisted {
				logger.Debug("Token已被撤销", "error", err)
				response.Unauthorized(c, "令牌已被撤销，请重新登录")
				c.Abort()
				return
			}

			logger.Debug("Token verification failed", "error", err)
			response.Unauthorized(c, "无效的认证信息")
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("userId", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("permissions", claims.Permissions)

		c.Next()
	}
}

// AuthorizePermission 权限检查中间件
func AuthorizePermission(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取用户权限
		permissions, exists := c.Get("permissions")
		if !exists {
			response.Unauthorized(c, "未授权的访问")
			c.Abort()
			return
		}

		// 检查权限
		userPermissions := permissions.([]string)
		hasPermission := false
		for _, permission := range userPermissions {
			if permission == requiredPermission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			response.Forbidden(c, "权限不足")
			c.Abort()
			return
		}

		c.Next()
	}
} 