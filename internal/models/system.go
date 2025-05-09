package models

import (
	"time"
)

// SystemSetting 系统设置模型
type SystemSetting struct {
	ID          string     `json:"id" gorm:"primaryKey;type:varchar(32)"`
	Category    string     `json:"category" gorm:"type:varchar(50);index"` // 设置类别
	Key         string     `json:"key" gorm:"type:varchar(100);uniqueIndex"` // 设置键名
	Value       string     `json:"value" gorm:"type:text"` // 设置值
	Description string     `json:"description" gorm:"type:varchar(255)"` // 描述
	CreatedAt   *time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   *time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

// TableName 表名
func (SystemSetting) TableName() string {
	return "system_settings"
}

// AuditLog 操作日志模型
type AuditLog struct {
	ID          uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID      string     `json:"userId" gorm:"type:varchar(32);index"`
	UserName    string     `json:"userName" gorm:"type:varchar(50)"`
	Action      string     `json:"action" gorm:"type:varchar(50);index"` // 操作类型
	Resource    string     `json:"resource" gorm:"type:varchar(50);index"` // 资源类型
	ResourceID  string     `json:"resourceId" gorm:"type:varchar(32);index"` // 资源ID
	Description string     `json:"description" gorm:"type:text"` // 详细描述
	IPAddress   string     `json:"ipAddress" gorm:"type:varchar(50)"` // IP地址
	UserAgent   string     `json:"userAgent" gorm:"type:varchar(255)"` // 用户代理
	CreatedAt   *time.Time `json:"createdAt" gorm:"autoCreateTime;index"`
}

// TableName 表名
func (AuditLog) TableName() string {
	return "audit_logs"
} 