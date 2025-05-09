package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sinker/ssop/internal/models"
	"github.com/sinker/ssop/internal/services"
	"github.com/sinker/ssop/pkg/logger"
	"github.com/sinker/ssop/pkg/response"
)

// RegisterSystemRoutes 注册系统相关路由
func RegisterSystemRoutes(router *gin.RouterGroup, systemService services.SystemService, authMiddleware gin.HandlerFunc) {
	systemHandler := &SystemHandler{systemService: systemService}
	
	// 系统管理路由(需要管理员权限)
	system := router.Group("/system")
	system.Use(authMiddleware, AdminRequired())
	{
		// 系统设置
		settings := system.Group("/settings")
		{
			settings.GET("", systemHandler.GetSettings)
			settings.GET("/:category", systemHandler.GetSettingsByCategory)
			settings.PUT("/:key", systemHandler.UpdateSetting)
		}
		
		// 审计日志
		logs := system.Group("/logs")
		{
			logs.GET("", systemHandler.GetAuditLogs)
		}
	}
}

// SystemHandler 系统处理器
type SystemHandler struct {
	systemService services.SystemService
}

// AdminRequired 管理员权限中间件
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户角色
		role, exists := c.Get("userRole")
		if !exists {
			response.Fail(c, http.StatusUnauthorized, "未登录")
			c.Abort()
			return
		}
		
		// 检查是否为管理员
		if role != "admin" {
			response.Fail(c, http.StatusForbidden, "需要管理员权限")
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// GetSettings 获取所有系统设置
func (h *SystemHandler) GetSettings(c *gin.Context) {
	// 获取所有设置分类
	categories := []string{"general", "security", "storage", "notification", "analysis"}
	
	result := make(map[string][]*models.SystemSetting)
	
	// 获取每个分类的设置
	for _, category := range categories {
		settings, err := h.systemService.GetSettingsByCategory(category)
		if err != nil {
			logger.Error("Failed to get settings for category", "error", err, "category", category)
			continue
		}
		
		if len(settings) > 0 {
			result[category] = settings
		}
	}
	
	response.Success(c, result, "获取成功")
}

// GetSettingsByCategory 获取特定分类的设置
func (h *SystemHandler) GetSettingsByCategory(c *gin.Context) {
	category := c.Param("category")
	
	settings, err := h.systemService.GetSettingsByCategory(category)
	if err != nil {
		logger.Error("Failed to get settings", "error", err, "category", category)
		response.Fail(c, http.StatusInternalServerError, "获取系统设置失败")
		return
	}
	
	response.Success(c, settings, "获取成功")
}

// UpdateSetting 更新系统设置
func (h *SystemHandler) UpdateSetting(c *gin.Context) {
	key := c.Param("key")
	
	// 解析请求
	var req struct {
		Value       string `json:"value"`
		Category    string `json:"category"`
		Description string `json:"description"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to parse request body", "error", err)
		response.Fail(c, http.StatusBadRequest, "请求格式错误")
		return
	}
	
	// 更新设置
	if err := h.systemService.UpdateSetting(key, req.Value, req.Category, req.Description); err != nil {
		logger.Error("Failed to update setting", "error", err, "key", key)
		response.Fail(c, http.StatusInternalServerError, "更新系统设置失败")
		return
	}
	
	response.Success(c, gin.H{"message": "更新成功"}, "更新成功")
}

// GetAuditLogs 获取审计日志
func (h *SystemHandler) GetAuditLogs(c *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	
	// 解析过滤参数
	filters := make(map[string]interface{})
	
	if userID := c.Query("userId"); userID != "" {
		filters["userId"] = userID
	}
	
	if action := c.Query("action"); action != "" {
		filters["action"] = action
	}
	
	if resource := c.Query("resource"); resource != "" {
		filters["resource"] = resource
	}
	
	if startTime := c.Query("startTime"); startTime != "" {
		filters["startTime"] = startTime
	}
	
	if endTime := c.Query("endTime"); endTime != "" {
		filters["endTime"] = endTime
	}
	
	// 获取日志
	logs, total, err := h.systemService.GetAuditLogs(page, size, filters)
	if err != nil {
		logger.Error("Failed to get audit logs", "error", err)
		response.Fail(c, http.StatusInternalServerError, "获取审计日志失败")
		return
	}
	
	response.Success(c, gin.H{
		"total": total,
		"page":  page,
		"size":  size,
		"logs":  logs,
	}, "获取成功")
} 