package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sinker/ssop/internal/models"
	"github.com/sinker/ssop/internal/repository"
	"github.com/sinker/ssop/pkg/logger"
	"github.com/sinker/ssop/pkg/utils"
)

// AnalysisService 分析功能服务接口
type AnalysisService interface {
	// 任务管理
	CreateTask(task *models.AnalysisTask) (string, error)
	GetTaskByID(id string) (*models.AnalysisTask, error)
	ListTasks(page, size int, userID string, status string) ([]*models.AnalysisTask, int64, error)
	UpdateTask(task *models.AnalysisTask) error
	DeleteTask(id string) error
	
	// 特定分析功能
	GetTemperatureSalinityTimeSeries(datasetID, lat, lng, depth, startDate, endDate, interval string) (map[string]interface{}, error)
	GetTemperatureSalinitySpatial(datasetID, date, depth, bounds, resolution string) (map[string]interface{}, error)
	
	// 结果管理
	CreateResult(result *models.AnalysisResult) (string, error)
	GetResultByID(id string) (*models.AnalysisResult, error)
	ListResultsByTaskID(taskID string) ([]*models.AnalysisResult, error)
	DeleteResult(id string) error
}

// analysisService 分析功能服务实现
type analysisService struct {
	analysisRepo repository.AnalysisRepository
	datasetRepo  repository.DatasetRepository
	resultsDir   string // 分析结果存储目录
}

// NewAnalysisService 创建分析功能服务
func NewAnalysisService(
	analysisRepo repository.AnalysisRepository,
	datasetRepo repository.DatasetRepository,
	resultsDir string,
) AnalysisService {
	// 确保结果目录存在
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		logger.Error("Failed to create analysis results directory", "error", err)
	}
	
	return &analysisService{
		analysisRepo: analysisRepo,
		datasetRepo:  datasetRepo,
		resultsDir:   resultsDir,
	}
}

// CreateTask 创建分析任务
func (s *analysisService) CreateTask(task *models.AnalysisTask) (string, error) {
	// 生成唯一ID
	if task.ID == "" {
		task.ID = utils.GenerateID("task")
	}
	
	// 设置初始状态
	if task.Status == "" {
		task.Status = "pending"
	}
	
	// 保存任务
	if err := s.analysisRepo.CreateTask(task); err != nil {
		return "", fmt.Errorf("failed to create analysis task: %w", err)
	}
	
	// 在实际应用中，这里应该启动异步任务处理
	// 为了演示，我们在这里简化处理
	go s.processTask(task)
	
	return task.ID, nil
}

