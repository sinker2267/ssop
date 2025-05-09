package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sinker/ssop/internal/services"
	"github.com/sinker/ssop/pkg/logger"
	"github.com/sinker/ssop/pkg/response"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService services.UserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetCurrentUser 获取当前用户信息
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		response.Unauthorized(c, "未授权的访问")
		return
	}

	userInfo, err := h.userService.GetCurrentUserInfo(userID.(string))
	if err != nil {
		logger.Error("获取用户信息失败", "error", err)
		if err == services.ErrUserNotFound {
			response.NotFound(c, "用户不存在")
			return
		}
		response.ServerError(c, "获取用户信息失败")
		return
	}

	response.Success(c, userInfo, "成功")
}

// RegisterUserRoutes 注册用户路由
func RegisterUserRoutes(router *gin.RouterGroup, userService services.UserService, authMiddleware gin.HandlerFunc) {
	handler := NewUserHandler(userService)

	users := router.Group("/users")
	{
		// 需要认证的路由
		current := users.Group("/current").Use(authMiddleware)
		{
			current.GET("", handler.GetCurrentUser)
		}
	}
}
