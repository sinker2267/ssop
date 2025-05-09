package handlers

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sinker/ssop/internal/services"
	"github.com/sinker/ssop/pkg/logger"
	"github.com/sinker/ssop/pkg/response"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService services.AuthService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register 注册处理
func (h *AuthHandler) Register(c *gin.Context) {
	var registerDTO services.RegisterDTO
	if err := c.ShouldBindJSON(&registerDTO); err != nil {
		response.BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	user, err := h.authService.Register(registerDTO)
	if err != nil {
		logger.Error("用户注册失败", "error", err)
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"userId":     user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"role":       user.Role,
		"createTime": user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, "注册成功")
}

// Login 登录处理
func (h *AuthHandler) Login(c *gin.Context) {
	var loginDTO services.LoginDTO
	if err := c.ShouldBindJSON(&loginDTO); err != nil {
		response.BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	tokenPair, err := h.authService.Login(loginDTO)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			response.Unauthorized(c, "用户名或密码错误")
			return
		}
		logger.Error("用户登录失败", "error", err)
		response.ServerError(c, "登录失败")
		return
	}

	response.Success(c, tokenPair, "登录成功")
}

// RefreshToken 刷新令牌
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	tokenPair, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		if err == services.ErrInvalidRefreshToken || err == services.ErrRefreshTokenExpired || err == services.ErrTokenBlacklisted {
			response.Unauthorized(c, err.Error())
			return
		}
		logger.Error("刷新令牌失败", "error", err)
		response.ServerError(c, "刷新令牌失败")
		return
	}

	response.Success(c, gin.H{
		"token":     tokenPair.Token,
		"expiresIn": tokenPair.ExpiresIn,
	}, "刷新成功")
}

// Logout 退出登录
func (h *AuthHandler) Logout(c *gin.Context) {
	// 从请求头获取token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		response.Success(c, nil, "登出成功")
		return
	}

	// 解析token
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		response.Success(c, nil, "登出成功")
		return
	}

	// 获取token
	tokenString := parts[1]

	// 将token加入黑名单
	if err := h.authService.Logout(tokenString); err != nil {
		logger.Error("登出失败", "error", err)
		response.ServerError(c, "登出失败")
		return
	}

	response.Success(c, nil, "登出成功")
}

// GuestLogin 游客登录
func (h *AuthHandler) GuestLogin(c *gin.Context) {
	tokenPair, err := h.authService.GuestLogin()
	if err != nil {
		logger.Error("游客登录失败", "error", err)
		response.ServerError(c, "游客登录失败")
		return
	}

	response.Success(c, tokenPair, "游客登录成功")
}

// RegisterAuthRoutes 注册认证路由
func RegisterAuthRoutes(router *gin.RouterGroup, authService services.AuthService) {
	handler := NewAuthHandler(authService)

	auth := router.Group("/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
		auth.POST("/refresh-token", handler.RefreshToken)
		auth.POST("/logout", handler.Logout)
		auth.POST("/guest-login", handler.GuestLogin)
	}
}
