package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sinker/ssop/internal/models"
	"github.com/sinker/ssop/internal/services"
	"github.com/sinker/ssop/pkg/logger"
	"github.com/sinker/ssop/pkg/response"
)

// RegisterDatasetRoutes 注册数据集相关路由
func RegisterDatasetRoutes(router *gin.RouterGroup, datasetService services.DatasetService, authMiddleware gin.HandlerFunc) {
	datasetHandler := &DatasetHandler{datasetService: datasetService}
	
	datasets := router.Group("/datasets")
	{
		// 公开接口
		datasets.GET("", datasetHandler.GetDatasets)
		datasets.GET("/:datasetId", datasetHandler.GetDatasetByID)
		
		// 需要认证的接口
		authenticated := datasets.Group("")
		authenticated.Use(authMiddleware)
		{
			authenticated.POST("/upload", datasetHandler.UploadDataset)
			authenticated.PUT("/:datasetId", datasetHandler.UpdateDataset)
			authenticated.DELETE("/:datasetId", datasetHandler.DeleteDataset)
			authenticated.GET("/:datasetId/download", datasetHandler.DownloadDataset)
		}
	}
}

// DatasetHandler 数据集处理器
type DatasetHandler struct {
	datasetService services.DatasetService
}

// GetDatasets 获取数据集列表
func (h *DatasetHandler) GetDatasets(c *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	
	// 解析过滤参数
	filters := make(map[string]interface{})
	
	if dataType := c.Query("type"); dataType != "" {
		filters["type"] = dataType
	}
	
	if startDate := c.Query("startDate"); startDate != "" {
		filters["startDate"] = startDate
	}
	
	if endDate := c.Query("endDate"); endDate != "" {
		filters["endDate"] = endDate
	}
	
	if region := c.Query("region"); region != "" {
		filters["region"] = region
	}
	
	if keyword := c.Query("keyword"); keyword != "" {
		filters["keyword"] = keyword
	}
	
	// 获取数据集列表
	datasets, total, err := h.datasetService.GetDatasets(page, size, filters)
	if err != nil {
		logger.Error("Failed to get datasets", "error", err)
		response.Fail(c, http.StatusInternalServerError, "获取数据集列表失败")
		return
	}
	
	// 构造响应
	response.Success(c, gin.H{
		"total":    total,
		"page":     page,
		"size":     size,
		"datasets": datasets,
	}, "获取成功")
}

// GetDatasetByID 获取数据集详情
func (h *DatasetHandler) GetDatasetByID(c *gin.Context) {
	datasetID := c.Param("datasetId")
	
	dataset, err := h.datasetService.GetDatasetByID(datasetID)
	if err != nil {
		logger.Error("Failed to get dataset", "error", err, "datasetId", datasetID)
		response.Fail(c, http.StatusNotFound, "数据集不存在")
		return
	}
	
	response.Success(c, dataset, "获取成功")
}

// UploadDataset 上传数据集
func (h *DatasetHandler) UploadDataset(c *gin.Context) {
	// 获取上传的文件
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		logger.Error("Failed to get file", "error", err)
		response.Fail(c, http.StatusBadRequest, "文件上传失败")
		return
	}
	defer file.Close()
	
	// 解析元数据
	metadataStr := c.PostForm("metadata")
	var dataset models.Dataset
	
	if err := json.Unmarshal([]byte(metadataStr), &dataset); err != nil {
		logger.Error("Failed to parse metadata", "error", err)
		response.Fail(c, http.StatusBadRequest, "元数据格式错误")
		return
	}
	
	// 设置创建者
	userID, _ := c.Get("userId")
	dataset.CreatedBy = userID.(string)
	
	// 创建数据集
	datasetID, err := h.datasetService.CreateDataset(&dataset, file, fileHeader.Filename)
	if err != nil {
		logger.Error("Failed to create dataset", "error", err)
		response.Fail(c, http.StatusInternalServerError, "创建数据集失败")
		return
	}
	
	// 返回数据集ID
	response.Success(c, gin.H{
		"datasetId":  datasetID,
		"name":       dataset.Name,
		"uploadTime": dataset.CreatedAt,
		"size":       dataset.Size,
		"status":     "processing", // 实际应用中可能需要处理数据集的状态
	}, "上传成功")
}

// UpdateDataset 更新数据集
func (h *DatasetHandler) UpdateDataset(c *gin.Context) {
	datasetID := c.Param("datasetId")
	
	// 获取请求体
	var updateData models.Dataset
	if err := c.ShouldBindJSON(&updateData); err != nil {
		logger.Error("Failed to parse request body", "error", err)
		response.Fail(c, http.StatusBadRequest, "请求格式错误")
		return
	}
	
	// 设置ID
	updateData.ID = datasetID
	
	// 更新数据集
	if err := h.datasetService.UpdateDataset(&updateData); err != nil {
		logger.Error("Failed to update dataset", "error", err, "datasetId", datasetID)
		response.Fail(c, http.StatusInternalServerError, "更新数据集失败")
		return
	}
	
	response.Success(c, gin.H{"message": "更新成功"}, "更新成功")
}

// DeleteDataset 删除数据集
func (h *DatasetHandler) DeleteDataset(c *gin.Context) {
	datasetID := c.Param("datasetId")
	
	// 删除数据集
	if err := h.datasetService.DeleteDataset(datasetID); err != nil {
		logger.Error("Failed to delete dataset", "error", err, "datasetId", datasetID)
		response.Fail(c, http.StatusInternalServerError, "删除数据集失败")
		return
	}
	
	response.Success(c, gin.H{"message": "删除成功"}, "删除成功")
}

// DownloadDataset 下载数据集
func (h *DatasetHandler) DownloadDataset(c *gin.Context) {
	datasetID := c.Param("datasetId")
	
	// 获取数据集文件路径
	filePath, err := h.datasetService.DownloadDataset(datasetID)
	if err != nil {
		logger.Error("Failed to get dataset file", "error", err, "datasetId", datasetID)
		response.Fail(c, http.StatusNotFound, "数据集文件不存在")
		return
	}
	
	// 获取文件名
	parts := strings.Split(filePath, "/")
	fileName := parts[len(parts)-1]
	
	// 设置下载头
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Description", "File Transfer")
	c.File(filePath)
} 