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

// RegisterAnalysisRoutes 注册分析功能相关路由
func RegisterAnalysisRoutes(router *gin.RouterGroup, analysisService services.AnalysisService, authMiddleware gin.HandlerFunc) {
	analysisHandler := &AnalysisHandler{analysisService: analysisService}
	
	// 需要认证的接口
	analysis := router.Group("/analysis")
	analysis.Use(authMiddleware)
	{
		// 任务管理
		tasks := analysis.Group("/tasks")
		{
			tasks.POST("", analysisHandler.CreateTask)
			tasks.GET("", analysisHandler.ListTasks)
			tasks.GET("/:taskId", analysisHandler.GetTaskByID)
			tasks.PUT("/:taskId", analysisHandler.UpdateTask)
			tasks.DELETE("/:taskId", analysisHandler.DeleteTask)
			tasks.GET("/:taskId/results", analysisHandler.ListResultsByTaskID)
		}
		
		// 结果管理
		results := analysis.Group("/results")
		{
			results.GET("/:resultId", analysisHandler.GetResultByID)
			results.DELETE("/:resultId", analysisHandler.DeleteResult)
		}
		
		// 温盐分析
		ts := analysis.Group("/temperature-salinity")
		{
			ts.GET("/timeseries", analysisHandler.GetTemperatureSalinityTimeSeries)
			ts.GET("/spatial", analysisHandler.GetTemperatureSalinitySpatial)
		}
	}
}

// AnalysisHandler 分析功能处理器
type AnalysisHandler struct {
	analysisService services.AnalysisService
}

// CreateTask 创建分析任务
func (h *AnalysisHandler) CreateTask(c *gin.Context) {
	var task models.AnalysisTask
	if err := c.ShouldBindJSON(&task); err != nil {
		logger.Error("Failed to parse request body", "error", err)
		response.Fail(c, http.StatusBadRequest, "请求格式错误")
		return
	}
	
	// 设置创建者
	userID, _ := c.Get("userId")
	task.CreatedBy = userID.(string)
	
	// 创建任务
	taskID, err := h.analysisService.CreateTask(&task)
	if err != nil {
		logger.Error("Failed to create analysis task", "error", err)
		response.Fail(c, http.StatusInternalServerError, "创建分析任务失败")
		return
	}
	
	response.Success(c, gin.H{
		"taskId":   taskID,
		"status":   "pending",
		"progress": 0,
	}, "创建成功")
}

// ListTasks 获取分析任务列表
func (h *AnalysisHandler) ListTasks(c *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	
	// 获取筛选条件
	status := c.Query("status")
	
	// 根据当前用户获取任务列表
	userID, _ := c.Get("userId")
	
	tasks, total, err := h.analysisService.ListTasks(page, size, userID.(string), status)
	if err != nil {
		logger.Error("Failed to get analysis tasks", "error", err)
		response.Fail(c, http.StatusInternalServerError, "获取分析任务列表失败")
		return
	}
	
	response.Success(c, gin.H{
		"total": total,
		"page":  page,
		"size":  size,
		"tasks": tasks,
	}, "获取成功")
}

// GetTaskByID 获取分析任务详情
func (h *AnalysisHandler) GetTaskByID(c *gin.Context) {
	taskID := c.Param("taskId")
	
	task, err := h.analysisService.GetTaskByID(taskID)
	if err != nil {
		logger.Error("Failed to get analysis task", "error", err, "taskId", taskID)
		response.Fail(c, http.StatusNotFound, "分析任务不存在")
		return
	}
	
	// 检查权限(只能查看自己的任务)
	userID, _ := c.Get("userId")
	if task.CreatedBy != userID.(string) {
		response.Fail(c, http.StatusForbidden, "无权访问此任务")
		return
	}
	
	response.Success(c, task, "获取成功")
}

