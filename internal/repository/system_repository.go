package repository

import (
	"github.com/sinker/ssop/internal/models"
	"gorm.io/gorm"
)

// SystemRepository 系统仓库接口
type SystemRepository interface {
	// 系统设置
	GetSetting(key string) (*models.SystemSetting, error)
	GetSettingsByCategory(category string) ([]*models.SystemSetting, error)
	UpdateSetting(setting *models.SystemSetting) error
	
	// 操作日志
	CreateAuditLog(log *models.AuditLog) error
	ListAuditLogs(page, size int, filters map[string]interface{}) ([]*models.AuditLog, int64, error)
}

// systemRepository 系统仓库实现
type systemRepository struct {
	db *gorm.DB
}

// NewSystemRepository 创建系统仓库
func NewSystemRepository(db *gorm.DB) SystemRepository {
	return &systemRepository{db: db}
}

// GetSetting 获取系统设置
func (r *systemRepository) GetSetting(key string) (*models.SystemSetting, error) {
	var setting models.SystemSetting
	err := r.db.Where("key = ?", key).First(&setting).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

// GetSettingsByCategory 获取特定类别的系统设置
func (r *systemRepository) GetSettingsByCategory(category string) ([]*models.SystemSetting, error) {
	var settings []*models.SystemSetting
	err := r.db.Where("category = ?", category).Find(&settings).Error
	if err != nil {
		return nil, err
	}
	return settings, nil
}

// UpdateSetting 更新系统设置
func (r *systemRepository) UpdateSetting(setting *models.SystemSetting) error {
	// 使用Upsert操作，如果存在则更新，不存在则创建
	return r.db.Save(setting).Error
}

// CreateAuditLog 创建操作日志
func (r *systemRepository) CreateAuditLog(log *models.AuditLog) error {
	return r.db.Create(log).Error
}

// ListAuditLogs 获取操作日志列表
func (r *systemRepository) ListAuditLogs(page, size int, filters map[string]interface{}) ([]*models.AuditLog, int64, error) {
	var logs []*models.AuditLog
	var total int64

	query := r.db.Model(&models.AuditLog{})

	// 应用过滤条件
	if filters != nil {
		if userID, ok := filters["userId"]; ok && userID != "" {
			query = query.Where("user_id = ?", userID)
		}
		if action, ok := filters["action"]; ok && action != "" {
			query = query.Where("action = ?", action)
		}
		if resource, ok := filters["resource"]; ok && resource != "" {
			query = query.Where("resource = ?", resource)
		}
		if startTime, ok := filters["startTime"]; ok && startTime != "" {
			query = query.Where("created_at >= ?", startTime)
		}
		if endTime, ok := filters["endTime"]; ok && endTime != "" {
			query = query.Where("created_at <= ?", endTime)
		}
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
	err := query.Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
} 