// 处理分析任务
func (s *analysisService) processTask(task *models.AnalysisTask) {
	// 更新任务状态为运行中
	task.Status = "running"
	task.Progress = 10
	startTime := time.Now()
	task.StartedAt = &startTime
	
	if err := s.analysisRepo.UpdateTask(task); err != nil {
		logger.Error("Failed to update task status", "error", err, "taskId", task.ID)
		return
	}
	
	// 创建结果目录
	resultDir := filepath.Join(s.resultsDir, task.ID)
	if err := os.MkdirAll(resultDir, 0755); err != nil {
		logger.Error("Failed to create result directory", "error", err, "path", resultDir)
		task.Status = "failed"
		task.ErrorMsg = "Failed to create result directory: " + err.Error()
		s.analysisRepo.UpdateTask(task)
		return
	}
	
	// 更新进度
	task.Progress = 30
	if err := s.analysisRepo.UpdateTask(task); err != nil {
		logger.Error("Failed to update task progress", "error", err, "taskId", task.ID)
	}
	
	// 解析任务参数
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(task.Parameters), &params); err != nil {
		logger.Error("Failed to parse task parameters", "error", err, "taskId", task.ID)
		task.Status = "failed"
		task.ErrorMsg = "Invalid parameters: " + err.Error()
		s.analysisRepo.UpdateTask(task)
		return
	}
	
	// 根据任务类型执行分析
	var result map[string]interface{}
	var err error
	
	switch task.Type {
	case "temperature-salinity-timeseries":
		result, err = s.executeTemperatureSalinityTimeSeries(params)
	case "temperature-salinity-spatial":
		result, err = s.executeTemperatureSalinitySpatial(params)
	default:
		err = fmt.Errorf("unsupported analysis type: %s", task.Type)
	}
	
	// 处理分析结果
	if err != nil {
		logger.Error("Analysis task failed", "error", err, "taskId", task.ID)
		task.Status = "failed"
		task.ErrorMsg = err.Error()
		s.analysisRepo.UpdateTask(task)
		return
	}
	
	// 更新进度
	task.Progress = 70
	if err := s.analysisRepo.UpdateTask(task); err != nil {
		logger.Error("Failed to update task progress", "error", err, "taskId", task.ID)
	}
	
	// 保存结果
	resultFilePath := filepath.Join(resultDir, "result.json")
	resultData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		logger.Error("Failed to serialize result", "error", err, "taskId", task.ID)
		task.Status = "failed"
		task.ErrorMsg = "Failed to serialize result: " + err.Error()
		s.analysisRepo.UpdateTask(task)
		return
	}
	
	if err := os.WriteFile(resultFilePath, resultData, 0644); err != nil {
		logger.Error("Failed to write result file", "error", err, "taskId", task.ID)
		task.Status = "failed"
		task.ErrorMsg = "Failed to write result file: " + err.Error()
		s.analysisRepo.UpdateTask(task)
		return
	}
	
	// 创建分析结果记录
	analysisResult := &models.AnalysisResult{
		ID:          utils.GenerateID("result"),
		TaskID:      task.ID,
		Title:       task.Name + " Result",
		Description: "Result for " + task.Name,
		Type:        "json",
		Format:      "json",
		FilePath:    resultFilePath,
		PreviewData: string(resultData[:min(1000, len(resultData))]), // 保存结果预览(最多1000字节)
	}
	
	if _, err := s.CreateResult(analysisResult); err != nil {
		logger.Error("Failed to create result record", "error", err, "taskId", task.ID)
	}
	
	// 完成任务
	task.Progress = 100
	task.Status = "completed"
	task.ResultPath = resultFilePath
	completedTime := time.Now()
	task.CompletedAt = &completedTime
	
	if err := s.analysisRepo.UpdateTask(task); err != nil {
		logger.Error("Failed to update task status", "error", err, "taskId", task.ID)
	}
}

// 辅助函数，返回最小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// 执行温盐时间序列分析
func (s *analysisService) executeTemperatureSalinityTimeSeries(params map[string]interface{}) (map[string]interface{}, error) {
	// 在实际应用中，这里应该实现真正的分析逻辑
	// 这里仅做简单演示，返回模拟数据
	
	// 检查必要参数
	datasetID, _ := params["datasetId"].(string)
	if datasetID == "" {
		return nil, errors.New("missing required parameter: datasetId")
	}
	
	// 获取数据集信息
	_ , err := s.datasetRepo.GetByID(datasetID)
	if err != nil {
		return nil, fmt.Errorf("dataset not found: %w", err)
	}
	
	// 模拟时间序列数据
	lat, _ := params["lat"].(string)
	lng, _ := params["lng"].(string)
	depth, _ := params["depth"].(string)
	startDate, _ := params["startDate"].(string)
	endDate, _ := params["endDate"].(string)
	interval, _ := params["interval"].(string)
	
	// 生成模拟结果
	series := []map[string]interface{}{}
	
	// 解析起始时间和结束时间
	start, _ := time.Parse(time.RFC3339, startDate)
	end, _ := time.Parse(time.RFC3339, endDate)
	
	// 根据间隔设置步长
	var step time.Duration
	switch interval {
	case "hour":
		step = time.Hour
	case "day":
		step = 24 * time.Hour
	case "week":
		step = 7 * 24 * time.Hour
	case "month":
		step = 30 * 24 * time.Hour
	default:
		step = 24 * time.Hour
	}
	
	// 生成数据点
	for t := start; t.Before(end) || t.Equal(end); t = t.Add(step) {
		// 模拟温度和盐度数据
		temp := 20.0 + 5.0*utils.RandomFloat64(-1, 1) // 20±5°C
		salinity := 35.0 + 2.0*utils.RandomFloat64(-1, 1) // 35±2 PSU
		
		series = append(series, map[string]interface{}{
			"timestamp":   t.Format(time.RFC3339),
			"temperature": temp,
			"salinity":    salinity,
		})
	}
	
	return map[string]interface{}{
		"location": map[string]interface{}{
			"lat":   lat,
			"lng":   lng,
			"depth": depth,
		},
		"timeRange": map[string]interface{}{
			"start": startDate,
			"end":   endDate,
		},
		"interval": interval,
		"series":   series,
	}, nil
}

