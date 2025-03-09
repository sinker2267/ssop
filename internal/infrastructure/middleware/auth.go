package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sinker/ssop/pkg/auth"
	"github.com/sinker/ssop/pkg/errors"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware(jwtConfig auth.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取令牌
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"code":    401,
				"message": "未提供认证令牌",
			})
			return
		}

		// 解析令牌
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.AbortWithStatusJSON(401, gin.H{
				"code":    401,
				"message": "认证令牌格式无效",
			})
			return
		}

		// 验证令牌
		claims, err := auth.ParseToken(parts[1], jwtConfig)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"code":    401,
				"message": "无效的认证令牌",
			})
			return
		}

		// 将用户信息存储到上下文
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

// ErrorHandler 错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				// 处理应用错误
				if appErr, ok := e.Err.(*errors.AppError); ok {
					c.AbortWithStatusJSON(appErr.Code, gin.H{
						"code":    appErr.Code,
						"message": appErr.Message,
					})
					return
				}
			}
			
			// 处理其他错误
			c.AbortWithStatusJSON(500, gin.H{
				"code":    500,
				"message": "服务器内部错误",
			})
		}
	}
}

// CORS 跨域中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
} 