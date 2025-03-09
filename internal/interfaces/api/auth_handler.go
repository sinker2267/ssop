package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sinker/ssop/internal/application/dto"
	"github.com/sinker/ssop/internal/application/service"
	"github.com/sinker/ssop/pkg/errors"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register 注册路由
func (h *AuthHandler) Register(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.RefreshToken)
		auth.GET("/profile", h.GetProfile)
	}
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "无效的请求参数",
		})
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{
				"code":    appErr.Code,
				"message": appErr.Message,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "服务器内部错误",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "注册成功",
		"data":    user,
	})
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "无效的请求参数",
		})
		return
	}

	token, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{
				"code":    appErr.Code,
				"message": appErr.Message,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "服务器内部错误",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "登录成功",
		"data":    token,
	})
}

// RefreshToken 刷新令牌
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "无效的请求参数",
		})
		return
	}

	token, err := h.authService.RefreshToken(c.Request.Context(), req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{
				"code":    appErr.Code,
				"message": appErr.Message,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "服务器内部错误",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "刷新令牌成功",
		"data":    token,
	})
}

// GetProfile 获取用户信息
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "未认证",
		})
		return
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), userID.(uint))
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{
				"code":    appErr.Code,
				"message": appErr.Message,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "服务器内部错误",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "成功",
		"data":    user,
	})
} 