// 执行温盐空间分布分析
func (s *analysisService) executeTemperatureSalinitySpatial(params map[string]interface{}) (map[string]interface{}, error) {
	// 在实际应用中，这里应该实现真正的分析逻辑
	// 这里仅做简单演示，返回模拟数据
	
	// 检查必要参数
	datasetID, _ := params["datasetId"].(string)
	if datasetID == "" {
		return nil, errors.New("missing required parameter: datasetId")
	}
	
	// 获取数据集信息
	_, err := s.datasetRepo.GetByID(datasetID)
	if err != nil {
		return nil, fmt.Errorf("dataset not found: %w", err)
	}
	
	// 提取参数
	date, _ := params["date"].(string)
	depth, _ := params["depth"].(string)
	bounds, _ := params["bounds"].(string)
	resolution, _ := params["resolution"].(string)
	
	// 解析边界范围
	var minLat, minLng, maxLat, maxLng float64 = 20.0, 110.0, 25.0, 120.0
	fmt.Sscanf(bounds, "%f,%f,%f,%f", &minLat, &minLng, &maxLat, &maxLng)
	
	// 计算网格大小
	var latStep, lngStep float64
	var latCount, lngCount int
	
	switch resolution {
	case "low":
		latStep, lngStep = 1.0, 1.0
	case "medium":
		latStep, lngStep = 0.5, 0.5
	case "high":
		latStep, lngStep = 0.1, 0.1
	default:
		latStep, lngStep = 0.5, 0.5
	}
	
	latCount = int((maxLat - minLat) / latStep) + 1
	lngCount = int((maxLng - minLng) / lngStep) + 1
	
	// 生成温度和盐度数据
	temperatureData := make([][]float64, latCount)
	salinityData := make([][]float64, latCount)
	
	for i := 0; i < latCount; i++ {
		temperatureData[i] = make([]float64, lngCount)
		salinityData[i] = make([]float64, lngCount)
		
		for j := 0; j < lngCount; j++ {
			// 模拟温度和盐度空间分布
			lat := minLat + float64(i)*latStep
			// lng := minLng + float64(j)*lngStep
			
			// 简单模拟: 温度随纬度降低，盐度随纬度升高
			temp := 30.0 - 0.2*(lat-minLat) + 2.0*utils.RandomFloat64(-1, 1)
			salt := 33.0 + 0.1*(lat-minLat) + 1.0*utils.RandomFloat64(-1, 1)
			
			temperatureData[i][j] = temp
			salinityData[i][j] = salt
		}
	}
	
	return map[string]interface{}{
		"time":   date,
		"depth":  depth,
		"bounds": []float64{minLat, minLng, maxLat, maxLng},
		"resolution": resolution,
		"grid": map[string]interface{}{
			"latCount": latCount,
			"lngCount": lngCount,
			"latStep":  latStep,
			"lngStep":  lngStep,
			"startLat": minLat,
			"startLng": minLng,
		},
		"data": map[string]interface{}{
			"temperature": temperatureData,
			"salinity":    salinityData,
		},
	}, nil
}

// GetTaskByID 根据ID获取分析任务
func (s *analysisService) GetTaskByID(id string) (*models.AnalysisTask, error) {
	return s.analysisRepo.GetTaskByID(id)
}