// UpdateTask 更新分析任务
func (h *AnalysisHandler) UpdateTask(c *gin.Context) {
	taskID := c.Param("taskId")
	
	// 获取原任务信息
	originalTask, err := h.analysisService.GetTaskByID(taskID)
	if err != nil {
		logger.Error("Failed to get analysis task", "error", err, "taskId", taskID)
		response.Fail(c, http.StatusNotFound, "分析任务不存在")
		return
	}
	
	// 检查权限(只能更新自己的任务)
	userID, _ := c.Get("userId")
	if originalTask.CreatedBy != userID.(string) {
		response.Fail(c, http.StatusForbidden, "无权更新此任务")
		return
	}
	
	// 解析请求体
	var updateData models.AnalysisTask
	if err := c.ShouldBindJSON(&updateData); err != nil {
		logger.Error("Failed to parse request body", "error", err)
		response.Fail(c, http.StatusBadRequest, "请求格式错误")
		return
	}
	
	// 设置不可变更的字段
	updateData.ID = taskID
	
	// 更新任务
	if err := h.analysisService.UpdateTask(&updateData); err != nil {
		logger.Error("Failed to update analysis task", "error", err, "taskId", taskID)
		response.Fail(c, http.StatusInternalServerError, "更新分析任务失败")
		return
	}
	
	response.Success(c, gin.H{"message": "更新成功"}, "更新成功")
}

// DeleteTask 删除分析任务
func (h *AnalysisHandler) DeleteTask(c *gin.Context) {
	taskID := c.Param("taskId")
	
	// 获取任务信息
	task, err := h.analysisService.GetTaskByID(taskID)
	if err != nil {
		logger.Error("Failed to get analysis task", "error", err, "taskId", taskID)
		response.Fail(c, http.StatusNotFound, "分析任务不存在")
		return
	}
	
	// 检查权限(只能删除自己的任务)
	userID, _ := c.Get("userId")
	if task.CreatedBy != userID.(string) {
		response.Fail(c, http.StatusForbidden, "无权删除此任务")
		return
	}
	
	// 删除任务
	if err := h.analysisService.DeleteTask(taskID); err != nil {
		logger.Error("Failed to delete analysis task", "error", err, "taskId", taskID)
		response.Fail(c, http.StatusInternalServerError, "删除分析任务失败")
		return
	}
	
	response.Success(c, gin.H{"message": "删除成功"}, "删除成功")
}

// ListResultsByTaskID 获取任务的结果列表
func (h *AnalysisHandler) ListResultsByTaskID(c *gin.Context) {
	taskID := c.Param("taskId")
	
	// 获取任务信息
	task, err := h.analysisService.GetTaskByID(taskID)
	if err != nil {
		logger.Error("Failed to get analysis task", "error", err, "taskId", taskID)
		response.Fail(c, http.StatusNotFound, "分析任务不存在")
		return
	}
	
	// 检查权限(只能查看自己的任务结果)
	userID, _ := c.Get("userId")
	if task.CreatedBy != userID.(string) {
		response.Fail(c, http.StatusForbidden, "无权访问此任务结果")
		return
	}
	
	// 获取结果列表
	results, err := h.analysisService.ListResultsByTaskID(taskID)
	if err != nil {
		logger.Error("Failed to get results", "error", err, "taskId", taskID)
		response.Fail(c, http.StatusInternalServerError, "获取分析结果失败")
		return
	}
	
	response.Success(c, gin.H{
		"taskId":  taskID,
		"results": results,
	}, "获取成功")
}

