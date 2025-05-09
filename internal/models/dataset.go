package models

import (
	"time"
)

// Dataset 数据集模型
type Dataset struct {
	ID          string    `json:"id" gorm:"primaryKey;type:varchar(32)"`
	Name        string    `json:"name" gorm:"type:varchar(100);index"`
	Description string    `json:"description" gorm:"type:text"`
	Type        string    `json:"type" gorm:"type:varchar(50);index"`
	Format      string    `json:"format" gorm:"type:varchar(20)"`
	
	// 区域信息JSON存储
	RegionName  string    `json:"regionName" gorm:"type:varchar(50)"`
	RegionBounds string   `json:"regionBounds" gorm:"type:varchar(100)"` // JSON格式: [minLat, minLng, maxLat, maxLng]
	
	// 时间范围
	StartTime   *time.Time `json:"startTime" gorm:"index"`
	EndTime     *time.Time `json:"endTime" gorm:"index"`
	
	// 分辨率
	SpatialResolution  string `json:"spatialResolution" gorm:"type:varchar(50)"`
	TemporalResolution string `json:"temporalResolution" gorm:"type:varchar(50)"`
	
	// 元数据
	Size        int64     `json:"size" gorm:"default:0"` // 字节大小
	Variables   string    `json:"variables" gorm:"type:text"` // JSON格式存储变量列表
	Source      string    `json:"source" gorm:"type:varchar(100)"`
	Methodology string    `json:"methodology" gorm:"type:varchar(255)"`
	
	// 文件路径
	FilePath    string    `json:"filePath" gorm:"type:varchar(255)"`
	
	// 统计信息
	DownloadCount int       `json:"downloadCount" gorm:"default:0"`
	
	// 标签
	Tags        string    `json:"tags" gorm:"type:varchar(255)"` // 以逗号分隔
	
	// 创建和更新信息
	CreatedBy   string    `json:"createdBy" gorm:"type:varchar(32)"`
	CreatedAt   *time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   *time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

// TableName 表名
func (Dataset) TableName() string {
	return "datasets"
}

// VariableInfo 变量信息
type VariableInfo struct {
	Name        string    `json:"name"`
	Unit        string    `json:"unit"`
	Description string    `json:"description"`
	Range       [2]float64 `json:"range"` // [min, max]
}

// Region 区域信息
type Region struct {
	Name   string    `json:"name"`
	Bounds [4]float64 `json:"bounds"` // [minLat, minLng, maxLat, maxLng]
} 