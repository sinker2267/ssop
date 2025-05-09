package repository

import (
	_ "time"

	"github.com/sinker/ssop/internal/models"
	"gorm.io/gorm"
)

// AnalysisRepository 分析功能仓库接口
type AnalysisRepository interface {
	// 任务相关
	CreateTask(task *models.AnalysisTask) error
	GetTaskByID(id string) (*models.AnalysisTask, error)
	ListTasks(page, size int, userID string, status string) ([]*models.AnalysisTask, int64, error)
	UpdateTask(task *models.AnalysisTask) error
	DeleteTask(id string) error
	
	// 结果相关
	CreateResult(result *models.AnalysisResult) error
	GetResultByID(id string) (*models.AnalysisResult, error)
	ListResultsByTaskID(taskID string) ([]*models.AnalysisResult, error)
	DeleteResult(id string) error
}

// analysisRepository 分析功能仓库实现
type analysisRepository struct {
	db *gorm.DB
}

// NewAnalysisRepository 创建分析功能仓库
func NewAnalysisRepository(db *gorm.DB) AnalysisRepository {
	return &analysisRepository{db: db}
}

// CreateTask 创建分析任务
func (r *analysisRepository) CreateTask(task *models.AnalysisTask) error {
	return r.db.Create(task).Error
}

// GetTaskByID 根据ID获取分析任务
func (r *analysisRepository) GetTaskByID(id string) (*models.AnalysisTask, error) {
	var task models.AnalysisTask
	err := r.db.Where("id = ?", id).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// ListTasks 获取分析任务列表
func (r *analysisRepository) ListTasks(page, size int, userID string, status string) ([]*models.AnalysisTask, int64, error) {
	var tasks []*models.AnalysisTask
	var total int64

	query := r.db.Model(&models.AnalysisTask{})

	// 应用过滤条件
	if userID != "" {
		query = query.Where("created_by = ?", userID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 计算总数
	query.Count(&total)

	// 分页
	if page > 0 && size > 0 {
		offset := (page - 1) * size
		query = query.Offset(offset).Limit(size)
	}

	// 排序(默认按创建时间倒序)
	query = query.Order("created_at DESC")

	// 执行查询
	err := query.Find(&tasks).Error
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// UpdateTask 更新分析任务
func (r *analysisRepository) UpdateTask(task *models.AnalysisTask) error {
	return r.db.Save(task).Error
}

// DeleteTask 删除分析任务
func (r *analysisRepository) DeleteTask(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.AnalysisTask{}).Error
}

// CreateResult 创建分析结果
func (r *analysisRepository) CreateResult(result *models.AnalysisResult) error {
	return r.db.Create(result).Error
}

// GetResultByID 根据ID获取分析结果
func (r *analysisRepository) GetResultByID(id string) (*models.AnalysisResult, error) {
	var result models.AnalysisResult
	err := r.db.Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListResultsByTaskID 根据任务ID获取分析结果列表
func (r *analysisRepository) ListResultsByTaskID(taskID string) ([]*models.AnalysisResult, error) {
	var results []*models.AnalysisResult
	err := r.db.Where("task_id = ?", taskID).Order("created_at DESC").Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

// DeleteResult 删除分析结果
func (r *analysisRepository) DeleteResult(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.AnalysisResult{}).Error
} 