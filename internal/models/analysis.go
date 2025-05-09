package models

import (
	"time"
)

// AnalysisTask 分析任务模型
type AnalysisTask struct {
	ID          string     `json:"id" gorm:"primaryKey;type:varchar(32)"`
	Type        string     `json:"type" gorm:"type:varchar(50);index"` // 如: temperature-salinity, wave, current等
	Name        string     `json:"name" gorm:"type:varchar(100)"`
	Description string     `json:"description" gorm:"type:text"`
	
	// 任务参数 (JSON)
	Parameters  string     `json:"parameters" gorm:"type:text"`
	
	// 关联数据集
	DatasetID   string     `json:"datasetId" gorm:"type:varchar(32);index"`
	
	// 任务状态
	Status      string     `json:"status" gorm:"type:varchar(20);index"` // pending, running, completed, failed
	Progress    int        `json:"progress" gorm:"default:0"`            // 进度百分比: 0-100
	
	// 结果和错误信息
	ResultPath  string     `json:"resultPath" gorm:"type:varchar(255)"` // 结果文件路径
	ErrorMsg    string     `json:"errorMsg" gorm:"type:text"`
	
	// 创建和更新信息
	CreatedBy   string     `json:"createdBy" gorm:"type:varchar(32);index"`
	CreatedAt   *time.Time `json:"createdAt" gorm:"autoCreateTime;index"`
	StartedAt   *time.Time `json:"startedAt"`
	CompletedAt *time.Time `json:"completedAt"`
	UpdatedAt   *time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

// TableName 表名
func (AnalysisTask) TableName() string {
	return "analysis_tasks"
}

// AnalysisResult 分析结果模型
type AnalysisResult struct {
	ID          string     `json:"id" gorm:"primaryKey;type:varchar(32)"`
	TaskID      string     `json:"taskId" gorm:"type:varchar(32);index"`
	Title       string     `json:"title" gorm:"type:varchar(100)"`
	Description string     `json:"description" gorm:"type:text"`
	
	// 结果类型和文件
	Type        string     `json:"type" gorm:"type:varchar(20)"` // 如: chart, table, map, file
	Format      string     `json:"format" gorm:"type:varchar(20)"` // 如: json, csv, netcdf, png
	FilePath    string     `json:"filePath" gorm:"type:varchar(255)"`
	
	// 预览数据
	PreviewData string     `json:"previewData" gorm:"type:text"` // 预览数据的JSON表示
	
	// 元数据
	Metadata    string     `json:"metadata" gorm:"type:text"` // 存储元数据的JSON对象
	
	// 创建和更新信息
	CreatedAt   *time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   *time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

// TableName 表名
func (AnalysisResult) TableName() string {
	return "analysis_results"
} 