// GetResultByID 获取分析结果详情
func (h *AnalysisHandler) GetResultByID(c *gin.Context) {
	resultID := c.Param("resultId")
	
	// 获取结果信息
	result, err := h.analysisService.GetResultByID(resultID)
	if err != nil {
		logger.Error("Failed to get result", "error", err, "resultId", resultID)
		response.Fail(c, http.StatusNotFound, "分析结果不存在")
		return
	}
	
	// 获取任务信息以验证权限
	task, err := h.analysisService.GetTaskByID(result.TaskID)
	if err != nil {
		logger.Error("Failed to get task for result", "error", err, "taskId", result.TaskID)
		response.Fail(c, http.StatusInternalServerError, "获取关联任务失败")
		return
	}
	
	// 检查权限(只能查看自己任务的结果)
	userID, _ := c.Get("userId")
	if task.CreatedBy != userID.(string) {
		response.Fail(c, http.StatusForbidden, "无权访问此分析结果")
		return
	}
	
	response.Success(c, result, "获取成功")
}

// DeleteResult 删除分析结果
func (h *AnalysisHandler) DeleteResult(c *gin.Context) {
	resultID := c.Param("resultId")
	
	// 获取结果信息
	result, err := h.analysisService.GetResultByID(resultID)
	if err != nil {
		logger.Error("Failed to get result", "error", err, "resultId", resultID)
		response.Fail(c, http.StatusNotFound, "分析结果不存在")
		return
	}
	
	// 获取任务信息以验证权限
	task, err := h.analysisService.GetTaskByID(result.TaskID)
	if err != nil {
		logger.Error("Failed to get task for result", "error", err, "taskId", result.TaskID)
		response.Fail(c, http.StatusInternalServerError, "获取关联任务失败")
		return
	}
	
	// 检查权限(只能删除自己任务的结果)
	userID, _ := c.Get("userId")
	if task.CreatedBy != userID.(string) {
		response.Fail(c, http.StatusForbidden, "无权删除此分析结果")
		return
	}
	
	// 删除结果
	if err := h.analysisService.DeleteResult(resultID); err != nil {
		logger.Error("Failed to delete result", "error", err, "resultId", resultID)
		response.Fail(c, http.StatusInternalServerError, "删除分析结果失败")
		return
	}
	
	response.Success(c, gin.H{"message": "删除成功"}, "删除成功")
}

// GetTemperatureSalinityTimeSeries 获取温盐时间序列
func (h *AnalysisHandler) GetTemperatureSalinityTimeSeries(c *gin.Context) {
	// 获取参数
	datasetID := c.Query("datasetId")
	lat := c.Query("lat")
	lng := c.Query("lng")
	depth := c.DefaultQuery("depth", "0")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	interval := c.DefaultQuery("interval", "day")
	
	// 验证必要参数
	if datasetID == "" || lat == "" || lng == "" || startDate == "" || endDate == "" {
		response.Fail(c, http.StatusBadRequest, "缺少必要参数")
		return
	}
	
	// 执行分析
	result, err := h.analysisService.GetTemperatureSalinityTimeSeries(datasetID, lat, lng, depth, startDate, endDate, interval)
	if err != nil {
		logger.Error("Failed to get temperature-salinity timeseries", "error", err)
		response.Fail(c, http.StatusInternalServerError, "获取温盐时间序列失败: "+err.Error())
		return
	}
	
	response.Success(c, result, "获取成功")
}

// GetTemperatureSalinitySpatial 获取温盐空间分布
func (h *AnalysisHandler) GetTemperatureSalinitySpatial(c *gin.Context) {
	// 获取参数
	datasetID := c.Query("datasetId")
	date := c.Query("date")
	depth := c.DefaultQuery("depth", "0")
	bounds := c.Query("bounds")
	resolution := c.DefaultQuery("resolution", "medium")
	
	// 验证必要参数
	if datasetID == "" || date == "" || bounds == "" {
		response.Fail(c, http.StatusBadRequest, "缺少必要参数")
		return
	}
	
	// 执行分析
	result, err := h.analysisService.GetTemperatureSalinitySpatial(datasetID, date, depth, bounds, resolution)
	if err != nil {
		logger.Error("Failed to get temperature-salinity spatial distribution", "error", err)
		response.Fail(c, http.StatusInternalServerError, "获取温盐空间分布失败: "+err.Error())
		return
	}
	
	response.Success(c, result, "获取成功")
} 