// ListTasks 获取分析任务列表
func (s *analysisService) ListTasks(page, size int, userID string, status string) ([]*models.AnalysisTask, int64, error) {
	return s.analysisRepo.ListTasks(page, size, userID, status)
}

// UpdateTask 更新分析任务
func (s *analysisService) UpdateTask(task *models.AnalysisTask) error {
	// 获取原始任务信息
	originalTask, err := s.analysisRepo.GetTaskByID(task.ID)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}
	
	// 保留不可修改的字段
	task.CreatedBy = originalTask.CreatedBy
	task.CreatedAt = originalTask.CreatedAt
	task.StartedAt = originalTask.StartedAt
	task.CompletedAt = originalTask.CompletedAt
	
	return s.analysisRepo.UpdateTask(task)
}

// DeleteTask 删除分析任务
func (s *analysisService) DeleteTask(id string) error {
	// 获取任务信息
	_, err := s.analysisRepo.GetTaskByID(id)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}
	
	// 删除关联结果
	results, err := s.analysisRepo.ListResultsByTaskID(id)
	if err != nil {
		logger.Error("Failed to get results for task", "error", err, "taskId", id)
	} else {
		for _, result := range results {
			if err := s.DeleteResult(result.ID); err != nil {
				logger.Error("Failed to delete result", "error", err, "resultId", result.ID)
			}
		}
	}
	
	// 删除结果目录
	resultDir := filepath.Join(s.resultsDir, id)
	if err := os.RemoveAll(resultDir); err != nil && !os.IsNotExist(err) {
		logger.Error("Failed to delete result directory", "error", err, "path", resultDir)
	}
	
	// 从数据库中删除任务
	return s.analysisRepo.DeleteTask(id)
}

// GetTemperatureSalinityTimeSeries 获取温盐时间序列
func (s *analysisService) GetTemperatureSalinityTimeSeries(datasetID, lat, lng, depth, startDate, endDate, interval string) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"datasetId": datasetID,
		"lat":       lat,
		"lng":       lng,
		"depth":     depth,
		"startDate": startDate,
		"endDate":   endDate,
		"interval":  interval,
	}
	
	return s.executeTemperatureSalinityTimeSeries(params)
}

// GetTemperatureSalinitySpatial 获取温盐空间分布
func (s *analysisService) GetTemperatureSalinitySpatial(datasetID, date, depth, bounds, resolution string) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"datasetId":  datasetID,
		"date":       date,
		"depth":      depth,
		"bounds":     bounds,
		"resolution": resolution,
	}
	
	return s.executeTemperatureSalinitySpatial(params)
}

// CreateResult 创建分析结果
func (s *analysisService) CreateResult(result *models.AnalysisResult) (string, error) {
	// 生成唯一ID
	if result.ID == "" {
		result.ID = utils.GenerateID("result")
	}
	
	// 保存结果
	if err := s.analysisRepo.CreateResult(result); err != nil {
		return "", fmt.Errorf("failed to create analysis result: %w", err)
	}
	
	return result.ID, nil
}

// GetResultByID 根据ID获取分析结果
func (s *analysisService) GetResultByID(id string) (*models.AnalysisResult, error) {
	return s.analysisRepo.GetResultByID(id)
}

// ListResultsByTaskID 根据任务ID获取分析结果列表
func (s *analysisService) ListResultsByTaskID(taskID string) ([]*models.AnalysisResult, error) {
	return s.analysisRepo.ListResultsByTaskID(taskID)
}

// DeleteResult 删除分析结果
func (s *analysisService) DeleteResult(id string) error {
	// 获取结果信息
	result, err := s.analysisRepo.GetResultByID(id)
	if err != nil {
		return fmt.Errorf("result not found: %w", err)
	}
	
	// 删除结果文件
	if result.FilePath != "" {
		if err := os.Remove(result.FilePath); err != nil && !os.IsNotExist(err) {
			logger.Error("Failed to delete result file", "error", err, "path", result.FilePath)
		}
	}
	
	// 从数据库中删除
	return s.analysisRepo.DeleteResult(id)
}