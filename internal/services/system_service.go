package services

import (
	"errors"
	"strconv"

	"github.com/sinker/ssop/internal/models"
	"github.com/sinker/ssop/internal/repository"
	"github.com/sinker/ssop/pkg/utils"
)

// SystemService 系统服务接口
type SystemService interface {
	// 设置管理
	GetSetting(key string) (string, error)
	GetSettingInt(key string, defaultValue int) int
	GetSettingBool(key string, defaultValue bool) bool
	GetSettingsByCategory(category string) ([]*models.SystemSetting, error)
	UpdateSetting(key, value, category, description string) error
	
	// 审计日志
	CreateAuditLog(log *models.AuditLog) error
	GetAuditLogs(page, size int, filters map[string]interface{}) ([]*models.AuditLog, int64, error)
}

// systemService 系统服务实现
type systemService struct {
	systemRepo repository.SystemRepository
}

// NewSystemService 创建系统服务
func NewSystemService(systemRepo repository.SystemRepository) SystemService {
	return &systemService{
		systemRepo: systemRepo,
	}
}

// GetSetting 获取系统设置
func (s *systemService) GetSetting(key string) (string, error) {
	setting, err := s.systemRepo.GetSetting(key)
	if err != nil {
		return "", err
	}
	return setting.Value, nil
}

// GetSettingInt 获取整数类型的系统设置
func (s *systemService) GetSettingInt(key string, defaultValue int) int {
	value, err := s.GetSetting(key)
	if err != nil {
		return defaultValue
	}
	
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	
	return intValue
}

// GetSettingBool 获取布尔类型的系统设置
func (s *systemService) GetSettingBool(key string, defaultValue bool) bool {
	value, err := s.GetSetting(key)
	if err != nil {
		return defaultValue
	}
	
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	
	return boolValue
}

// GetSettingsByCategory 获取特定类别的系统设置
func (s *systemService) GetSettingsByCategory(category string) ([]*models.SystemSetting, error) {
	return s.systemRepo.GetSettingsByCategory(category)
}

// UpdateSetting 更新系统设置
func (s *systemService) UpdateSetting(key, value, category, description string) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}
	
	// 尝试获取现有设置
	existingSetting, err := s.systemRepo.GetSetting(key)
	if err == nil {
		// 更新现有设置
		existingSetting.Value = value
		if category != "" {
			existingSetting.Category = category
		}
		if description != "" {
			existingSetting.Description = description
		}
		return s.systemRepo.UpdateSetting(existingSetting)
	}
	
	// 创建新设置
	newSetting := &models.SystemSetting{
		ID:          utils.GenerateID("setting"),
		Key:         key,
		Value:       value,
		Category:    category,
		Description: description,
	}
	
	return s.systemRepo.UpdateSetting(newSetting)
}

// CreateAuditLog 创建审计日志
func (s *systemService) CreateAuditLog(log *models.AuditLog) error {
	return s.systemRepo.CreateAuditLog(log)
}

// GetAuditLogs 获取审计日志
func (s *systemService) GetAuditLogs(page, size int, filters map[string]interface{}) ([]*models.AuditLog, int64, error) {
	return s.systemRepo.ListAuditLogs(page, size, filters)
} 